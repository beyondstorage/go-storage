package coreutils

import (
	"fmt"

	"github.com/Xuanwo/storage/types"
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
