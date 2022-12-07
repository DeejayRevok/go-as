package updateUserRoles

import (
	"context"
	"fmt"
	"go-as/src/domain/internals"
	"go-as/src/domain/role"
	"go-as/src/domain/user"
)

type UpdateUserRolesUseCase struct {
	userRepository user.UserRepository
	roleRepository role.RoleRepository
	logger         internals.Logger
}

func (useCase *UpdateUserRolesUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*UpdateUserRolesRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting adding roles to %s", validatedRequest.UserEmail))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished adding roles to %s", validatedRequest.UserEmail))

	user, err := useCase.userRepository.FindByEmail(ctx, validatedRequest.UserEmail)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if user == nil {
		return internals.ErrorUseCaseResponse(fmt.Errorf("user %s not found", validatedRequest.UserEmail))
	}

	roles, err := useCase.roleRepository.FindByNames(ctx, validatedRequest.RoleNames)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if len(roles) != len(validatedRequest.RoleNames) {
		return internals.ErrorUseCaseResponse(fmt.Errorf("roles %s not found", validatedRequest.RoleNames))
	}

	user.Roles = roles
	err = useCase.userRepository.Save(ctx, *user)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (*UpdateUserRolesUseCase) RequiredPermissions() []string {
	return []string{user.UpdateUserPermission}
}

func NewUpdateUserRolesUseCase(userRepository user.UserRepository, roleRepository role.RoleRepository, logger internals.Logger) *UpdateUserRolesUseCase {
	return &UpdateUserRolesUseCase{
		userRepository: userRepository,
		roleRepository: roleRepository,
		logger:         logger,
	}
}
