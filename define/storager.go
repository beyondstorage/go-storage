package define

import (
	"io"
)

// Storager is the interface for storage service.
//
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY",
// and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.
type Storager interface {
	// Caller SHOULD check capability before use file, stream or segment related operations.
	// This function will be generated via meta.json
	Capability() Capability
	// Caller SHOULD check option availability before use any option.
	// This function will be generated via meta.json
	IsOptionAvailable(action, option string) bool

	// Service MUST return a *File, *Stream or *Dir.
	Stat(path string, option ...Option) (i *Informer, err error)
	// Service MAY accept a recursive options to support delete recursively.
	Delete(path string, option ...Option) (err error)
	// Service MAY accept a recursive options to support copy recursively.
	Copy(src, dst string, option ...Option) (err error)
	// Service MAY accept a recursive options to support move recursively.
	Move(src, dst string, option ...Option) (err error)

	// Service MUST close all channel while no more items listed.
	// Service MAY start goroutines to transfer items.
	ListDir(path string, option ...Option) (dir chan *Dir, file chan *File, stream chan *Stream, err error)

	// Caller MUST close reader while error happened or all data read.
	ReadFile(path string, option ...Option) (r io.ReadCloser, err error)
	// Service MUST close reader while error happened or all data written.
	WriteFile(path string, size int64, r io.ReadCloser, option ...Option) (err error)

	// Caller MUST close reader while error happened or all data read.
	ReadStream(path string, option ...Option) (r io.ReadCloser, err error)
	// Service MUST close reader while error happened or all data written.
	WriteStream(path string, r io.ReadCloser, option ...Option) (err error)

	// Service MUST maintain whole segment operation runtime data, including upload_id and any other similar things,
	// caller will have no idea about those underlying implements.
	InitSegment(path string, size int64, option ...Option) (err error)
	// Caller SHOULD NOT call InitSegment before ReadSegment.
	ReadSegment(path string, offset, size int64, option ...Option) (r io.ReadCloser, err error)
	// Service SHOULD return error while caller call WriteSegment without init.
	// Caller SHOULD call InitSegment before WriteSegment.
	WriteSegment(path string, offset, size int64, r io.ReadCloser, option ...Option) (err error)
	// Service SHOULD return error while caller call CompleteSegment without init.
	// Caller SHOULD call InitSegment before CompleteSegment.
	CompleteSegment(path string, option ...Option) (err error)
	// Service SHOULD return error while caller call AbortSegment without init.
	// Caller SHOULD call InitSegment before AbortSegment.
	AbortSegment(path string, option ...Option) (err error)
}
