package app

import (
	"fmt"
	"go-as/src/application/createUser"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func RunEventConsumers(container *dig.Container) {
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(consumer *createUser.UserCreatedEventConsumer) {
			if err := consumer.Consume(); err != nil {
				logger.Fatal(fmt.Sprintf("Error running UserCreatedEventConsumer: %s", err.Error()))
			}
		}), logger)
	}); err != nil {
		panic(fmt.Sprintf("Error adding event consumers to the dependency container %s", err.Error()))
	}
}
