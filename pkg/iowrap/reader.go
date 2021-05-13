package iowrap

import (
	"io"
)

// LimitReadCloser will return a limited hasCall closer.
func LimitReadCloser(r io.ReadCloser, n int64) *LimitedReadCloser {
	return &LimitedReadCloser{r, io.LimitReader(r, n)}
}

// LimitedReadCloser hasCall from underlying r and provide Close as well.
type LimitedReadCloser struct {
	c  io.Closer
	lr io.Reader
}

var (
	_ io.Reader = &LimitedReadCloser{}
	_ io.Closer = &LimitedReadCloser{}
)

// Read is copied from io.LimitedReader's Read.
func (l *LimitedReadCloser) Read(p []byte) (n int, err error) {
	return l.lr.Read(p)
}

// Close will close underlying reader.
func (l *LimitedReadCloser) Close() error {
	return l.c.Close()
}

// SectionReadCloser will return a sectioned hasCall closer.
func SectionReadCloser(r ReadAtCloser, off, n int64) *SectionedReadCloser {
	return &SectionedReadCloser{r, io.NewSectionReader(r, off, n)}
}

// SectionedReadCloser hasCall from underlying r and provide Close as well.
type SectionedReadCloser struct {
	c  io.Closer
	sr *io.SectionReader
}

var (
	_ io.Reader = &SectionedReadCloser{}
	_ io.Closer = &SectionedReadCloser{}
)

// Read is copied from io.SectionReader's Read.
func (s *SectionedReadCloser) Read(p []byte) (n int, err error) {
	return s.sr.Read(p)
}

// Close will close underlying reader.
func (s *SectionedReadCloser) Close() error {
	return s.c.Close()
}

// ReadSeekCloser wraps a io.Reader returning a SeekCloseableReader. Allows the
// SDK to accept an io.Reader that is not also an io.Seeker for unsigned
// streaming payload API operations.
//
// A ReadSeekCloser wrapping an nonseekable io.Reader used in an API
// operation's input will prevent that operation being retried in the case of
// network errors, and cause operation requests to fail if the operation
// requires payload signing.
//
// NOTES: Idea borrows from AWS Go SDK.
func ReadSeekCloser(r io.Reader) *SeekCloseableReader {
	return &SeekCloseableReader{r, 0}
}

// SizedReadSeekCloser will return a size featured SeekCloseableReader.
func SizedReadSeekCloser(r io.Reader, size int64) *SeekCloseableReader {
	return &SeekCloseableReader{r, size}
}

// SeekCloseableReader represents a reader that can also delegate io.Seeker and
// io.Closer interfaces to the underlying object if they are available.
type SeekCloseableReader struct {
	r    io.Reader
	size int64
}

var (
	_ io.Reader = &SeekCloseableReader{}
	_ io.Seeker = &SeekCloseableReader{}
	_ io.Closer = &SeekCloseableReader{}
)

// Read reads from the reader up to size of p. The number of bytes read, and
// error if it occurred will be returned.
//
// If the reader is not an io.Reader zero bytes read, and nil error will be
// returned.
//
// Performs the same functionality as io.Reader Read
func (r *SeekCloseableReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

// Seek sets the offset for the next Read to offset, interpreted according to
// whence: 0 means relative to the origin of the file, 1 means relative to the
// current offset, and 2 means relative to the end. Seek returns the new offset
// and an error, if any.
//
// If the SeekCloseableReader is not an io.Seeker nothing will be done to underlying Reader.
// For example: seek to end and then seek current will still return 0.
func (r *SeekCloseableReader) Seek(offset int64, whence int) (int64, error) {
	t, ok := r.r.(io.Seeker)
	if ok {
		return t.Seek(offset, whence)
	}

	// If underlying reader is not seekable, we will use internal size.
	switch whence {
	case io.SeekStart, io.SeekCurrent:
		return int64(0), nil
	case io.SeekEnd:
		return r.size, nil
	default:
		panic("invalid whence")
	}
}

// Close closes the SeekCloseableReader.
//
// If the SeekCloseableReader is not an io.Closer nothing will be done.
func (r *SeekCloseableReader) Close() error {
	t, ok := r.r.(io.Closer)
	if ok {
		return t.Close()
	}
	return nil
}

// CallbackReader will create a new CallbackifyReader.
func CallbackReader(r io.Reader, fn func([]byte)) *CallbackifyReader {
	return &CallbackifyReader{
		r:  r,
		fn: fn,
	}
}

// CallbackifyReader will execute callback func in Read.
type CallbackifyReader struct {
	r  io.Reader
	fn func([]byte)
}

var (
	_ io.Reader = &CallbackifyReader{}
)

// Read will read from underlying Reader.
func (r *CallbackifyReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.fn(p[:n])
	return n, err
}

// CallbackReadCloser will create a new CallbackifyReadCloser.
func CallbackReadCloser(r io.ReadCloser, fn func([]byte)) *CallbackifyReadCloser {
	return &CallbackifyReadCloser{
		r:  r,
		fn: fn,
	}
}

// CallbackifyReadCloser will execute callback func in Read.
type CallbackifyReadCloser struct {
	r  io.ReadCloser
	fn func([]byte)
}

var (
	_ io.Reader = &CallbackifyReadCloser{}
	_ io.Closer = &CallbackifyReadCloser{}
)

// Read will read from underlying Reader.
func (r *CallbackifyReadCloser) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.fn(p[:n])
	return n, err
}

// Close will close underlying Reader.
func (r *CallbackifyReadCloser) Close() error {
	return r.r.Close()
}
