package updateUserRoles

import (
	"context"
	"errors"
	"go-as/mocks"
	"go-as/src/domain/internals"
	"go-as/src/domain/role"
	"go-as/src/domain/user"
	"go-as/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo *mocks.UserRepository
	RoleRepo *mocks.RoleRepository
	UseCase  *UpdateUserRolesUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	roleRepoMock := mocks.NewRoleRepository(t)
	userRepoMock := mocks.NewUserRepository(t)
	return testCase{
		UserRepo: userRepoMock,
		RoleRepo: roleRepoMock,
		UseCase:  NewUpdateUserRolesUseCase(userRepoMock, roleRepoMock, logger),
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
	testCase.RoleRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteUserNotFound(t *testing.T) {
	testCase := setUp(t)
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: make([]string, 0),
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.RoleRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: make([]string, 0),
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
	testCase.RoleRepo.AssertNotCalled(t, "FindByNames")
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteFindRolesError(t *testing.T) {
	testCase := setUp(t)
	testRoleNames := []string{"testRole1", "testRole2"}
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: testRoleNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testError := errors.New("Test error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testCase.RoleRepo.On("FindByNames", mock.Anything, mock.Anything).Return(nil, testError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as the find roles error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.RoleRepo.AssertCalled(t, "FindByNames", ctx, request.RoleNames)
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteRolesNotFound(t *testing.T) {
	testCase := setUp(t)
	testRoleNames := []string{"testRole1", "testRole2"}
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: testRoleNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testCase.RoleRepo.On("FindByNames", mock.Anything, mock.Anything).Return(make([]role.Role, 0), nil)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, request.UserEmail)
	testCase.RoleRepo.AssertCalled(t, "FindByNames", ctx, request.RoleNames)
	testCase.UserRepo.AssertNotCalled(t, "Save")
}

func TestExecuteSaveUserError(t *testing.T) {
	testCase := setUp(t)
	testRoleNames := []string{"testRole1", "testRole2"}
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: testRoleNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testRoles := []role.Role{{Name: "testRole1"}, {Name: "testRole2"}}
	testCase.RoleRepo.On("FindByNames", mock.Anything, mock.Anything).Return(testRoles, nil)
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
	testCase.RoleRepo.AssertCalled(t, "FindByNames", ctx, request.RoleNames)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == request.UserEmail && reflect.DeepEqual(user.Roles, testRoles)
	}))
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testRoleNames := []string{"testRole1", "testRole2"}
	request := UpdateUserRolesRequest{
		UserEmail: "testEmail",
		RoleNames: testRoleNames,
	}
	testUser := &user.User{
		Email: "testEmail",
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(testUser, nil)
	testRoles := []role.Role{{Name: "testRole1"}, {Name: "testRole2"}}
	testCase.RoleRepo.On("FindByNames", mock.Anything, mock.Anything).Return(testRoles, nil)
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
	testCase.RoleRepo.AssertCalled(t, "FindByNames", ctx, request.RoleNames)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == request.UserEmail && reflect.DeepEqual(user.Roles, testRoles)
	}))
}
