package dtos

import "github.com/google/uuid"

type ClientDTO struct {
	SiteID    uuid.UUID `json:"site_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedBy string    `json:"created_by" validate:"required"`
	UpdatedBy string    `json:"updated_by" validate:"required"`
}
