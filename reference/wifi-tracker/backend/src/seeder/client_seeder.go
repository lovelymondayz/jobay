package seeder

import (
	"fmt"
	"log"
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedClients(db *gorm.DB) {
	var count int64
	db.Model(&models.Clients{}).Count(&count)

	if count > 0 {
		log.Println("Clients table already has data, skipping seeder.")
		return
	}

	var sites models.Sites
	if err := db.First(&sites, "name = ?", "default").Error; err != nil {
		log.Println("Role 'Admin' not found:", err)
		return
	}

	clients := models.Clients{
		ClientID: uuid.New(),
		SiteID:   sites.SiteID,
		Name:     "Internal",
		CreatedBy: "System",
		UpdatedBy: "System",
	}

	if err := db.Create(&clients).Error; err != nil {
		log.Fatalf("Failed to seed clients")
	} else {
		fmt.Println("Seeded clients successfully.")
	}
}
