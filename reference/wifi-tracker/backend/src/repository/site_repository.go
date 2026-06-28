package repository

import (
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SiteRepositoryInterface interface {
	FindById(tx *gorm.DB, id uuid.UUID) (*models.Sites, error)
	FindAll(tx *gorm.DB) ([]models.Sites, error)
	Create(tx *gorm.DB, site *models.Sites) error
	Update(tx *gorm.DB, site *models.Sites) error
	Delete(tx *gorm.DB, id uuid.UUID) error
}

type SiteRepository struct {}

func NewSiteRepository() SiteRepositoryInterface {
	return &SiteRepository{}
}

func (r *SiteRepository) FindById(tx *gorm.DB, id uuid.UUID) (*models.Sites, error) {
	var site models.Sites
	if err := tx.Where("site_id = ?", id).First(&site).Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *SiteRepository) FindAll(tx *gorm.DB) ([]models.Sites, error) {
	var sites []models.Sites
	if err := tx.Find(&sites).Error; err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *SiteRepository) Create(tx *gorm.DB, site *models.Sites) error {
	if err := tx.Create(site).Error; err != nil {
		return err
	}
	return nil
}

func (r *SiteRepository) Update(tx *gorm.DB, site *models.Sites) error {
	if err := tx.Save(site).Error; err != nil {
		return err
	}
	return nil
}

func (r *SiteRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	if err := tx.Where("site_id = ?", id).Delete(&models.Sites{}).Error; err != nil {
		return err
	}
	return nil
}
