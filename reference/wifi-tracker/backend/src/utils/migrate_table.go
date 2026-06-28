package utils

import (
	"log"
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/models"
)

func MigrateTable() {
	// Migrate
	if err := config.DB.AutoMigrate(
		&models.Sites{},
		&models.Clients{},
		&models.Roles{},
		&models.Users{},
		&models.Devices{},
		&models.AccessPoints{},
		&models.ConnectionHistories{},
	); err != nil {
		log.Fatalf(" Error migrating database: %v", err)
	}
}

func DropTable() {
	// Migrate
	if err := config.DB.Migrator().DropTable(
		&models.Sites{},
		&models.Clients{},
		&models.Roles{},
		&models.Users{},
		&models.Devices{},
		&models.AccessPoints{},
		&models.ConnectionHistories{},
	); err != nil {
		log.Fatalf(" Error migrating database: %v", err)
	}
}	