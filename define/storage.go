package define

import (
	"io"
)

// Metadata is the metadata for storage service
type Metadata struct {
}

// File represents a seekable file or object.
type File struct {
}

// Stream represents a not seekable stream.
type Stream struct {
}

// Dir represents a virtual directory which contains files or streams.
type Dir struct {
}

// Storager is the interface for storage service.
type Storager interface {
	Capability() Capability

	Stat(path string, option ...Option)
	Delete(path string, option ...Option) (err error)
	Copy(src, dst string, option ...Option) (err error)
	Move(src, dst string, option ...Option) (err error)

	ListDir(path string, option ...Option) (dir chan Dir, file chan File, err error)

	ReadFile(path string, option ...Option) (r io.Reader, err error)
	WriteFile(path string, size int64, r io.Reader, option ...Option) (err error)

	ReadStream(path string, option ...Option) (r io.Reader, err error)
	WriteStream(path string, r io.Reader, option ...Option) (err error)

	InitSegment(path string, size int64, option ...Option) (err error)
	ReadSegment(path string, offset, size int64, option ...Option) (r io.Reader, err error)
	WriteSegment(path string, offset, size int64, r io.Reader, option ...Option) (err error)
	CompleteSegment(path string, option ...Option) (err error)
	AbortSegment(path string, option ...Option) (err error)
}
