package commands

import (
	"context"
	"go-as/src/application/createPermission"
	"go-as/src/domain/permission"
	"go-as/src/domain/role"
	"go-as/src/domain/user"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var permissions [5]string = [...]string{
	permission.CreatePermissionPermission,
	role.CreateRolePermission,
	role.UpdateRolePermission,
	role.DeleteRolePermission,
	user.UpdateUserPermission,
}

type BoostrapPermissionsCLI struct {
	createPermissionUseCase *createPermission.CreatePermissionUseCase
	logger                  *zap.Logger
}

func (cli *BoostrapPermissionsCLI) Execute(_ *cli.Context) error {
	cli.logger.Info("Starting permissions bootstraping")
	defer cli.logger.Info("Finished permissions bootstraping")
	ctx := context.Background()
	for _, permission := range permissions {
		useCaseRequest := createPermission.CreatePermissionRequest{
			Name: permission,
		}
		response := cli.createPermissionUseCase.Execute(ctx, &useCaseRequest)
		if response.Err != nil {
			return response.Err
		}
	}
	return nil
}

func NewBoostrapPermissionsCLI(createPermissionUseCase *createPermission.CreatePermissionUseCase, logger *zap.Logger) *BoostrapPermissionsCLI {
	return &BoostrapPermissionsCLI{
		createPermissionUseCase: createPermissionUseCase,
		logger:                  logger,
	}
}
