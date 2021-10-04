package main

/*
	db
	env
	fiber
	html

*/

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	key, err := base64.RawURLEncoding.DecodeString("eXvpIA2UF32Ubf_-FgEh-wu-w_TD8ktW2nid-wUvpc0")
	if err != nil {
		fmt.Println("session key is not valid")
	}
	sessionKey = key // login signature

	time.Sleep(time.Second * 3)
	db, err := NewDB(os.Getenv("DB"))
	if err != nil {
		fmt.Println("error connecting to db", err)
		return
	}
	err = CreateTables(db)
	if err != nil {
		fmt.Println("error creating tables", err)
		return
	}

	err = PreloadUsers(db)
	if err != nil {
		Log(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Get("/", Frontend)
	app.Static("/", "./static")
	app.Post("/login", Login(db))
	app.Post("/register", Register(db))
	app.Post("/verify", CheckNIT(db))
	app.Get("/logout", Logout(db))
	app.Post("/roman", GetRoman(db))

	app.Listen(":" + os.Getenv("PORT"))

}
