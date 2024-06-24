package service

import (
	"errors"
	"fmt"
)

var (
	ErrRouteUrlNotFound          = errors.New("route URL not found")
	ErrResponseFromCalledService = errors.New("error response from called service")
	ErrAlreadyExists             = errors.New("already exists")
)

func calledServiceErrorResponse(statusCode int) error {
	return fmt.Errorf("%s: %d", ErrResponseFromCalledService, statusCode)
}
