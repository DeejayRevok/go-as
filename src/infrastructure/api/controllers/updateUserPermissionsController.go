package controllers

import (
	"go-as/src/application/updateUserPermissions"
	"go-as/src/domain/internals"
	"go-as/src/infrastructure/api"
	"go-as/src/infrastructure/dto"
	"go-as/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateUserPermissionsController struct {
	updateUserPermissionsUseCase *updateUserPermissions.UpdateUserPermissionsUseCase
	useCaseExecutor              *internals.AuthorizedUseCaseExecutor
	accessTokenFinder            *api.HTTPAccessTokenFinder
	dtoDeserializer              *dto.EchoDTODeserializer
	errorTransformer             *transformers.ErrorToEchoErrorTransformer
}

func (controller *UpdateUserPermissionsController) Handle(c echo.Context) error {
	userEmail := c.Param("email")
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var updateUserPermissionsDTO dto.UpdateUserPermissionsDTO
	if err := controller.dtoDeserializer.Deserialize(c, &updateUserPermissionsDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	updateUserPermissionsRequest := updateUserPermissions.UpdateUserPermissionsRequest{
		UserEmail:       userEmail,
		PermissionNames: updateUserPermissionsDTO.Permissions,
	}
	ctx := c.Request().Context()
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.updateUserPermissionsUseCase, &updateUserPermissionsRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusOK)
}

func NewUpdateUserPermissionsController(useCase *updateUserPermissions.UpdateUserPermissionsUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *UpdateUserPermissionsController {
	return &UpdateUserPermissionsController{
		updateUserPermissionsUseCase: useCase,
		useCaseExecutor:              useCaseExecutor,
		accessTokenFinder:            accessTokenFinder,
		dtoDeserializer:              dtoDeserializer,
		errorTransformer:             errorTransformer,
	}
}
