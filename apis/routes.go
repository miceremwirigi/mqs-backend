package apis

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/miceremwirigi/mqs-backend/apis/hospitals"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, world")
	})
	app.Get("/home", static.New("./public"))

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	api := app.Group("/api")

	hospitalsRoutes := api.Group("/hospitals")
	hospitalHandler := hospitals.Handler{
		DB: db,
	}

	// hospitalsRoutes.Post("/", hospitalHandler.AddHospital)
	hospitalsRoutes.Get("/", nil, hospitalHandler.GetAllHospitals)
	// hospitalsRoutes.Get("/:id", hospitalHandler.GetHospital)
	// hospitalsRoutes.Put("/:id", hospitalHandler.UpdateHospital)
	// hospitalsRoutes.Delete("/:id", hospitalHandler.DeleteHospital)

}
