package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type AccessPoints struct {
	AccessPointID uuid.UUID `json:"access_point_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ClientID      uuid.UUID `json:"client_id"`
	SSID          string    `json:"ssid"`
	Name          string    `json:"name"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	CreatedBy     string     `json:"created_by"`
	UpdatedBy     string     `json:"updated_by"`
	MACAddress    string     `json:"mac_address"`
}

func (AccessPoints) TableName() string {
	return "access_points"
}