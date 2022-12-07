package controllers

import (
	"go-as/src/application/updateUserRoles"
	"go-as/src/domain/internals"
	"go-as/src/infrastructure/api"
	"go-as/src/infrastructure/dto"
	"go-as/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateUserRolesController struct {
	updateUserRolesUseCase *updateUserRoles.UpdateUserRolesUseCase
	useCaseExecutor        *internals.AuthorizedUseCaseExecutor
	accessTokenFinder      *api.HTTPAccessTokenFinder
	dtoDeserializer        *dto.EchoDTODeserializer
	errorTransformer       *transformers.ErrorToEchoErrorTransformer
}

func (controller *UpdateUserRolesController) Handle(c echo.Context) error {
	userEmail := c.Param("email")
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var updateUserRolesDTO dto.UpdateUserRolesDTO
	if err := controller.dtoDeserializer.Deserialize(c, &updateUserRolesDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	updateUserRolesRequest := updateUserRoles.UpdateUserRolesRequest{
		UserEmail: userEmail,
		RoleNames: updateUserRolesDTO.Roles,
	}
	ctx := c.Request().Context()
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.updateUserRolesUseCase, &updateUserRolesRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	return c.NoContent(http.StatusOK)
}

func NewUpdateUserRolesController(useCase *updateUserRoles.UpdateUserRolesUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *UpdateUserRolesController {
	return &UpdateUserRolesController{
		updateUserRolesUseCase: useCase,
		useCaseExecutor:        useCaseExecutor,
		accessTokenFinder:      accessTokenFinder,
		dtoDeserializer:        dtoDeserializer,
		errorTransformer:       errorTransformer,
	}
}
