package createUser

import (
	"context"
	"errors"
	"go-as/mocks"
	"go-as/src/domain/user"
	"go-as/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo *mocks.UserRepository
	UseCase  *CreateUserUseCase
}

func setUp(t *testing.T) testCase {
	userRepositoryMock := mocks.NewUserRepository(t)
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	return testCase{
		UserRepo: userRepositoryMock,
		UseCase:  NewCreateUserUseCase(userRepositoryMock, logger),
	}
}

func TestExecuteWrongRequest(t *testing.T) {
	testCase := setUp(t)
	request := "wrongRequest"
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteSaveError(t *testing.T) {
	testCase := setUp(t)
	saveError := errors.New("Test save error")
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(saveError)
	testIsSuperuser := false
	testEmail := "testEmail"
	request := &CreateUserRequest{
		Email:     testEmail,
		Superuser: testIsSuperuser,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != saveError {
		t.Fatal("Expected use case to return user repository save error")
	}
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == testEmail && user.Superuser == testIsSuperuser
	}))
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	testIsSuperuser := false
	testEmail := "testEmail"
	request := &CreateUserRequest{
		Email:     testEmail,
		Superuser: testIsSuperuser,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == testEmail && user.Superuser == testIsSuperuser
	}))
}
