package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

var sessionKey []byte

const newSessionSQL = `insert into sessions ("token", "userId", "active") values ($1, $2, $3)`
const validateSessionSQL = `select count(s) from sessions s where s."userId" = $1 and active = true`

func NewSession(db *pgxpool.Pool, user User) (string, error) {
	var count int
	c := context.Background()
	row := db.QueryRow(c, validateSessionSQL, user.ID)
	err := row.Scan(&count)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", ErrAlreadyLogged
	}
	tokenInfo := fmt.Sprintf(`{"id":%v, "username":%q, "firstName": %q}`, user.ID, user.Username, user.FirstName)
	token, err := NewSignature([]byte(tokenInfo), sessionKey)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(c, newSessionSQL, token, user.ID, true)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewSignature(value []byte, key []byte) (string, error) {
	hasher := hmac.New(sha256.New, key)
	_, err := hasher.Write(value)
	if err != nil {
		return "", err
	}
	plain := base64.RawURLEncoding.EncodeToString(value)
	signed := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return plain + "." + signed, nil
}

//Verifies signature and returns the decoded value of the token
func VerifySignature(value string, key []byte) ([]byte, error) {
	split := strings.Split(value, ".")
	if len(split) != 2 {
		return nil, ErrWrongSignatureFormat
	}
	bts, err := base64.RawURLEncoding.DecodeString(split[0])
	if err != nil {
		return nil, err
	}
	token, err := NewSignature(bts, key)
	if err != nil {
		return nil, err
	}

	if subtle.ConstantTimeCompare([]byte(value), []byte(token)) == 1 {
		return bts, nil
	}
	return nil, ErrInvalidSignature
}
