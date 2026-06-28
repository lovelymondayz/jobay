package dtos

import (
	"github.com/google/uuid"
)

type UserCreateDTO struct {
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	RoleID    uint      `json:"role_id" validate:"required"`
	ClientID  uuid.UUID `json:"client_id" validate:"required"`
	CreatedBy string    `json:"created_by" validate:"required"`
	UpdatedBy string    `json:"updated_by" validate:"required"`
}

type UserUpdateDTO struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	RoleID    uint      `json:"role_id" validate:"required"`
	ClientID  uuid.UUID `json:"client_id" validate:"required"`
	CreatedBy string    `json:"created_by" validate:"required"`
	UpdatedBy string    `json:"updated_by" validate:"required"`
}

