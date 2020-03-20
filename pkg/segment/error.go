package segment

import (
	"errors"
	"fmt"
)

// All errors that segment could return.
var (
	ErrSegmentNotFound = errors.New("segment not found")
)

// Error represents error related to segment.
type Error struct {
	Op  string
	Err error

	Segment *Segment
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Op, e.Segment, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *Error) Unwrap() error {
	return e.Err
}
