package controllers

import (
	"go-as/src/application/createRole"
	"go-as/src/domain/internals"
	"go-as/src/infrastructure/api"
	"go-as/src/infrastructure/dto"
	"go-as/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateRoleController struct {
	createRoleUseCase *createRole.CreateRoleUseCase
	useCaseExecutor   *internals.AuthorizedUseCaseExecutor
	accessTokenFinder *api.HTTPAccessTokenFinder
	dtoDeserializer   *dto.EchoDTODeserializer
	errorTransformer  *transformers.ErrorToEchoErrorTransformer
}

func (controller *CreateRoleController) Handle(c echo.Context) error {
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var creationRequestDTO dto.RoleCreationRequestDTO
	if err := controller.dtoDeserializer.Deserialize(c, &creationRequestDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	ctx := c.Request().Context()
	createRoleRequest := createRole.CreateRoleRequest{
		Name:        creationRequestDTO.Name,
		Permissions: creationRequestDTO.Permissions,
	}
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.createRoleUseCase, &createRoleRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusCreated)
}

func NewCreateRoleController(useCase *createRole.CreateRoleUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CreateRoleController {
	return &CreateRoleController{
		createRoleUseCase: useCase,
		useCaseExecutor:   useCaseExecutor,
		accessTokenFinder: accessTokenFinder,
		dtoDeserializer:   dtoDeserializer,
		errorTransformer:  errorTransformer,
	}
}
