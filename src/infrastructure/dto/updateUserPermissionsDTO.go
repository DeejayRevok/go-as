package dto

type UpdateUserPermissionsDTO struct {
	Permissions []string `json:"permissions" validate:"required"`
}
