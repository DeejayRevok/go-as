package role

import (
	"context"
)

type RoleRepository interface {
	Save(ctx context.Context, role Role) error
	FindByNames(ctx context.Context, roleNames []string) ([]Role, error)
}
