package updateUserPermissions

import (
	"context"
	"errors"
	"go-as/mocks"
	"go-as/src/domain/internals"
	"go-as/src/domain/permission"
	"go-as/src/domain/user"
	"go-as/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo       *mocks.UserRepository
	PermissionRepo *mocks.PermissionRepository
	UseCase        *UpdateUserPermissionsUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	permissionRepoMock := mocks.NewPermissionRepository(t)
	userRepoMock := mocks.NewUserRepository(t)
	return testCase{
		UserRepo:       userRepoMock,
		PermissionRepo: permissionRepoMock,
		UseCase:        NewUpdateUserPermissionsUseCase(userRepoMock, permissionRepoMock, logger),
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
	testCase.PermissionRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteUserNotFound(t *testing.T) {
	testCase := setUp(t)
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: make([]string, 0),
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: make([]string, 0),
	}
	testError := errors.New("Test error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, testError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as the find user error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteFindPermissionsError(t *testing.T) {
	testCase := setUp(t)
	testPermissionNames := []string{"testPermission1", "testPermission2"}
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: testPermissionNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testError := errors.New("Test error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(nil, testError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as the find permissions error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, request.PermissionNames)
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecutePermissionsNotFound(t *testing.T) {
	testCase := setUp(t)
	testPermissionNames := []string{"testPermission1", "testPermission2"}
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: testPermissionNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(make([]permission.Permission, 0), nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, request.PermissionNames)
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteSaveUserError(t *testing.T) {
	testCase := setUp(t)
	testPermissionNames := []string{"testPermission1", "testPermission2"}
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: testPermissionNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testPermissions := []permission.Permission{{Name: "testPermission1"}, {Name: "testPermission2"}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(testPermissions, nil)
	testError := errors.New("Test error")
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(testError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as save user error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, request.PermissionNames)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == request.UserEmail && reflect.DeepEqual(user.Permissions, testPermissions)
	}))
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testPermissionNames := []string{"testPermission1", "testPermission2"}
	request := UpdateUserPermissionsRequest{
		UserEmail:       "testEmail",
		PermissionNames: testPermissionNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testPermissions := []permission.Permission{{Name: "testPermission1"}, {Name: "testPermission2"}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(testPermissions, nil)
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if response != internals.EmptyUseCaseResponse() {
		t.Fatal("Expected use case to return an empty response")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, request.PermissionNames)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == request.UserEmail && reflect.DeepEqual(user.Permissions, testPermissions)
	}))
}
