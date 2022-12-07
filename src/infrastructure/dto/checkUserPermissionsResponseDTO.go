package dto

type CheckUserPermissionsResponseDTO struct {
	Result bool `json:"result" validate:"required"`
}
