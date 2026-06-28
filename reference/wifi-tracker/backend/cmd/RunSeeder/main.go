package main

import (
	"log"
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/seeder"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to DB
	config.ConnectDatabase()

	if config.DB == nil {
		log.Fatal("Database initialization failed")
	}

	// Seed data
	seeder.SeedData(config.DB)
}