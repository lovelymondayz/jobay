package service

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientServiceInterface interface {
	GetAll() ([]models.Clients, error)
	GetByID(id uuid.UUID) (*models.Clients, error)
	Create(dto *dtos.ClientDTO) error
	Update(id uuid.UUID, dto *dtos.ClientDTO) error
	Delete(id uuid.UUID) error
}

type ClientService struct {
	Repo repository.ClientRepositoryInterface
	DB   *gorm.DB
}

func NewClientService(repo repository.ClientRepositoryInterface, db *gorm.DB) ClientServiceInterface {
	return &ClientService{Repo: repo, DB: db}
}

func (s *ClientService) GetAll() ([]models.Clients, error) {
	return s.Repo.FindAll(s.DB)
}

func (s *ClientService) GetByID(id uuid.UUID) (*models.Clients, error) {
	return s.Repo.FindByID(s.DB, id)
}

func (s *ClientService) Create(dto *dtos.ClientDTO) error {
	client := &models.Clients{
		SiteID:    dto.SiteID,
		Name:      dto.Name,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
	return s.Repo.Create(s.DB, client)
}

func (s *ClientService) Update(id uuid.UUID, dto *dtos.ClientDTO) error {
	client, err := s.Repo.FindByID(s.DB, id)
	if err != nil {
		return err
	}

	client.Name = dto.Name
	client.SiteID = dto.SiteID
	client.UpdatedBy = dto.UpdatedBy

	return s.Repo.Update(s.DB, client)
}

func (s *ClientService) Delete(id uuid.UUID) error {
	return s.Repo.Delete(s.DB, id)
}
