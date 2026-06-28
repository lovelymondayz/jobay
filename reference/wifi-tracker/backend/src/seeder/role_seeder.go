package seeder

import (
	"log"
	"wifi-tracker-be/src/models"

	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) {
	var count int64
	db.Model(&models.Roles{}).Count(&count)

	if count > 0 {
		log.Println("Roles table already has data, skipping seeder.")
		return
	}

	roles := []models.Roles{
		{
			Name: "Admin",
			CreatedBy: "System",
			UpdatedBy: "System",
		},
		{
			Name: "User",
			CreatedBy: "System",
			UpdatedBy: "System",
		},
	}

	if err := db.Create(&roles).Error; err != nil {
		log.Fatalf("Failed to seed roles")
	} else {
		log.Println("Seeded roles successfully.")
	}
}
