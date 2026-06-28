package service

import (
	"fmt"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/repository"
	"wifi-tracker-be/src/utils"

	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(data *dtos.LoginRequest) (string, error)
}

type authService struct {
	userRepo repository.UserRepositoryInterface
	db       *gorm.DB
}

func NewAuthService(db *gorm.DB,userRepo repository.UserRepositoryInterface) AuthServiceInterface {
	return &authService{
		userRepo: userRepo,
		db:       db,
	}
}

func (a *authService) Login(data *dtos.LoginRequest) (string, error) {

	user, err := a.userRepo.FindByEmail(a.db, data.Email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, _ := utils.GenerateJWT(user)
	return token, nil
}
