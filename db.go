package main

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

var preload = []byte(`<users>
<user id="111">
<firstName>Lokesh</firstName>
<lastName>Gupta</lastName>
<location>India</location>
<username>nickname1</username>
<password>admin1234</password>
<role>admin</role>
</user>
<user id="222">
<firstName>Alex</firstName>
<lastName>Gussin</lastName>
<location>Russia</location>
<username>nickname2</username>
<password>admin1235</password>
<role>user</role>
</user>
<user id="333">
<firstName>David</firstName>
<lastName>Feezor</lastName>
<location>USA</location>
<username>nickname3</username>
<password>admin1734</password>
<role>user</role>
</user>
</users>`)

const preloadSQL = `
	insert into users ("id", "firstName","lastName","location","username","password","role") values ($1,$2,$3,$4,$5,$6,$7)
`

func NewDB(connString string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.ConnectConfig(ctx, config)
	return conn, err
}

func PreloadUsers(db *pgxpool.Pool) error {
	var users = Users{}
	err := xml.Unmarshal(preload, &users)
	if err != nil {
		return err
	}
	c := context.Background()
	tx, err := db.Begin(c)
	if err != nil {
		return err
	}
	errs := []error{}
	defer tx.Rollback(c)
	for _, user := range users.Users {
		_, err := tx.Exec(c, preloadSQL, user.ID, user.FirstName, user.LastName, user.Location, user.Username, user.Password, user.Role)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		for _, e := range errs {
			Log(e)
		}
		return fmt.Errorf("errors loading users: %v", len(errs))
	}
	err = tx.Commit(c)
	if err != nil {
		return err
	}
	return nil
}
