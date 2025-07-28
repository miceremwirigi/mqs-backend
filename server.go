package main

import (
	"os"
	"os/exec"

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

	// cmd := exec.Command("sh", "-c", "go run migrations/up/up.go")
	// if err := cmd.Run(); err != nil {
	// 	panic("Failed to run migrations: " + err.Error())
	// }


	// Connect to the database
	db := databases.StartDatabase()

	// Register routes
	apis.RegisterRoutes(app, db)

	// Start the cron jobs
	utils.RunCronJobs(db)

	// Start the server
	app.Listen(":3000")
}
