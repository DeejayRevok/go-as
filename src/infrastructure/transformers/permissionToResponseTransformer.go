package transformers

import (
	"go-as/src/domain/permission"
	"go-as/src/infrastructure/dto"
)

type PermissionToResponseTransformer struct{}

func (transformer *PermissionToResponseTransformer) Transform(permission *permission.Permission) *dto.PermissionResponseDTO {
	permissionResponse := dto.PermissionResponseDTO{
		Name: permission.Name,
	}
	return &permissionResponse
}

func NewPermissionToResponseTransformer() *PermissionToResponseTransformer {
	transformer := PermissionToResponseTransformer{}
	return &transformer
}
