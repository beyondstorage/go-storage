package qingstor

import (
	"errors"
)

var (
	// ErrInvalidBucketName will be returned while bucket name is invalid.
	ErrInvalidBucketName = errors.New("invalid bucket name")
	// ErrInvalidWorkDir will be returned while work dir is invalid.
	ErrInvalidWorkDir = errors.New("invalid work dir, must start and end with only one '/'")
)
