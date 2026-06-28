package repository

import (
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceRepositoryInterface interface {
	FindByID(tx *gorm.DB, id uuid.UUID) (*models.Devices, error)
	FindAll(tx *gorm.DB) ([]models.Devices, error)
	Create(tx *gorm.DB, device *models.Devices) error
	Update(tx *gorm.DB, device *models.Devices) error
	Delete(tx *gorm.DB, id uuid.UUID) error
	GetUDeviceIDByMAC(tx *gorm.DB,macs []string) (map[string]uuid.UUID, error)
}

type DeviceRepository struct {
}

func NewDeviceRepository() DeviceRepositoryInterface {
	return &DeviceRepository{}
}


func (r *DeviceRepository) FindAll(tx *gorm.DB) ([]models.Devices, error) {
	var devices []models.Devices
	if err := tx.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *DeviceRepository) FindByID(tx *gorm.DB, id uuid.UUID) (*models.Devices, error) {
	var device models.Devices
	if err := tx.Where("device_id = ?", id).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Create(tx *gorm.DB, device *models.Devices) error {
	return tx.Create(device).Error
}

func (r *DeviceRepository) Update(tx *gorm.DB, device *models.Devices) error {
	return tx.Save(device).Error
}

func (r *DeviceRepository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("device_id = ?", id).Delete(&models.Devices{}).Error
}

func (r *DeviceRepository) GetUDeviceIDByMAC(tx *gorm.DB,macs []string) (map[string]uuid.UUID, error) {
	var results []struct {
		MACAddress string
		DeviceID   uuid.UUID
	}

	err := tx.Table("devices").
		Select("devices.device_id as device_id, devices.mac_address").
		Where("devices.mac_address IN ?", macs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	deviceMap := make(map[string]uuid.UUID)
	for _, r := range results {
		deviceMap[r.MACAddress] = r.DeviceID
	}

	return deviceMap, nil
}



