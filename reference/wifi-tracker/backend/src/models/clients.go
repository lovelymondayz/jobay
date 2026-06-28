package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Clients struct {
	ClientID     uuid.UUID `json:"client_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SiteID       uuid.UUID
	Site         Sites          `gorm:"references:SiteID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name         string         `json:"name"`
	Users        []Users        `gorm:"foreignKey:ClientID"`
	AccessPoints []AccessPoints `gorm:"foreignKey:ClientID;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	CreatedBy    string         `json:"created_by"`
	UpdatedBy    string         `json:"updated_by"`
}

func (Clients) TableName() string {
	return "clients"
}


