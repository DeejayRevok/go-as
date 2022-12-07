package dto

type UpdateUserRolesDTO struct {
	Roles []string `json:"roles" validate:"required"`
}
