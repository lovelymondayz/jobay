package dtos

import (
	"github.com/google/uuid"
)

type AccessPointDTO struct {
	ClientID   uuid.UUID `json:"client_id" validate:"required"`
	SSID       string    `json:"ssid" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	MACAddress string    `json:"mac_address" validate:"required,mac"`
	CreatedBy  string    `json:"created_by" validate:"required"`
	UpdatedBy  string    `json:"updated_by" validate:"required"`
}
