package repository

import (
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepositoryInterface interface {
	FindAll(tx *gorm.DB) ([]models.Clients, error)
	FindByID(tx *gorm.DB, id uuid.UUID) (*models.Clients, error)
	Create(tx *gorm.DB, client *models.Clients) error
	Update(tx *gorm.DB, client *models.Clients) error
	Delete(tx *gorm.DB, id uuid.UUID) error
}

type ClientRepository struct{}

func NewClientRepository() ClientRepositoryInterface {
	return &ClientRepository{}
}

func (r *ClientRepository) FindAll(tx *gorm.DB) ([]models.Clients, error) {
	var clients []models.Clients
	if err := tx.Preload("Site").Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) FindByID(tx *gorm.DB, id uuid.UUID) (*models.Clients, error) {
	var client models.Clients
	if err := tx.Preload("Site").Where("client_id = ?", id).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) Create(tx *gorm.DB, client *models.Clients) error {
	return tx.Create(client).Error
}

func (r *ClientRepository) Update(tx *gorm.DB, client *models.Clients) error {
	return tx.Save(client).Error
}

func (r *ClientRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("client_id = ?", id).Delete(&models.Clients{}).Error
}
