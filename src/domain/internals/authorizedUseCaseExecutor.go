package internals

import (
	"context"
	"errors"
	"go-as/src/domain/auth"
	"go-as/src/domain/user"
)

type AuthorizedUseCaseExecutor struct {
	userRepository user.UserRepository
}

func (executor *AuthorizedUseCaseExecutor) Execute(ctx context.Context, useCase UseCase, useCaseRequest any, accessToken *auth.AccessToken) *UseCaseResponse {
	requiredPermissions := useCase.RequiredPermissions()
	if len(requiredPermissions) > 0 {
		if err := executor.checkPermissions(ctx, useCase, accessToken, requiredPermissions); err != nil {
			useCaseResponse := UseCaseResponse{
				Err: err,
			}
			return &useCaseResponse
		}
	}

	useCaseResponse := useCase.Execute(ctx, useCaseRequest)
	return &useCaseResponse
}

func (executor *AuthorizedUseCaseExecutor) checkPermissions(ctx context.Context, useCase UseCase, token *auth.AccessToken, permissions []string) error {
	if token == nil {
		return errors.New("authentication required")
	}

	var user *user.User
	var err error
	if token != nil {
		user, err = executor.getUserFromAccessToken(ctx, token)
	}
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("authentication required")
	}

	if user.Superuser {
		return nil
	}

	for _, permissionName := range permissions {
		if !user.HasPermission(permissionName) {
			return UseCaseAuthorizationError{
				Email:      user.Email,
				Permission: permissionName,
			}
		}
	}

	return nil
}

func (executor *AuthorizedUseCaseExecutor) getUserFromAccessToken(ctx context.Context, token *auth.AccessToken) (*user.User, error) {
	return executor.userRepository.FindByEmail(ctx, token.Sub)
}

func NewAuthorizedUseCaseExecutor(userRepository user.UserRepository) *AuthorizedUseCaseExecutor {
	return &AuthorizedUseCaseExecutor{
		userRepository: userRepository,
	}
}
