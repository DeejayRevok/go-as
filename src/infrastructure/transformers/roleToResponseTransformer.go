package transformers

import (
	"go-as/src/domain/permission"
	"go-as/src/domain/role"
	"go-as/src/infrastructure/dto"
)

type RoleToResponseTransformer struct {
	permissionTransformer *PermissionToResponseTransformer
}

func (transformer *RoleToResponseTransformer) Transform(role *role.Role) *dto.RoleResponseDTO {
	roleResponse := dto.RoleResponseDTO{
		Name:        role.Name,
		Permissions: transformer.transformPermissions(role.Permissions),
	}
	return &roleResponse
}

func (transformer *RoleToResponseTransformer) transformPermissions(permissions []permission.Permission) []dto.PermissionResponseDTO {
	var permissionResponses []dto.PermissionResponseDTO
	for _, permission := range permissions {
		permissionResponses = append(permissionResponses, *transformer.permissionTransformer.Transform(&permission))
	}
	return permissionResponses
}

func NewRoleToResponseTransformer(permissionTransformer *PermissionToResponseTransformer) *RoleToResponseTransformer {
	transformer := RoleToResponseTransformer{
		permissionTransformer: permissionTransformer,
	}
	return &transformer
}
