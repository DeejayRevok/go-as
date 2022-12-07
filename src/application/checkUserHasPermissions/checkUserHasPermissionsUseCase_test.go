package checkUserHasPermissions

import (
	"context"
	"errors"
	"go-as/mocks"
	"go-as/src/domain/permission"
	"go-as/src/domain/role"
	"go-as/src/domain/user"
	"go-as/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo *mocks.UserRepository
	UseCase  *CheckUserHasPermissionUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepoMock := mocks.NewUserRepository(t)
	return testCase{
		UserRepo: userRepoMock,
		UseCase:  NewCheckUserHasPermissionUseCase(userRepoMock, logger),
	}
}

func TestExecuteWrongRequest(t *testing.T) {
	testCase := setUp(t)
	request := "wrongRequest"
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{},
	}
	ctx := context.Background()
	testError := errors.New("Test error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as the find user one")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}

func TestExecuteUserNotFound(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{},
	}
	ctx := context.Background()
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}

func TestExecuteSuperuserHasPermissions(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{"testPermission1", "testPermission2"},
	}
	ctx := context.Background()
	testUser := user.User{
		Email:     "testEmail",
		Superuser: true,
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if !response.Content.(bool) {
		t.Fatal("Expected use case to return true")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}

func TestExecuteUserHasPermissions(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{"testPermission1", "testPermission2"},
	}
	ctx := context.Background()
	testUser := user.User{
		Email:       "testEmail",
		Superuser:   false,
		Permissions: []permission.Permission{{Name: "testPermission1"}, {Name: "testPermission2"}},
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if !response.Content.(bool) {
		t.Fatal("Expected use case to return true")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}

func TestExecuteUserHasPermissionsInRoles(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{"testPermission1", "testPermission2"},
	}
	ctx := context.Background()
	testUser := user.User{
		Email:       "testEmail",
		Superuser:   false,
		Permissions: []permission.Permission{},
		Roles:       []role.Role{{Name: "testRole", Permissions: []permission.Permission{{Name: "testPermission1"}, {Name: "testPermission2"}}}},
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if !response.Content.(bool) {
		t.Fatal("Expected use case to return true")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}

func TestExecuteUserHasNotPermissions(t *testing.T) {
	testCase := setUp(t)
	request := CheckUserHasPermissionRequest{
		UserEmail:       "testEmail",
		PermissionNames: []string{"testPermission1", "testPermission2"},
	}
	ctx := context.Background()
	testUser := user.User{
		Email:       "testEmail",
		Superuser:   false,
		Permissions: []permission.Permission{},
		Roles:       []role.Role{},
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if response.Content.(bool) {
		t.Fatal("Expected use case to return false")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
}
