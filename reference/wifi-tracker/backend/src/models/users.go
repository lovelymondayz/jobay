package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoleID    uint      `json:"role_id"`
	Role      Roles     `gorm:"references:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"role"`
	ClientID  uuid.UUID `json:"client_id"`
	Client    Clients   `gorm:"references:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"client"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" gorm:"unique" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Devices   []Devices `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsActive  int       `json:"is_active" gorm:"column:is_active;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string
	UpdatedBy string
}

func (Users) TableName() string {
	return "users"
}
