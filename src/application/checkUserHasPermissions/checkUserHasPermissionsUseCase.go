package checkUserHasPermissions

import (
	"context"
	"fmt"
	"go-as/src/domain/internals"
	"go-as/src/domain/user"
)

type CheckUserHasPermissionUseCase struct {
	userRepo user.UserRepository
	logger   internals.Logger
}

func (useCase *CheckUserHasPermissionUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CheckUserHasPermissionRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting checking permissions from user %s", validatedRequest.UserEmail))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished checking permissions from user %s", validatedRequest.UserEmail))

	user, err := useCase.userRepo.FindByEmail(ctx, validatedRequest.UserEmail)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if user == nil {
		return internals.ErrorUseCaseResponse(fmt.Errorf("user %s not found", validatedRequest.UserEmail))
	}

	return internals.UseCaseResponse{
		Content: useCase.checkUserHasPermissions(user, validatedRequest.PermissionNames),
		Err:     nil,
	}
}

func (*CheckUserHasPermissionUseCase) checkUserHasPermissions(user *user.User, permissionNames []string) bool {
	for _, permissionName := range permissionNames {
		if !user.HasPermission(permissionName) {
			return false
		}
	}
	return true
}

func (*CheckUserHasPermissionUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewCheckUserHasPermissionUseCase(userRepo user.UserRepository, logger internals.Logger) *CheckUserHasPermissionUseCase {
	return &CheckUserHasPermissionUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}
