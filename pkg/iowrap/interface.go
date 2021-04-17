package iowrap

import (
	"io"
)

// ReadAtCloser is the composition of io.Closer and io.ReaderAt
type ReadAtCloser interface {
	io.Closer
	io.ReaderAt
}
