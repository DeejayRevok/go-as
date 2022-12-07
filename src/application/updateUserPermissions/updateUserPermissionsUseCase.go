package updateUserPermissions

import (
	"context"
	"fmt"
	"go-as/src/domain/internals"
	"go-as/src/domain/permission"
	"go-as/src/domain/user"
)

type UpdateUserPermissionsUseCase struct {
	userRepository       user.UserRepository
	permissionRepository permission.PermissionRepository
	logger               internals.Logger
}

func (useCase *UpdateUserPermissionsUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*UpdateUserPermissionsRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting updating permissions to %s", validatedRequest.UserEmail))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished updating permissions to %s", validatedRequest.UserEmail))

	user, err := useCase.userRepository.FindByEmail(ctx, validatedRequest.UserEmail)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if user == nil {
		return internals.ErrorUseCaseResponse(fmt.Errorf("user %s not found", validatedRequest.UserEmail))
	}

	permissions, err := useCase.permissionRepository.FindByNames(ctx, validatedRequest.PermissionNames)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if len(permissions) != len(validatedRequest.PermissionNames) {
		return internals.ErrorUseCaseResponse(fmt.Errorf("permissions %s not found", validatedRequest.PermissionNames))
	}

	user.Permissions = permissions
	err = useCase.userRepository.Save(ctx, *user)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (*UpdateUserPermissionsUseCase) RequiredPermissions() []string {
	return []string{user.UpdateUserPermission}
}

func NewUpdateUserPermissionsUseCase(userRepository user.UserRepository, permissionRepository permission.PermissionRepository, logger internals.Logger) *UpdateUserPermissionsUseCase {
	return &UpdateUserPermissionsUseCase{
		userRepository:       userRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}
