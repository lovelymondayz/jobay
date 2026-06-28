package dtos
type RoleDTO struct {
	Name      string `json:"name" validate:"required"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
}