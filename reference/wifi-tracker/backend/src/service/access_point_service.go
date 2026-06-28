package service

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessPointServiceInterface interface {
	GetAll() ([]models.AccessPoints, error)
	GetByID(id uuid.UUID) (*models.AccessPoints, error)
	Create(dto *dtos.AccessPointDTO) error
	Update(id uuid.UUID, dto *dtos.AccessPointDTO) error
	Delete(id uuid.UUID) error
}

type AccessPointService struct {
	Repo repository.AccessPointRepositoryInterface
	DB   *gorm.DB
}

func NewAccessPointService(repo repository.AccessPointRepositoryInterface, db *gorm.DB) AccessPointServiceInterface {
	return &AccessPointService{
		Repo: repo,
		DB:   db,
	}
}

func (s *AccessPointService) GetAll() ([]models.AccessPoints, error) {
	return s.Repo.FindAll(s.DB)
}

func (s *AccessPointService) GetByID(id uuid.UUID) (*models.AccessPoints, error) {
	return s.Repo.FindByID(s.DB, id)
}

func (s *AccessPointService) Create(dto *dtos.AccessPointDTO) error {
	ap := &models.AccessPoints{
		ClientID:   dto.ClientID,
		SSID:       dto.SSID,
		Name:       dto.Name,
		MACAddress: dto.MACAddress,
		CreatedBy:  dto.CreatedBy,
		UpdatedBy:  dto.UpdatedBy,
	}
	return s.Repo.Create(s.DB, ap)
}

func (s *AccessPointService) Update(id uuid.UUID, dto *dtos.AccessPointDTO) error {
	ap, err := s.Repo.FindByID(s.DB, id)
	if err != nil {
		return err
	}

	ap.ClientID = dto.ClientID
	ap.SSID = dto.SSID
	ap.Name = dto.Name
	ap.MACAddress = dto.MACAddress
	ap.UpdatedBy = dto.UpdatedBy

	return s.Repo.Update(s.DB, ap)
}

func (s *AccessPointService) Delete(id uuid.UUID) error {
	return s.Repo.Delete(s.DB, id)
}
