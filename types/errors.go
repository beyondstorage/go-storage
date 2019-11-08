package types

import (
	"errors"
	"fmt"
)

// All possible error could be return by services.
var (
	// user handleable error
	ErrConfigIncorrect  = errors.New("config incorrect")
	ErrPermissionDenied = errors.New("permission denied")

	// caller handleable error
	ErrPairRequired    = errors.New("pair required")
	ErrObjectNotExist  = errors.New("object not exist")
	ErrDirAlreadyExist = errors.New("dir already exist")
	ErrDirNotEmpty     = errors.New("dir not empty")

	// unhandleable error
	ErrUnhandledError = errors.New("unhandled error")
)

// NewErrPairRequired will create a new pair required error.
func NewErrPairRequired(pair string) error {
	return fmt.Errorf("%s is required but missing: %w", pair, ErrPairRequired)
}
