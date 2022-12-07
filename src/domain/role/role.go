package role

import (
	"go-as/src/domain/permission"
)

type Role struct {
	Name        string                  `gorm:"column:name;primaryKey"`
	Permissions []permission.Permission `gorm:"many2many:role_permission"`
}

func (role *Role) HasPermission(permission string) bool {
	if role.Permissions == nil {
		return false
	}

	return role.hasPermissionInPermissions(permission)
}

func (role *Role) hasPermissionInPermissions(permission string) bool {
	for _, rolePermission := range role.Permissions {
		if rolePermission.Name == permission {
			return true
		}
	}
	return false
}
