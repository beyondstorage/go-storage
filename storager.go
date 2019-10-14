package storage

import (
	"io"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
)

// All actions that storager used.
const (
	ActionStat   = "stat"
	ActionDelete = "delete"
	ActionCopy   = "copy"
	ActionMove   = "move"

	ActionReach = "reach"

	ActionCreateDir = "create_dir"
	ActionListDir   = "list_dir"

	ActionRead        = "read"
	ActionWriteFile   = "write_file"
	ActionWriteStream = "write_stream"

	ActionInitSegment     = "init_segment"
	ActionReadSegment     = "read_segment"
	ActionWriteSegment    = "write_segment"
	ActionCompleteSegment = "complete_segment"
	ActionAbortSegment    = "abort_segment"
)

/*
Servicer can maintain multipart storage services.

Implementer can choose to implement this interface or not.
*/
type Servicer interface {
	// Init will init service itself.
	Init(pairs ...*types.Pair) (err error)

	// List will list all storager instances under this service.
	List(pairs ...*types.Pair) ([]Storager, error)
	// Get will get a valid storager instance for service.
	Get(name string, pairs ...*types.Pair) (Storager, error)
	// Create will create a new storager instance.
	Create(name string, pairs ...*types.Pair) (Storager, error)
	// Delete will delete a storager instance.
	Delete(name string, pairs ...*types.Pair) (err error)
}

/*
Storager is the interface for storage service.

Currently we support two different type storage services: prefix based and directory based. Prefix based storage
service is usually a object storage service, such as AWS; And directory based service is often a POSIX file system.
We used to treat them as different abstract level services, but in this project, we will unify both of them to make a
production ready high performance vendor lock free storage layer.

Every service will implement the same interface but with different capability and operation pairs set.

Everything in a service is an Object with three types: File, Stream, Dir.
Both File and Stream are smallest unit in service, they will have content and metadata. The difference is File has
determined size but Stream's size is uncertain. Dir is a container for File, Stream and Dir. In prefix based storage
service, Dir is usually an empty key end with "/" or with special content type. And for directory based service, Dir
will be corresponded to the real directory on file system.

In the comments of every method, we will use following rules to standardize the Storager's behavior:

  - The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY",
    and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.
  - Implementer is the provider of the service, while you trying to implement Storager interface, you need to follow.
  - Caller is the user of the service, while you trying to use the Storager interface, you need to follow.
*/
type Storager interface {
	// Meta Operations.
	// They will be generated via meta.json, don't need to implement manually.

	// Capability will return the capability of current storage services.
	//
	// Caller:
	//   - SHOULD check capability before use file, stream or segment related operations.
	Capability() types.Capability
	// IsPairAvailable will check whether current service support this pair.
	//
	// Caller:
	//   - SHOULD check pairs availability before use any pairs.
	IsPairAvailable(action, key string) bool

	// Base Operations.

	// Metadata will return current storager's metadata.
	//
	// Implementer:
	//   - MAY return following data: Name,Location, Size, Count
	Metadata() (types.Metadata, error)
	// Stat will stat a path to get info of an object.
	//
	// Implementer:
	//   - MUST fill object's name and type.
	Stat(path string, pairs ...*types.Pair) (o *types.Object, err error)
	// Delete will delete an Object or multiple object from service.
	// path could be File, Stream or Dir
	//
	// Implementer:
	//   - MAY accept a recursive pair to support delete Dir recursively.
	Delete(path string, pairs ...*types.Pair) (err error)
	// Copy will copy an Object or multiple object in the service.
	//
	// Implementer:
	//   - MAY accept a recursive pairs to support copy recursively.
	Copy(src, dst string, pairs ...*types.Pair) (err error)
	// Move will move an object or multiple object in the service.
	//
	// Implementer:
	//   - MAY accept a recursive pairs to support move recursively.
	Move(src, dst string, pairs ...*types.Pair) (err error)

	// Fetch Operation.
	// Fetch(src, dst string) ?

	// Reach will provide a way which can reach the object.
	//
	// Implementer:
	//   - SHOULD return a publicly reachable http url.
	Reach(path string, pairs ...*types.Pair) (url string, err error)

	// Dir Operations.

	// CreateDir will create a Dir in the services.
	CreateDir(path string, pairs ...*types.Pair) (err error)
	// ListDir will return an Iterator which can list all object under the Dir.
	ListDir(path string, pairs ...*types.Pair) iterator.Iterator

	// File && Stream Operations.

	// Read will read the file's data.
	//
	// Caller:
	//   - MUST close reader while error happened or all data read.
	Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error)

	// File Only Operations.

	// WriteFile will write data into file.
	//
	// Caller:
	//   - MUST close reader while error happened or all data written.
	WriteFile(path string, size int64, r io.Reader, pairs ...*types.Pair) (err error)

	// Stream Operations.

	// WriteStream will write data into file.
	//
	// Caller:
	//   - MUST close reader while error happened or all data written.
	WriteStream(path string, r io.Reader, pairs ...*types.Pair) (err error)

	// Segment Operations.

	// InitSegment will init a segment which could be a File after complete.
	//
	// Implementer:
	//   - MUST maintain whole segment operation runtime data, including upload_id and any other similar things.
	// Caller:
	//   - SHOULD call InitSegment before Write, Complete or Abort.
	InitSegment(path string, pairs ...*types.Pair) (err error)
	// ReadSegment will read a segment from an Object.
	//
	// Implementer:
	//   - MAY support ReadSegment for Stream is it's seekable.
	// Caller:
	//   - SHOULD NOT call InitSegment before ReadSegment.
	ReadSegment(path string, offset, size int64, pairs ...*types.Pair) (r io.ReadCloser, err error)
	// WriteSegment will read data into segment.
	//
	// Implementer:
	//   - SHOULD return error while caller call WriteSegment without init.
	// Caller:
	//   - SHOULD call InitSegment before WriteSegment.
	WriteSegment(path string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error)
	// CompleteSegment will complete a segment and merge them into a File.
	//
	// Implementer:
	//   - SHOULD return error while caller call CompleteSegment without init.
	// Caller:
	//   - SHOULD call InitSegment before CompleteSegment.
	CompleteSegment(path string, pairs ...*types.Pair) (err error)
	// AbortSegment will abort a segment.
	//
	// Implementer:
	//   - SHOULD return error while caller call AbortSegment without init.
	// Caller:
	//   - SHOULD call InitSegment before AbortSegment.
	AbortSegment(path string, pairs ...*types.Pair) (err error)
}
