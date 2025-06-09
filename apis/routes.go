package apis

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/miceremwirigi/mqs-backend/apis/auth"
	"github.com/miceremwirigi/mqs-backend/apis/departments"
	"github.com/miceremwirigi/mqs-backend/apis/engineers"
	"github.com/miceremwirigi/mqs-backend/apis/equipments"
	"github.com/miceremwirigi/mqs-backend/apis/hospitals"
	"github.com/miceremwirigi/mqs-backend/apis/micerejwt"
	"github.com/miceremwirigi/mqs-backend/apis/services"
	"github.com/miceremwirigi/mqs-backend/apis/users"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your_secret_key") // Use the same secret as in auth.go

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/*", static.New("./public"))
	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("Pong!")
	})

	// --- Authentication routes ---
	authHandler := &auth.Handler{DB: db}
	app.Post("/auth/login", func(c fiber.Ctx) error {
		token, err := authHandler.LoginUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"token": token})
	})
	app.Post("/auth/logout", authHandler.LogoutUser)

	// --- Registration route ---
	userHandler := users.Handler{DB: db}
	// Allow registration without token ONLY if no users exist (logic is in AddUser)
	app.Post("/auth/register", userHandler.AddUser)

	// Apply JWT middleware to all /api routes
	api := app.Group("/api", micerejwt.New(micerejwt.Config{
		SigningKey:    jwtSecret,
		ContextKey:    "user",  // Store the user in context
		SigningMethod: "HS256", // Ensure this matches your JWT signing method
	}))

	// --- Protected API routes ---
	// api := app.Group("/api")

	// --- Hospitals ---
	hospitalsRoutes := api.Group("/hospitals")
	hospitalHandler := hospitals.Handler{DB: db}
	hospitalsRoutes.Post("/", hospitalHandler.AddHospital)
	hospitalsRoutes.Get("/", hospitalHandler.GetAllHospitals)
	hospitalsRoutes.Get("/:id", hospitalHandler.GetHospital)
	hospitalsRoutes.Get("/details/:id", hospitalHandler.GetHospitalHtml)
	hospitalsRoutes.Put("/:id", hospitalHandler.UpdateHospital)
	hospitalsRoutes.Post("/details/:id", hospitalHandler.UpdateHospital)
	hospitalsRoutes.Delete("/delete/:id", hospitalHandler.DeleteHospital)

	// --- Engineers ---
	engineersRoutes := api.Group("/engineers")
	engineerHandler := engineers.Handler{DB: db}
	engineersRoutes.Post("/", engineerHandler.AddEngineer)
	engineersRoutes.Get("/", engineerHandler.GetAllEngineers)
	engineersRoutes.Get("/:id", engineerHandler.GetEngineer)
	engineersRoutes.Get("/details/:id", engineerHandler.GetEngineerHtml)
	engineersRoutes.Put("/:id", engineerHandler.UpdateEngineer)
	engineersRoutes.Post("/details/:id", engineerHandler.UpdateEngineer)
	engineersRoutes.Delete("/delete/:id", engineerHandler.DeleteEngineer)

	// --- Equipments ---
	equipmentsRoutes := api.Group("/equipments")
	equipmentHandler := equipments.Handler{DB: db}
	equipmentsRoutes.Post("/", equipmentHandler.AddEquipment)
	equipmentsRoutes.Get("/", equipmentHandler.GetAllEquipments)
	equipmentsRoutes.Get("/:id", equipmentHandler.GetEquipment)
	equipmentsRoutes.Get("/details/:id", equipmentHandler.GetEquipmentHtml)
	equipmentsRoutes.Put("/:id", equipmentHandler.UpdateEquipment)
	equipmentsRoutes.Post("/details/:id", equipmentHandler.UpdateEquipment)
	equipmentsRoutes.Delete("/delete/:id", equipmentHandler.DeleteEquipment)

	// --- Services ---
	servicesRoutes := api.Group("/services")
	serviceHandler := services.Handler{DB: db}
	servicesRoutes.Post("/", serviceHandler.AddService)
	servicesRoutes.Get("/", serviceHandler.GetAllServices)
	servicesRoutes.Get("/:id", serviceHandler.GetService)
	servicesRoutes.Get("/details/:id", serviceHandler.GetServiceHtml)
	servicesRoutes.Put("/:id", serviceHandler.UpdateService)
	servicesRoutes.Post("/details/:id", serviceHandler.UpdateService)
	servicesRoutes.Delete("/delete/:id", serviceHandler.DeleteService)

	// --- Users ---
	usersRoutes := api.Group("/users")
	usersRoutes.Post("/", userHandler.AddUser)
	usersRoutes.Get("/", userHandler.GetAllUsers)
	usersRoutes.Get("/email/:email", userHandler.GetUserByEmail)
	usersRoutes.Get("/id/:id", userHandler.GetUserById)
	usersRoutes.Get("/username/:username", userHandler.GetUserByUsername)
	usersRoutes.Put("/:id", userHandler.UpdateUser)
	usersRoutes.Delete("/:id", userHandler.DeleteUser)

	// --- Departments ---
	departmentsHandler := departments.Handler{DB: db}
	departmentsRoutes := api.Group("/departments")
	departmentsRoutes.Get("/", departmentsHandler.GetAllDepartments)
	departmentsRoutes.Post("/", departmentsHandler.AddDepartment)
}
