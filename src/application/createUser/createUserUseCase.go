package createUser

import (
	"context"
	"fmt"
	"go-as/src/domain/internals"
	"go-as/src/domain/user"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
	logger         internals.Logger
}

func (useCase *CreateUserUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreateUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting user creation for %s", validatedRequest.Email))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished user creation for %s", validatedRequest.Email))

	user := user.User{
		Email:     validatedRequest.Email,
		Superuser: validatedRequest.Superuser,
	}
	if err := useCase.userRepository.Save(ctx, user); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (*CreateUserUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewCreateUserUseCase(userRepository user.UserRepository, logger internals.Logger) *CreateUserUseCase {
	useCase := CreateUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
	return &useCase
}
