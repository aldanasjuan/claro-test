package main

import (
	"context"
	"encoding/xml"

	"github.com/jackc/pgx/v4/pgxpool"
)

var userTable = `
create table if not exists "users" (
	"id" bigserial unique primary key,
	"firstName" text,
	"lastName" text,
	"location" text,
	"username" text unique not null,
	"password" text,
	"role" text
);
`

var sessionTable = `
create table if not exists "sessions" (
	"id" bigserial unique primary key,
	"userId" bigint references users(id),
	"token" text not null,
	"active" boolean default true
)
`

type User struct {
	XMLName   xml.Name `xml:"user"`
	ID        int      `json:"id,omitempty" xml:"id,attr"`
	FirstName string   `json:"firstName,omitempty" xml:"firstName,omitempty"`
	LastName  string   `json:"lastName,omitempty" xml:"lastName,omitempty"`
	Location  string   `json:"location,omitempty" xml:"location,omitempty"`
	Username  string   `json:"username,omitempty" xml:"username,omitempty"`
	Password  string   `json:"password,omitempty" xml:"password,omitempty"`
	Role      string   `json:"role,omitempty" xml:"role,omitempty"`
}

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type UserLogin struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Session struct {
	ID     int    `json:"id,omitempty"`
	UserID int    `json:"userId,omitempty"`
	Token  string `json:"token,omitempty"`
}

func (login *UserLogin) Validate() error {
	if login.Username == "" || login.Password == "" {
		return ErrMissingRequired
	}
	return nil
}

func (user *User) ValidateRegister() error {
	if user.FirstName == "" || user.LastName == "" || user.Location == "" || user.Username == "" || user.Password == "" || user.Role == "" {
		return ErrMissingRequired
	}
	return nil
}
func (user *User) GetRegisterValues() []interface{} {
	values := []interface{}{}
	values = append(values, user.FirstName)
	values = append(values, user.LastName)
	values = append(values, user.Location)
	values = append(values, user.Username)
	values = append(values, user.Password)
	values = append(values, user.Role)
	return values
}
func CreateTables(db *pgxpool.Pool) error {
	c := context.Background()
	_, err := db.Exec(c, userTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(c, sessionTable)
	if err != nil {
		return err
	}
	return nil
}
