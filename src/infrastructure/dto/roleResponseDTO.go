package dto

type RoleResponseDTO struct {
	Name        string                  `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions"`
}
