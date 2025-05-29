package apis

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/miceremwirigi/mqs-backend/apis/engineers"
	"github.com/miceremwirigi/mqs-backend/apis/equipments"
	"github.com/miceremwirigi/mqs-backend/apis/hospitals"
	"github.com/miceremwirigi/mqs-backend/apis/services"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, world")
	})
	app.Get("/*", static.New("./public"))

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	api := app.Group("/api")

	hospitalsRoutes := api.Group("/hospitals")
	hospitalHandler := hospitals.Handler{
		DB: db,
	}

	hospitalsRoutes.Post("/", hospitalHandler.AddHospital)
	hospitalsRoutes.Get("/", hospitalHandler.GetAllHospitals)
	hospitalsRoutes.Get("/:id", hospitalHandler.GetHospital)
	hospitalsRoutes.Get("/details/:id", hospitalHandler.GetHospitalHtml)
	hospitalsRoutes.Put("/:id", hospitalHandler.UpdateHospital)
	hospitalsRoutes.Post("/details/:id", hospitalHandler.UpdateHospital)
	hospitalsRoutes.Delete("/delete/:id", hospitalHandler.DeleteHospital)

	engineersRoutes := api.Group("/engineers")
	engineerHandler := engineers.Handler{
		DB: db,
	}

	engineersRoutes.Post("/", engineerHandler.AddEngineer)
	engineersRoutes.Get("/", engineerHandler.GetAllEngineers)
	engineersRoutes.Get("/:id", engineerHandler.GetEngineer)
	engineersRoutes.Get("/details/:id", engineerHandler.GetEngineerHtml)
	engineersRoutes.Put("/:id", engineerHandler.UpdateEngineer)
	engineersRoutes.Post("/details/:id", engineerHandler.UpdateEngineer)
	engineersRoutes.Delete("/delete/:id", engineerHandler.DeleteEngineer)

	equipmentsRoutes := api.Group("/equipments")
	equipmentHandler := equipments.Handler{
		DB: db,
	}

	equipmentsRoutes.Post("/", equipmentHandler.AddEquipment)
	equipmentsRoutes.Get("/", equipmentHandler.GetAllEquipments)
	equipmentsRoutes.Get("/:id", equipmentHandler.GetEquipment)
	equipmentsRoutes.Get("/details/:id", equipmentHandler.GetEquipmentHtml)
	equipmentsRoutes.Put("/:id", equipmentHandler.UpdateEquipment)
	equipmentsRoutes.Post("/details/:id", equipmentHandler.UpdateEquipment)
	equipmentsRoutes.Delete("/delete/:id", equipmentHandler.DeleteEquipment)

	servicesRoutes := api.Group("/services")
	serviceHandler := services.Handler{
		DB: db,
	}

	servicesRoutes.Post("/", serviceHandler.AddService)
	servicesRoutes.Get("/", serviceHandler.GetAllServices)
	servicesRoutes.Get("/:id", serviceHandler.GetService)
	servicesRoutes.Get("/details/:id", serviceHandler.GetServiceHtml)
	servicesRoutes.Put("/:id", serviceHandler.UpdateService)
	servicesRoutes.Post("/details/:id", serviceHandler.UpdateService)
	servicesRoutes.Delete("/delete/:id", serviceHandler.DeleteService)

}
