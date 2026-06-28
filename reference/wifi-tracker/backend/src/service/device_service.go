package service

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DevicesServiceInterface interface {
	GetAll() ([]models.Devices, error)
	GetByID(id uuid.UUID) (*models.Devices, error)
	Create(dto *dtos.DeviceDTO) error
	Update(id uuid.UUID, dto *dtos.DeviceDTO) error
	Delete(id uuid.UUID) error
}

type DevicesService struct {
	Repo repository.DeviceRepositoryInterface
	DB   *gorm.DB
}

func NewDevicesService(repo repository.DeviceRepositoryInterface, db *gorm.DB) DevicesServiceInterface {
	return &DevicesService{
		Repo: repo,
		DB:   db,
	}
}

func (s *DevicesService) GetAll() ([]models.Devices, error) {
	return s.Repo.FindAll(s.DB)
}

func (s *DevicesService) GetByID(id uuid.UUID) (*models.Devices, error) {
	return s.Repo.FindByID(s.DB, id)
}

func (s *DevicesService) Create(dto *dtos.DeviceDTO) error {
	device := &models.Devices{
		DeviceID:   uuid.New(),
		UserID:     dto.UserID,
		Name:       dto.Name,
		MacAddress: dto.MacAddress,
		DeviceType: dto.DeviceType,
		ConnectedAt: dto.ConnectedAt,
		IsActive:   dto.IsActive,
		IsOnline:   dto.IsOnline,
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
	}
	return s.Repo.Create(s.DB, device)
}

func (s *DevicesService) Update(id uuid.UUID, dto *dtos.DeviceDTO) error {
	device, err := s.Repo.FindByID(s.DB, id)
	if err != nil {
		return err
	}

	device.UserID = dto.UserID
	device.Name = dto.Name
	device.MacAddress = dto.MacAddress
	device.DeviceType = dto.DeviceType
	device.ConnectedAt = dto.ConnectedAt
	device.IsActive = dto.IsActive
	device.IsOnline = dto.IsOnline
	device.UpdatedBy = dto.UpdatedBy

	return s.Repo.Update(s.DB, device)
}

func (s *DevicesService) Delete(id uuid.UUID) error {
	return s.Repo.Delete(s.DB, id)
}
