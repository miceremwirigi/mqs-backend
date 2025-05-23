package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/miceremwirigi/mqs-backend/databases"
	"github.com/miceremwirigi/mqs-backend/routes"
)

func main() {
	app := fiber.New()
	os.Setenv("env", "production")

	_ = databases.StartDatabase()
	routes.RegisterRoutes()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, world")
	})
	app.Get("/*", static.New("./front"))

	app.Listen(":3000")
}
