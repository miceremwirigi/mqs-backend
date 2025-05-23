package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/apis"
	"github.com/miceremwirigi/mqs-backend/databases"
)

func main() {
	app := fiber.New()
	os.Setenv("env", "production")

	db := databases.StartDatabase()
	apis.RegisterRoutes(app, db)

	app.Listen(":3000")
}
