package controllers

import (
	"go-as/src/application/checkUserHasPermissions"
	"go-as/src/domain/internals"
	"go-as/src/infrastructure/api"
	"go-as/src/infrastructure/dto"
	"go-as/src/infrastructure/transformers"

	"github.com/labstack/echo/v4"
)

type CheckPermissionsController struct {
	checkUserPermissionsUseCase *checkUserHasPermissions.CheckUserHasPermissionUseCase
	useCaseExecutor             *internals.AuthorizedUseCaseExecutor
	accessTokenFinder           *api.HTTPAccessTokenFinder
	dtoDeserializer             *dto.EchoDTODeserializer
	dtoSerializer               *dto.EchoDTOSerializer
	errorTransformer            *transformers.ErrorToEchoErrorTransformer
}

func (controller *CheckPermissionsController) Handle(c echo.Context) error {
	request := c.Request()
	accessToken, err := controller.accessTokenFinder.Find(request)
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}
	if err != nil {
		return controller.errorTransformer.Transform(err)
	}

	var checkPermissionsRequestDTO dto.CheckUserPermissionsRequestDTO
	if err := controller.dtoDeserializer.Deserialize(c, &checkPermissionsRequestDTO); err != nil {
		return controller.errorTransformer.Transform(err)
	}
	checkPermissionsRequest := checkUserHasPermissions.CheckUserHasPermissionRequest{
		UserEmail:       accessToken.Sub,
		PermissionNames: checkPermissionsRequestDTO.Permissions,
	}
	ctx := c.Request().Context()
	useCaseResponse := controller.useCaseExecutor.Execute(ctx, controller.checkUserPermissionsUseCase, &checkPermissionsRequest, accessToken)
	if useCaseResponse.Err != nil {
		return controller.errorTransformer.Transform(useCaseResponse.Err)
	}
	checkResponse := dto.CheckUserPermissionsResponseDTO{
		Result: useCaseResponse.Content.(bool),
	}
	return controller.dtoSerializer.Serialize(c, checkResponse)
}

func NewCheckPermissionsController(checkUserPermissionsUseCase *checkUserHasPermissions.CheckUserHasPermissionUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, dtoDeserializer *dto.EchoDTODeserializer, dtoSerializer *dto.EchoDTOSerializer, errorTransformer *transformers.ErrorToEchoErrorTransformer) *CheckPermissionsController {
	return &CheckPermissionsController{
		checkUserPermissionsUseCase: checkUserPermissionsUseCase,
		useCaseExecutor:             useCaseExecutor,
		accessTokenFinder:           accessTokenFinder,
		dtoDeserializer:             dtoDeserializer,
		dtoSerializer:               dtoSerializer,
		errorTransformer:            errorTransformer,
	}
}
