package qingstor

import (
	"errors"
)

var (
	// ErrInvalidBucketName will be returned while bucket name is invalid.
	ErrInvalidBucketName = errors.New("invalid bucket name")
)
