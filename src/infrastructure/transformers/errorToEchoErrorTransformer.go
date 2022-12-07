package transformers

import (
	"go-as/src/domain/internals"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorToEchoErrorTransformer struct{}

func (transformer *ErrorToEchoErrorTransformer) Transform(err error) *echo.HTTPError {
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(transformer.getHTTPStatusCode(err), err.Error())
}

func (*ErrorToEchoErrorTransformer) getHTTPStatusCode(err error) int {
	switch err.(type) {
	case internals.UseCaseAuthorizationError:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func NewErrorToEchoErrorTransformer() *ErrorToEchoErrorTransformer {
	return &ErrorToEchoErrorTransformer{}
}
