package main

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
)

var userFromUsernameSQL = `
	select row_to_json(u) from users u where username = $1
`

var deleteSessionSQL = `
	update sessions set active = false where token = $1
`

func Login(db *pgxpool.Pool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var login UserLogin
		err := ctx.BodyParser(&login)
		if err != nil {
			Log(err)
			return fiber.ErrBadRequest
		}
		err = login.Validate()
		if err != nil {
			Log(err)
			return fiber.ErrBadRequest
		}
		c := context.Background()

		row := db.QueryRow(c, userFromUsernameSQL, login.Username)
		var user User
		err = row.Scan(&user)
		if err != nil {
			Log(err)
			return fiber.ErrUnauthorized
		}
		if user.Password != login.Password {
			return fiber.ErrForbidden
		}
		token, err := NewSession(db, user)
		if err != nil {
			Log(err)
			return err
		}
		ctx.Set("authorization", token)
		ctx.Set("Access-Control-Expose-Headers", "authorization")
		return nil
	}
}

const registerSQL = `
	insert into users ("firstName","lastName","location","username","password","role") values ($1,$2,$3,$4,$5,$6) returning id
`

func Register(db *pgxpool.Pool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user User
		err := ctx.BodyParser(&user)
		if err != nil {
			Log(err)
			return fiber.ErrBadRequest
		}
		err = user.ValidateRegister()
		if err != nil {
			Log(err)
			return err
		}
		row := db.QueryRow(context.Background(), registerSQL, user.GetRegisterValues()...)
		err = row.Scan(&user.ID)
		if err != nil {
			Log(err)
			return ErrUsernameTaken
		}
		token, err := NewSession(db, user)
		if err != nil {
			Log(err)
			return err
		}
		ctx.Set("authorization", token)
		ctx.Set("Access-Control-Expose-Headers", "authorization")
		return ctx.JSON(user)
	}
}
func Logout(db *pgxpool.Pool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("authorization", "")
		if token == "" {
			return nil
		}
		fmt.Println("logging out", token)

		_, err := db.Exec(context.Background(), deleteSessionSQL, token)
		if err != nil {
			Log(err)
			return fiber.ErrInternalServerError
		}
		return nil
	}
}

func CheckNIT(db *pgxpool.Pool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("authorization")
		if token == "" {
			return fiber.ErrUnauthorized
		}
		value, err := VerifySignature(token, sessionKey)
		if err != nil {
			Log(err)
			return fiber.ErrUnauthorized
		}
		fmt.Println(jsoniter.MarshalToString(value))
		var nit NitRequest
		err = ctx.BodyParser(&nit)
		if err != nil {
			Log(err)
			return fiber.ErrBadRequest
		}
		res, err := ValidateNIT(nit.NIT)
		if err != nil {
			Log(err)
			return err
		}
		return ctx.JSON(res)
	}
}

func GetRoman(db *pgxpool.Pool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := struct {
			Number int `json:"number,omitempty"`
		}{}
		err := ctx.BodyParser(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}
		res := IntToRoman(req.Number)
		return ctx.SendString(res)
	}
}
