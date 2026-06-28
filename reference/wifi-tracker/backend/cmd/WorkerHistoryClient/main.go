package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/worker"

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

	// Context & worker start
	ctx, cancel := context.WithCancel(context.Background())
	go worker.StartAllWorkers(ctx)

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	log.Println("Shutdown signal received.")
	cancel()

	// Tunggu worker selesai
	time.Sleep(1 * time.Second)
	log.Println("Worker shutdown complete.")
}
