package types

import (
	"errors"
	"fmt"
)

var (
	// ErrNotImplemented will be returned while this operation is not
	// implemented by services.
	ErrNotImplemented = errors.New("not implemented")
)

// OperationError is the error for operation related errors.
type OperationError struct {
	op  string
	err error
}

func (oe OperationError) Error() string {
	return fmt.Sprintf("operation %s: %v", oe.op, oe.err)
}

func (oe OperationError) Unwrap() error {
	return oe.err
}

// NewOperationNotImplementedError will create a new NotImplemented error.
func NewOperationNotImplementedError(op string) error {
	return OperationError{
		op:  op,
		err: ErrNotImplemented,
	}
}
