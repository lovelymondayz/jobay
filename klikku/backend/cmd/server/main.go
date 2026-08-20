package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"klikku/internal/config"
	"klikku/internal/routes"
	"klikku/internal/scheduler"
	"klikku/internal/utils"
)

func main() {
	cfg := config.Load()

	// Database connection
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := utils.RunMigrations(db, cfg.DBName); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Storage (local filesystem)
	storage, err := utils.NewStorage("/app/storage")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:    50 * 1024 * 1024, // 50MB
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendURL + ",http://localhost:3009,http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Device-Token",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Routes
	routes.Setup(app, db, storage, cfg)

	// Start campaign scheduler
	campaignScheduler := scheduler.NewCampaignScheduler(db)
	campaignScheduler.Start()
	defer campaignScheduler.Stop()

	// Seed super admin
	if err := utils.SeedSuperAdmin(db, cfg); err != nil {
		log.Printf("Warning: Could not seed super admin: %v", err)
	}

	port := cfg.AppPort
	log.Printf("🚀 Klikku backend running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
		"code":  code,
	})
}

var _ = fmt.Sprintf
var _ = os.Exit
