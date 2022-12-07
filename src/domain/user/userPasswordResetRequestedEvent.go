package user

import "errors"

type UserCreatedEvent struct {
	Email     string
	Superuser bool
}

func UserCreatedEventFromMap(eventMap map[string]interface{}) (*UserCreatedEvent, error) {
	superuser, typeCheck := eventMap["Superuser"].(bool)
	if !typeCheck {
		return nil, errors.New("malformed data for creating UserCreatedEvent")
	}
	return &UserCreatedEvent{
		Email:     eventMap["Email"].(string),
		Superuser: superuser,
	}, nil
}
