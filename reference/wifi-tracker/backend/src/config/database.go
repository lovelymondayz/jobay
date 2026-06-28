package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SSLMODE"),
	)

	log.Println("DSN:", dsn) // Debug DSN
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_SSLMODE:", os.Getenv("DB_SSLMODE"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // slow SQL Threshold
			LogLevel:      logger.Info, // log level
			Colorful:      true,        // disable color
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Check connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf(" Failed to ping database: %v", err)
	}

	log.Println("Database connected and migration successful")
}
