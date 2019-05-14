package sonic

import (
	"errors"
)

var (
	// ErrUnexpectedResponse occurs when the server sends an unexpected response
	ErrUnexpectedResponse = errors.New("Unexpected response from server")

	// ErrInvalidOptions occurs when invalid options are provided to a method
	ErrInvalidOptions = errors.New("invalid options")
)
