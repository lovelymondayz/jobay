package service

import (
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetAll()([]models.Roles, error)
	GetByID(id uint) (*models.Roles, error)
	Create(dto *dtos.RoleDTO) error
	Update(id uint, dto *dtos.RoleDTO) error
	Delete(id uint) error
}

type RoleService struct {
	Repo repository.RoleRepositoryInterface
	DB   *gorm.DB
}

func NewRoleService(repo repository.RoleRepositoryInterface, db *gorm.DB) RoleServiceInterface {
	return &RoleService{
		Repo: repo,
		DB:   db,	
	}
}

func (s *RoleService) GetAll() ([]models.Roles, error) {
	return s.Repo.FindAll(s.DB)
}

func (s *RoleService) GetByID(id uint) (*models.Roles, error) {
	return s.Repo.FindByID(s.DB, id)
}

func (s *RoleService) Create(dto *dtos.RoleDTO) error {
	role := &models.Roles{
		Name:      dto.Name,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
	return s.Repo.Create(s.DB, role)
}

func (s *RoleService) Update( id uint, dto *dtos.RoleDTO) error {
	role, err := s.Repo.FindByID(s.DB, id)
	if err != nil {
		return err
	}

	role.Name = dto.Name
	role.UpdatedBy = dto.UpdatedBy

	return s.Repo.Update(s.DB, role)
}

func (s *RoleService) Delete(id uint) error {
	return s.Repo.Delete(s.DB, id)
}
