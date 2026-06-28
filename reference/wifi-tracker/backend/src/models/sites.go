package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sites struct {
	SiteID    uuid.UUID      `json:"site_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `json:"name" validate:"required"`
	Desc      string         `json:"desc"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
}

func (Sites) TableName() string {
	return "sites"
}
