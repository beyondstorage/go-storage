// Package iowrap's Pipe Inspired by following project:
// - golang stdlib: io.Pipe and bytes.Buffer
// - https://github.com/djherbis/nio
package iowrap

import (
	"io"
	"sync"
)

// Pipe creates a synchronous in-memory pipe.
// It can be used to connect code expecting an io.Reader
// with code expecting an io.Writer.
//
// NOTES:
//   - PipeReader and PipeWriter is not thread safe
//   - Internal buffer will never be grow up, so write could be block while no space to write.
func Pipe() (r *PipeReader, w *PipeWriter) {
	bp := &bufpipe{
		buf:    make([]byte, 64*1024), // Set buffer to 64k
		length: 64 * 1024,
	}
	bp.rwait.L = &bp.l
	bp.wwait.L = &bp.l

	return &PipeReader{bp}, &PipeWriter{bp}
}

type bufpipe struct {
	buf    []byte
	length int // Buffer's length, which will not changed.
	size   int // Valid content size.
	offset int // Offset of consumed data.

	// rwait and wwait will reuse the global lock.
	l     sync.Mutex
	wwait sync.Cond
	rwait sync.Cond

	werr error //nolint:structcheck
	rerr error //nolint:structcheck
}

// Available space to write data.
func (p *bufpipe) gap() int {
	return p.length - p.size
}

// All data have been consumed.
func (p *bufpipe) empty() bool {
	return p.offset >= p.size
}

// Only size and offset need to be reset.
func (p *bufpipe) reset() {
	p.size = 0
	p.offset = 0
}

type PipeReader struct {
	*bufpipe
}

func (r *PipeReader) Read(p []byte) (n int, err error) {
	// Lock here to prevent concurrent read/write on buffer.
	r.l.Lock()
	// Send signal to wwait to allow next write.
	defer r.wwait.Signal()
	defer r.l.Unlock()

	for r.empty() {
		// Buffer is empty, reset to recover space.
		r.reset()

		if r.rerr != nil {
			return 0, io.ErrClosedPipe
		}

		if r.werr != nil {
			return 0, r.werr
		}

		// Buffer has consumed, allow next write.
		r.wwait.Signal()
		// Wait for read.
		r.rwait.Wait()
	}

	n = copy(p, r.buf[r.offset:r.size])
	r.offset += n
	return n, nil
}

func (r *PipeReader) Close() error {
	r.CloseWithError(nil)

	return nil
}

func (r *PipeReader) CloseWithError(err error) {
	if err == nil {
		err = io.ErrClosedPipe
	}

	r.l.Lock()
	defer r.l.Unlock()
	if r.rerr == nil {
		r.rerr = err
		r.rwait.Signal()
		r.wwait.Signal()
	}
}

type PipeWriter struct {
	*bufpipe
}

func (w *PipeWriter) Write(p []byte) (n int, err error) {
	var m int

	// Lock here to prevent concurrent read/write on buffer.
	w.l.Lock()
	// Send signal to rwait to allow next read.
	defer w.rwait.Signal()
	defer w.l.Unlock()

	l := len(p)

	for towrite := l; towrite > 0; towrite = l - n {
		for w.gap() == 0 {
			if w.rerr != nil {
				return 0, w.rerr
			}

			if w.werr != nil {
				return 0, io.ErrClosedPipe
			}

			// Buffer is full, allow next read.
			w.rwait.Signal()
			// Wait for write.
			w.wwait.Wait()
		}

		m = copy(w.buf, p[n:])
		w.size += m
		n += m
	}

	return
}

func (w *PipeWriter) Close() error {
	w.CloseWithError(nil)

	return nil
}

func (w *PipeWriter) CloseWithError(err error) {
	if err == nil {
		err = io.EOF
	}

	w.l.Lock()
	defer w.l.Unlock()
	if w.werr == nil {
		w.werr = err
		w.rwait.Signal()
		w.wwait.Signal()
	}
}
