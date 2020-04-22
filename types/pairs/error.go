package pairs

import (
	"errors"
	"fmt"
)

var (
	// ErrPairTypeMismatch means the pair's type is not match
	ErrPairTypeMismatch = errors.New("pair type mismatch")
)

// Error represents error related to a pair.
type Error struct {
	Op  string
	Err error

	Key   string
	Type  string
	Value interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: key %s, type %s, value %s: %s", e.Op, e.Key, e.Type, e.Value, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *Error) Unwrap() error {
	return e.Err
}
