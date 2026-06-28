package repository

import (
	"wifi-tracker-be/src/models"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	FindAll(tx *gorm.DB) ([]models.Roles, error)
	FindByID(tx *gorm.DB, id uint) (*models.Roles, error)
	Create(tx *gorm.DB, role *models.Roles) error
	Update(tx *gorm.DB, role *models.Roles) error
	Delete(tx *gorm.DB, id uint) error
}

type RoleRepository struct{}

func NewRoleRepository() RoleRepositoryInterface {
	return &RoleRepository{}
}

func (r *RoleRepository) FindAll(tx *gorm.DB) ([]models.Roles, error) {
	var roles []models.Roles
	if err := tx.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) FindByID(tx *gorm.DB, id uint) (*models.Roles, error) {
	var role models.Roles
	if err := tx.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(tx *gorm.DB, role *models.Roles) error {
	return tx.Create(role).Error
}

func (r *RoleRepository) Update(tx *gorm.DB, role *models.Roles) error {
	return tx.Save(role).Error
}

func (r *RoleRepository) Delete(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Roles{}, id).Error
}
