package createUser

import (
	"context"
	"fmt"
	"go-as/src/domain/events"
	"go-as/src/domain/user"
	"reflect"

	"go.uber.org/zap"
)

type UserCreatedEventConsumer struct {
	eventListener     events.EventListener
	createUserUseCase *CreateUserUseCase
	logger            *zap.Logger
}

func (consumer *UserCreatedEventConsumer) Consume() error {
	eventMapChannel := make(chan map[string]interface{})
	err := consumer.eventListener.Listen(eventMapChannel)
	if err != nil {
		return err
	}
	go consumer.consumeEventMaps(eventMapChannel)
	return nil
}

func (consumer *UserCreatedEventConsumer) consumeEventMaps(eventMapChannel chan map[string]interface{}) {
	for eventMap := range eventMapChannel {
		err := consumer.handleEventMap(eventMap)
		if err != nil {
			consumer.logger.Warn(fmt.Sprintf("Error handling event: %s", err.Error()))
		}
	}
}

func (consumer *UserCreatedEventConsumer) handleEventMap(eventMap map[string]interface{}) error {
	event, err := user.UserCreatedEventFromMap(eventMap)
	if err != nil {
		return err
	}
	createUserRequest := CreateUserRequest{
		Email:     event.Email,
		Superuser: event.Superuser,
	}
	context := context.Background()
	if useCaseResponse := consumer.createUserUseCase.Execute(context, &createUserRequest); useCaseResponse.Err != nil {
		return useCaseResponse.Err
	}
	return nil
}

func NewUserCreatedEventConsumer(eventListenerFactory events.EventListenerFactory, useCase *CreateUserUseCase, logger *zap.Logger) *UserCreatedEventConsumer {
	eventName := reflect.TypeOf(user.UserCreatedEvent{}).Name()
	eventListener, err := eventListenerFactory.CreateListener(eventName)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating listener %s: %s", eventName, err.Error()))
	}
	return &UserCreatedEventConsumer{
		eventListener:     eventListener,
		createUserUseCase: useCase,
		logger:            logger,
	}
}
