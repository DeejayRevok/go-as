package app

import (
	"fmt"
	"go-as/app/cli/commands"
	"go-as/src/application/checkUserHasPermissions"
	"go-as/src/application/createPermission"
	"go-as/src/application/createRole"
	"go-as/src/application/createUser"
	"go-as/src/application/getApplicationHealth"
	"go-as/src/application/updateUserPermissions"
	"go-as/src/application/updateUserRoles"
	"go-as/src/domain/auth"
	"go-as/src/domain/events"
	"go-as/src/domain/healthcheck"
	"go-as/src/domain/internals"
	"go-as/src/domain/permission"
	"go-as/src/domain/role"
	"go-as/src/domain/user"
	"go-as/src/infrastructure/api"
	"go-as/src/infrastructure/api/controllers"
	"go-as/src/infrastructure/api/middlewares"
	"go-as/src/infrastructure/database"
	"go-as/src/infrastructure/dto"
	"go-as/src/infrastructure/jwt"
	"go-as/src/infrastructure/logging"
	"go-as/src/infrastructure/messaging"
	"go-as/src/infrastructure/transformers"

	"github.com/streadway/amqp"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildDIContainer() dig.Container {
	container := dig.New()
	if err := container.Provide(NewAPMTracer); err != nil {
		panic(fmt.Sprintf("Error providing APM tracer to the dependency injection container: %s", err.Error()))
	}
	if err := container.Provide(logging.NewZapLogger); err != nil {
		panic(fmt.Sprintf("Error providing zap logger to the dependency injection container: %s", err.Error()))
	}
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Provide(logging.NewZapTracedLogger, dig.As(new(internals.Logger))), logger)
		handleError(container.Provide(logging.NewZapGormTracedLogger), logger)
		handleError(container.Provide(ConnectDatabase), logger)
		handleError(container.Provide(ConnectToAMQPServer), logger)
		handleError(container.Provide(LoadJWTSettings), logger)

		handleError(container.Provide(database.NewPermissionDbRepository, dig.As(new(permission.PermissionRepository))), logger)
		handleError(container.Provide(database.NewRoleDbRepository, dig.As(new(role.RoleRepository))), logger)
		handleError(container.Provide(database.NewUserDbRepository, dig.As(new(user.UserRepository))), logger)

		handleError(container.Provide(jwt.NewJWTClaimsToAccessTokenTransformer), logger)
		handleError(container.Provide(jwt.NewJWTAccessTokenDeserializer, dig.As(new(auth.AccessTokenDeserializer))), logger)

		handleError(container.Provide(transformers.NewRoleToResponseTransformer), logger)
		handleError(container.Provide(transformers.NewPermissionToResponseTransformer), logger)
		handleError(container.Provide(transformers.NewEventToAMQPMessageTransformer), logger)
		handleError(container.Provide(transformers.NewAMQPDeliveryToMapTransformer), logger)
		handleError(container.Provide(transformers.NewErrorToEchoErrorTransformer), logger)

		handleError(container.Provide(func(amqpConnection *amqp.Connection, logger *zap.Logger) *amqp.Channel {
			amqpChannel, err := amqpConnection.Channel()
			if err != nil {
				logger.Fatal("Error creating the AMQP channel")
				return nil
			}
			return amqpChannel
		}), logger)
		handleError(container.Provide(messaging.NewAMQPExchangeManager), logger)
		handleError(container.Provide(messaging.NewAMQPQueueEventListenerFactory, dig.As(new(events.EventListenerFactory))), logger)

		handleError(container.Provide(internals.NewAuthorizedUseCaseExecutor), logger)
		handleError(container.Provide(createUser.NewCreateUserUseCase), logger)
		handleError(container.Provide(createUser.NewUserCreatedEventConsumer), logger)
		handleError(container.Provide(createPermission.NewCreatePermissionUseCase), logger)
		handleError(container.Provide(createRole.NewCreateRoleUseCase), logger)
		handleError(container.Provide(checkUserHasPermissions.NewCheckUserHasPermissionUseCase), logger)
		handleError(container.Provide(updateUserPermissions.NewUpdateUserPermissionsUseCase), logger)
		handleError(container.Provide(updateUserRoles.NewUpdateUserRolesUseCase), logger)

		handleError(container.Provide(dto.NewEchoDTOSerializer), logger)
		handleError(container.Provide(dto.NewEchoDTODeserializer), logger)

		handleError(container.Provide(dto.NewDTOValidator), logger)
		handleError(container.Provide(middlewares.NewEchoLogMiddleware), logger)
		handleError(container.Provide(NewRedocConfiguration), logger)

		addHealthCheckDependencies(container, logger)

		handleError(container.Provide(api.NewHTTPAccessTokenFinder), logger)
		handleError(container.Provide(controllers.NewCreatePermissionController), logger)
		handleError(container.Provide(controllers.NewCreateRoleController), logger)
		handleError(container.Provide(controllers.NewCheckPermissionsController), logger)
		handleError(container.Provide(controllers.NewUpdateUserPermissionsController), logger)
		handleError(container.Provide(controllers.NewUpdateUserRolesController), logger)

		handleError(container.Provide(commands.NewBoostrapPermissionsCLI), logger)
	}); err != nil {
		panic(fmt.Sprintf("Error adding dependencies to the container: %s", err.Error()))
	}

	return *container
}

func addHealthCheckDependencies(diContainer *dig.Container, logger *zap.Logger) {
	type healthCheckersAggregator struct {
		dig.Out
		HealthChecker healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	handleError(diContainer.Provide(func(db *gorm.DB) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: database.NewDatabaseHealthChecker(db),
		}
	}), logger)
	handleError(diContainer.Provide(func(amqpConnection *amqp.Connection) healthCheckersAggregator {
		return healthCheckersAggregator{
			HealthChecker: messaging.NewAMQPHealthChecker(amqpConnection),
		}
	}), logger)

	type healthCheckersGroup struct {
		dig.In
		HealthCheckers []healthcheck.SingleHealthChecker `group:"healthcheckers"`
	}
	handleError(diContainer.Provide(func(checkersGroup healthCheckersGroup) *healthcheck.HealthChecker {
		return healthcheck.NewHealthChecker(checkersGroup.HealthCheckers)
	}), logger)

	handleError(diContainer.Provide(getApplicationHealth.NewGetApplicationHealthUseCase), logger)
	handleError(diContainer.Provide(controllers.NewGetStatusController), logger)
}
