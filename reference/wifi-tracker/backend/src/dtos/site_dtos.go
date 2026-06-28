package dtos

import "github.com/google/uuid"

type SiteResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SiteRequest struct {
	SiteID   uuid.UUID `json:"site_id"`
	Name string `json:"name" validate:"required"`
}