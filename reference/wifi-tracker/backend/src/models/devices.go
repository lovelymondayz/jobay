package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Devices struct {
	DeviceID            uuid.UUID             `json:"device_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID              uuid.UUID             `json:"user_id"`
	Name                string                `json:"name"`
	MacAddress          string                `json:"mac_address" validate:"required"`
	ConnectionHistories []ConnectionHistories `gorm:"foreignKey:DeviceID;references:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DeviceType          string                `json:"device_type"`
	ConnectedAt         string                `json:"connection_at"`
	IsActive            int                   `json:"is_active" gorm:"column:is_active;default:1"`
	IsOnline            int                   `json:"is_online" gorm:"column:is_online;default:1"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"` // before time.Time
	DeletedAt           gorm.DeletedAt        `gorm:"index"`
	CreatedBy           string                `json:"created_by"`
	UpdatedBy           string                `json:"updated_by"`
}

func (Devices) TableName() string {
	return "devices"
}
