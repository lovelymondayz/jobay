package seeder

import (
	"log"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedUsers mengisi tabel users jika masih kosong, dengan relasi ke roles dan clients.
func SeedUsers(db *gorm.DB) {
	// Cek apakah tabel users sudah terisi
	var count int64
	db.Model(&models.Users{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded.")
		return
	}

	// Ambil role
	var role models.Roles
	if err := db.First(&role, "name = ?", "Admin").Error; err != nil {
		log.Println("Role 'Admin' not found:", err)
		return
	}

	// Ambil client
	var client models.Clients
	if err := db.First(&client, "name = ?", "Internal").Error; err != nil {
		log.Println("Client 'Internal' not found:", err)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword("Admin123")

	if err != nil {
		log.Println("Failed to hash password:", err)
		return
	}

	// Buat user
	user := models.Users{
		UserID:   uuid.New(),
		Name:     "Super Admin",
		Email:    "admin@example.com",
		Password: hashedPassword,
		RoleID:   role.RoleID,
		ClientID: client.ClientID,
		CreatedBy: "System",
		UpdatedBy: "System",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to seed user")
	} else {
		log.Println("Seeded user:", user.Email)
	}
}
