package iowrap

import (
	"io"
)

// LimitReadCloser will return a limited hasCall closer.
func LimitReadCloser(r io.ReadCloser, n int64) *LimitedReadCloser {
	return &LimitedReadCloser{r, io.LimitReader(r, n)}
}

// LimitedReadCloser hasCall from underlying r and provide Close as well.
//go:generate mockgen -package iowrap -destination mock_test.go io Reader,Closer,ReaderAt
type LimitedReadCloser struct {
	r  io.ReadCloser
	lr io.Reader
}

// Read is copied from io.LimitedReader's Read.
func (l *LimitedReadCloser) Read(p []byte) (n int, err error) {
	return l.lr.Read(p)
}

// Close will close underlying reader.
func (l *LimitedReadCloser) Close() error {
	return l.r.Close()
}

// SectionReadCloser will return a sectioned hasCall closer.
func SectionReadCloser(r interface {
	io.Closer
	io.ReaderAt
}, off, n int64) *SectionedReadCloser {
	return &SectionedReadCloser{r, io.NewSectionReader(r, off, n)}
}

// SectionedReadCloser hasCall from underlying r and provide Close as well.
type SectionedReadCloser struct {
	r interface {
		io.Closer
		io.ReaderAt
	}
	sr *io.SectionReader
}

// Read is copied from io.SectionReader's Read.
func (s *SectionedReadCloser) Read(p []byte) (n int, err error) {
	return s.sr.Read(p)
}

// Close will close underlying reader.
func (s *SectionedReadCloser) Close() error {
	return s.r.Close()
}
