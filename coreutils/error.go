package coreutils

import (
	"errors"
	"fmt"

	"github.com/aos-dev/go-storage/v2/types"
)

var (
	// ErrServicerNotImplemented will return when service doesn't implement Servicer.
	ErrServicerNotImplemented = errors.New("servicer not implemented")
	// ErrStoragerNotImplemented will return when service doesn't implement Storager.
	ErrStoragerNotImplemented = errors.New("storager not implemented")
)

// OpenError returned while open related error happens.
type OpenError struct {
	Err error

	Type  string
	Pairs []*types.Pair
}

func (e *OpenError) Error() string {
	return fmt.Sprintf("open: %s, %v: %s", e.Type, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *OpenError) Unwrap() error {
	return e.Err
}
