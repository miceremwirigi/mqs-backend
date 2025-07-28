package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/miceremwirigi/mqs-backend/apis"
	"github.com/miceremwirigi/mqs-backend/databases"
	"github.com/miceremwirigi/mqs-backend/utils"
)

func main() {
	// Create new engine
	engine := html.New("./public", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	os.Setenv("env", "production")

	db := databases.StartDatabase()
	apis.RegisterRoutes(app, db)
	utils.RunCronJobs(db)

	app.Listen(":3000")
}
