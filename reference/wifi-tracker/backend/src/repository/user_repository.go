package repository

import (
	"fmt"
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindAll(tx *gorm.DB) ([]models.Users, error)
	FindByEmail(tx *gorm.DB, email string) (*models.Users, error)
	FindByID(tx *gorm.DB, id uuid.UUID) (*models.Users, error)
	Create(tx *gorm.DB, user *models.Users) error
	Update(tx *gorm.DB, role *models.Users) error
	Delete(tx *gorm.DB, id uuid.UUID) error
	GetUserByDeviceMAC(tx *gorm.DB,macs []string) (map[string]string, error)
	GetUserIdByDeviceMAC(tx *gorm.DB,macs []string) (map[string]uuid.UUID, error)
}

type userRepository struct{}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() UserRepositoryInterface {
	return &userRepository{}
}

// FindByEmail finds a user by email
func (r *userRepository) FindAll(tx *gorm.DB) ([]models.Users, error) {
	var users []models.Users
	err := tx.Preload("Role").Preload("Client").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(tx *gorm.DB, id uuid.UUID) (*models.Users, error) {
	var user models.Users
	err := tx.Preload("Role").Preload("Client").Where("user_id = ?", id).First(&user).Error
	return &user, err
}
func (r *userRepository) FindByEmail(tx *gorm.DB, email string) (*models.Users, error) {
	var user models.Users
	if err := tx.
	Preload("Client").
	Preload("Role").
	Where("email = ?", email).
	First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	return &user, nil
}

// Create creates a new user in the database
func (r *userRepository) Create(tx *gorm.DB, user *models.Users) error {
	return tx.Create(user).Error
}


/*************  ✨ Windsurf Command ⭐  *************/
// Update updates an existing user in the database
/*******  7d569d09-536c-481a-ba5e-a3684d722c3c  *******/func (r *userRepository) Update(tx *gorm.DB, user *models.Users) error {
	return tx.Save(user).Error
}

func (r *userRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("user_id = ?", id).Delete(&models.Users{}).Error
}

func (r *userRepository) GetUserByDeviceMAC(tx *gorm.DB,macs []string) (map[string]string, error) {
	var results []struct {
		MACAddress string
		UserName   string
	}

	err := tx.Table("devices").
		Select("devices.mac_address, users.name as user_name, users.user_id").
		Joins("JOIN users ON users.user_id = devices.user_id").
		Where("devices.mac_address IN ?", macs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)
	for _, r := range results {
		userMap[r.MACAddress] = r.UserName
	}

	return userMap, nil
}
func (r *userRepository) GetUserIdByDeviceMAC(tx *gorm.DB,macs []string) (map[string]uuid.UUID, error) {
	var results []struct {
		MACAddress string
		UserID   uuid.UUID
	}

	err := tx.Table("devices").
		Select("devices.mac_address, users.user_id as user_id").
		Joins("JOIN users ON users.user_id = devices.user_id").
		Where("devices.mac_address IN ?", macs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	userMap := make(map[string]uuid.UUID)
	for _, r := range results {
		userMap[r.MACAddress] = r.UserID
	}

	return userMap, nil
}
