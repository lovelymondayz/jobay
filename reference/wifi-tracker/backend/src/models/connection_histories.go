package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConnectionHistories struct {
	ConnectionHistoryID uuid.UUID `json:"connection_history_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:uuid column:user_id˝"`
	DeviceID            uuid.UUID `json:"device_id" gorm:"type:uuid"`
	// Devices             Devices      `gorm:"foreignKey:DeviceID;references:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FromAPs         uuid.UUID    `json:"from_aps" gorm:"column:from_aps;type:uuid"`
	FromAccessPoint AccessPoints `gorm:"foreignKey:FromAPs;references:AccessPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ToAPs           uuid.UUID    `json:"to_aps" gorm:"column:to_aps;type:uuid"`
	ToAccessPoint   AccessPoints `gorm:"foreignKey:ToAPs;references:AccessPointID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	CreatedBy       string         `json:"created_by"`
	UpdatedBy       string         `json:"updated_by"`
}

func (ConnectionHistories) TableName() string {
	return "connection_histories"
}
