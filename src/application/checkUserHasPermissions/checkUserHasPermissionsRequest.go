package checkUserHasPermissions

type CheckUserHasPermissionRequest struct {
	UserEmail       string
	PermissionNames []string
}
