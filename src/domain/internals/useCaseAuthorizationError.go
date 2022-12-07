package internals

import "fmt"

type UseCaseAuthorizationError struct {
	Email      string
	Permission string
}

func (err UseCaseAuthorizationError) Error() string {
	return fmt.Sprintf("User %s has not authorization for %s", err.Email, err.Permission)
}
