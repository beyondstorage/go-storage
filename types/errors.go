package types

import (
	"errors"
)

// All possible error could be return by services.
var (
	// handleable error
	ErrConfigIncorrect  = errors.New("config incorrect")
	ErrPermissionDenied = errors.New("permission denied")

	// unhandleable but information available
	// Object here could be Dir, File or Stream.
	ErrObjectNotExist  = errors.New("object not exist")
	ErrDirAlreadyExist = errors.New("dir already exist")
	ErrDirNotEmpty     = errors.New("dir not empty")

	// unhandleable error
	ErrUnhandledError = errors.New("unhandled error")
)
