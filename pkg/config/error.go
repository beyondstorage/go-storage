package config

import (
	"errors"
)

var (
	// ErrInvalidConfig will be returned when config is invalid.
	ErrInvalidConfig = errors.New("invalid config")
)

// Error is the error returned by config
type Error struct {
	Op     string
	Config string

	Err error
}

func (e *Error) Error() string {
	return e.Op + " " + e.Config + ": " + e.Err.Error()
}

// Unwrap implements xerrors.Wrapper
func (e *Error) Unwrap() error {
	return e.Err
}
