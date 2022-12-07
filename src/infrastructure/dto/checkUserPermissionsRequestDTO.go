package dto

type CheckUserPermissionsRequestDTO struct {
	Permissions []string `json:"permissions" validate:"required"`
}
