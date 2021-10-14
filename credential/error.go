package credential

import (
	"errors"
	"fmt"
)

var (
	// ErrUnsupportedProtocol means protocol is unsupported
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
	// ErrInvalidValue means value is invalid.
	ErrInvalidValue = errors.New("invalid value")
)

// Error represents error related to credential.
type Error struct {
	Op  string
	Err error

	Protocol string
	Values   []string
}

func (e Error) Error() string {
	if e.Values == nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Protocol, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Protocol, e.Values, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e Error) Unwrap() error {
	return e.Err
}
