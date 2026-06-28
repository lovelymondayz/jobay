package main

import (
	"log"
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/middleware"
	"wifi-tracker-be/src/routes"
	"wifi-tracker-be/src/utils"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	config.ConnectDatabase()
	if config.DB == nil {
		log.Fatal("Database initialization failed")
	}

	// Init UniFi clients
	config.InitUniFiClient()
	config.InitUniFiClient2()

	// Fiber config
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return utils.JSONError(c, code, err.Error())
		},
	})

	// Middleware
	app.Use("/ws", middleware.WebSocketProtected())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(middleware.CustomRecoverPanic())

	// Swagger
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	if err := app.Listen(":8081"); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
