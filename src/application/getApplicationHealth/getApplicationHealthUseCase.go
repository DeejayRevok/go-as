package getApplicationHealth

import (
	"context"
	"go-as/src/domain/healthcheck"
	"go-as/src/domain/internals"
)

type GetApplicationHealthUseCase struct {
	healthChecker *healthcheck.HealthChecker
	logger        internals.Logger
}

func (useCase *GetApplicationHealthUseCase) Execute(ctx context.Context, _ any) internals.UseCaseResponse {
	useCase.logger.Info(ctx, "Starting checking if application is healthy")
	defer useCase.logger.Info(ctx, "Finished checking if application is healthy")
	return internals.UseCaseResponse{
		Err: useCase.healthChecker.Check(),
	}
}

func (*GetApplicationHealthUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewGetApplicationHealthUseCase(healthChecker *healthcheck.HealthChecker, logger internals.Logger) *GetApplicationHealthUseCase {
	useCase := GetApplicationHealthUseCase{
		healthChecker: healthChecker,
		logger:        logger,
	}
	return &useCase
}
