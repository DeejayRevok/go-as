package app

import (
	"go-as/src/infrastructure/api/controllers"
	"go-as/src/infrastructure/api/middlewares"
	"go-as/src/infrastructure/dto"

	"github.com/labstack/echo/v4"
	"github.com/mvrilo/go-redoc"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildHTTPServer(container *dig.Container) *echo.Echo {
	server := echo.New()
	server.Use(middlewares.NewEchoCorsMiddleware())
	server.Use(middlewares.NewEchoAPMMiddleware())

	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(validator *dto.DTOValidator) {
			server.Validator = validator
		}), logger)
		handleError(container.Invoke(func(middleware *middlewares.EchoLogMiddleware) {
			server.Use(middleware.Middleware())
		}), logger)
		handleError(container.Invoke(func(redoc *redoc.Redoc) {
			server.Use(middlewares.NewEchoRedocMiddleware(*redoc))
		}), logger)

		handleError(container.Invoke(func(controller *controllers.CreateRoleController) {
			server.POST("/roles", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.CreatePermissionController) {
			server.POST("/permissions", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.GetStatusController) {
			server.GET("/status", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.CheckPermissionsController) {
			server.POST("/permissions/check", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.UpdateUserPermissionsController) {
			server.PUT("/user/:email/permissions", controller.Handle)
		}), logger)
		handleError(container.Invoke(func(controller *controllers.UpdateUserRolesController) {
			server.PUT("/user/:email/roles", controller.Handle)
		}), logger)
	}); err != nil {
		panic("Error adding HTTP API components to the dependency injection container")
	}

	return server
}
