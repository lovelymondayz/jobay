package dtos

import "github.com/google/uuid"

type DeviceDTO struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	MacAddress  string    `json:"mac_address" validate:"required,mac"`
	DeviceType  string    `json:"device_type" validate:"required"`
	ConnectedAt string    `json:"connected_at" validate:"required"`
	IsActive    int       `json:"is_active" validate:"oneof=0 1"`
	IsOnline    int       `json:"is_online" validate:"oneof=0 1"`
	CreatedBy   string    `json:"created_by" validate:"required"`
	UpdatedBy   string    `json:"updated_by" validate:"required"`
}
