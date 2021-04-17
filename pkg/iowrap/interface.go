package iowrap

import (
	"context"
	"io"
)

// CancelableReader allow user to cancel a read call.
type CancelableReader interface {
	ReadWithContext(ctx context.Context, p []byte) (n int, err error)
}

// CancelableReadFrom allow user to cancel a read_from call.
type CancelableReadFrom interface {
	ReadFromWithContext(ctx context.Context, r io.Reader) (n int64, err error)
}

// CancelableWriter allow user to cancel a write call.
type CancelableWriter interface {
	WriteWithContext(ctx context.Context, p []byte) (n int, err error)
}

// CancelableWriteTo allow user to cancel a write_to call.
type CancelableWriteTo interface {
	WriteToWithContext(ctx context.Context, w io.Writer) (n int64, err error)
}

// CancelableCloser allow user to cancel a close call.
type CancelableCloser interface {
	CloseWithContext(ctx context.Context) error
}

// CancelableWriteCloser allow user to cancel a write or close call.
type CancelableWriteCloser interface {
	CancelableWriter
	CancelableCloser
}

// ReadAtCloser is the composition of io.Closer and io.ReaderAt
type ReadAtCloser interface {
	io.Closer
	io.ReaderAt
}
