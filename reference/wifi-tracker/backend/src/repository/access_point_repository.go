package repository

import (
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessPointRepositoryInterface interface {
	FindAll(tx *gorm.DB) ([]models.AccessPoints, error)
	FindByID(tx *gorm.DB, id uuid.UUID) (*models.AccessPoints, error)
	Create(tx *gorm.DB, accessPoint *models.AccessPoints) error
	Update(tx *gorm.DB, accessPoint *models.AccessPoints) error
	Delete(tx *gorm.DB, id uuid.UUID) error
	GetAccessPointByMAC(tx *gorm.DB, macs []string) (map[string]string, error)
	GetAccessPointIdByMAC(tx *gorm.DB, macs []string) (map[string]uuid.UUID, error)
}

type AccessPointRepository struct{}

func NewAccessPointRepository() AccessPointRepositoryInterface {
	return &AccessPointRepository{}
}


func (r *AccessPointRepository) FindAll(tx *gorm.DB) ([]models.AccessPoints, error) {
	var aps []models.AccessPoints
	err := tx.Find(&aps).Error
	return aps, err
}

func (r *AccessPointRepository) FindByID(tx *gorm.DB, id uuid.UUID) (*models.AccessPoints, error) {
	var ap models.AccessPoints
	err := tx.Where("access_point_id = ?", id).First(&ap).Error
	return &ap, err
}

func (r *AccessPointRepository) Create(tx *gorm.DB, ap *models.AccessPoints) error {
	return tx.Create(ap).Error
}

func (r *AccessPointRepository) Update(tx *gorm.DB, ap *models.AccessPoints) error {
	return tx.Save(ap).Error
}

func (r *AccessPointRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("access_point_id = ?", id).Delete(&models.AccessPoints{}).Error
}

func (r *AccessPointRepository) GetAccessPointByMAC(tx *gorm.DB, macs []string) (map[string]string, error) {
	var results []struct {
		MACAddress string
		Name       string
	}

	err := tx.Table("access_points").
		Select("mac_address, name").
		Where("mac_address IN ?", macs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	apMap := make(map[string]string)
	for _, r := range results {
		apMap[r.MACAddress] = r.Name
	}

	return apMap, nil
}

func (r *AccessPointRepository) GetAccessPointIdByMAC(tx *gorm.DB, macs []string) (map[string]uuid.UUID, error) {
	var results []struct {
		MACAddress string
		AccessPointID   uuid.UUID
	}

	err := tx.Table("access_points").
		Select("mac_address, access_point_id").
		Where("mac_address IN ?", macs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	apMap := make(map[string]uuid.UUID)
	for _, r := range results {
		apMap[r.MACAddress] = r.AccessPointID
	}

	return apMap, nil
}