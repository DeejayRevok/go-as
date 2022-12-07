package user

import (
	"go-as/src/domain/permission"
	"go-as/src/domain/role"
)

type User struct {
	Email       string                  `gorm:"column:email;primaryKey"`
	Roles       []role.Role             `gorm:"many2many:user_role"`
	Superuser   bool                    `gorm:"column:superuser"`
	Permissions []permission.Permission `gorm:"many2many:user_permission"`
}

func (user *User) HasPermission(permission string) bool {
	if user.Superuser {
		return true
	}
	if user.Permissions == nil {
		return false
	}

	hasPermission := user.hasPermissionInPermissions(permission)
	if !hasPermission {
		hasPermission = user.hasPermissionInRoles(permission)
	}
	return hasPermission
}

func (user *User) hasPermissionInPermissions(permission string) bool {
	for _, userPermission := range user.Permissions {
		if userPermission.Name == permission {
			return true
		}
	}
	return false
}

func (user *User) hasPermissionInRoles(permission string) bool {
	for _, role := range user.Roles {
		if role.HasPermission(permission) {
			return true
		}
	}
	return false
}
