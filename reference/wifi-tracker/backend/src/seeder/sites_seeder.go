package seeder

import (
	"fmt"
	"log"
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedSites(db *gorm.DB) {
	var count int64
	db.Model(&models.Sites{}).Count(&count)

	if count > 0 {
		log.Println("Sites table already has data, skipping seeder.")
		return
	}
	sites := models.Sites{
		SiteID: uuid.New(),
		Name:     "default",
		Desc:     "-",
		CreatedBy: "System",
		UpdatedBy: "System",
	}

	if err := db.Create(&sites).Error; err != nil {
		log.Fatalf("Failed to seed clients")
	} else {
		fmt.Println("Seeded clients successfully.")
	}
}
