package types

import (
	"fmt"
)

// errorCode is the same from services.errorCode.
// We copy the code here to prevent depend cycle.
type errorCode struct {
	s string
}

func (e errorCode) Error() string {
	return e.s
}

// IsInternalError implements service.InternalError
func (e errorCode) IsInternalError() {}

var (
	// ErrNotImplemented will be returned while this operation is not
	// implemented by services.
	ErrNotImplemented = errorCode{"not implemented"}
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

func (oe OperationError) IsInternalError() {}

// NewOperationNotImplementedError will create a new NotImplemented error.
func NewOperationNotImplementedError(op string) error {
	return OperationError{
		op:  op,
		err: ErrNotImplemented,
	}
}
