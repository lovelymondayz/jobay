package service

import (
	"log"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"
	"wifi-tracker-be/src/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetAll() ([]models.Users, error)
	GetByID(id uuid.UUID) (*models.Users, error)
	Create(dto *dtos.UserCreateDTO) error
	Update(id uuid.UUID, dto *dtos.UserUpdateDTO) error
	Delete(id uuid.UUID) error
}

type UserService struct {
	Repo repository.UserRepositoryInterface
	DB   *gorm.DB
}

func NewUserService(repo repository.UserRepositoryInterface, db *gorm.DB) UserServiceInterface {
	return &UserService{Repo: repo, DB: db}
}

func (s *UserService) GetAll() ([]models.Users, error) {
	return s.Repo.FindAll(s.DB)
}

func (s *UserService) GetByID(id uuid.UUID) (*models.Users, error) {
	return s.Repo.FindByID(s.DB, id)
}

func (s *UserService) Create(dto *dtos.UserCreateDTO) error {

	hashedPassword, err := utils.HashPassword(dto.Password)

	if err != nil {
		log.Println("Failed to hash password:", err)
		return err
	}
	user := &models.Users{
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  hashedPassword,
		RoleID:    dto.RoleID,
		ClientID:  dto.ClientID,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}
	return s.Repo.Create(s.DB, user)
}

func (s *UserService) Update(id uuid.UUID, dto *dtos.UserUpdateDTO) error {
	user, err := s.Repo.FindByID(s.DB, id)
	if err != nil {
		return err
	}

	user.Name = dto.Name
	user.Email = dto.Email
	user.RoleID = dto.RoleID
	user.ClientID = dto.ClientID
	user.UpdatedBy = dto.UpdatedBy

	return s.Repo.Update(s.DB, user)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.Repo.Delete(s.DB, id)
}
