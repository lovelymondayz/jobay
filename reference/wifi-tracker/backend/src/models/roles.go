package models

import (
	"time"

	"gorm.io/gorm"
)

// type Roles struct {
// 	RoleID    uint    `json:"role_id" gorm:"primaryKey"`
// 	Users     []Users `gorm:"foreignKey:RoleID"`
// 	Name      string  `json:"name"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// 	CreatedBy string  `json:"created_by"`
// 	UpdatedBy string  `json:"updated_by"`
// }

type Roles struct {
	RoleID    uint   `json:"role_id" gorm:"primaryKey"`
	Name      string `json:"name" validate:"required"`
	Users     []Users `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
}

func (Roles) TableName() string {
	return "roles"
}
