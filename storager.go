package storage

import (
	"context"
	"io"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// StoragerFunc will handle a storager.
type StoragerFunc func(Storager)

/*
Storager is the interface for storage service.

Currently, we support two different types of storage services: prefix based and directory based. Prefix based storage
service is usually an object storage service, such as AWS; And directory based service is often a POSIX file system.
We used to treat them as different abstract level services, but in this project, we will unify both of them to make a
production ready high performance vendor lock free storage layer.

Every storager will implement the same interface but with different capability and operation pairs set.

Everything in a storager is an Object with two types: File, Dir.
File is the smallest unit in service, it will have content and metadata. Dir is a container for File and Dir.
In prefix-based storage service, Dir is usually an empty key end with "/" or with special content type.
For directory-based service, Dir will be corresponded to the real directory on file system.

In the comments of every method, we will use following rules to standardize the Storager's behavior:

  - The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY",
    and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.
  - Implementer is the provider of the service, while trying to implement Storager interface, you need to follow.
  - Caller is the user of the service, while trying to use the Storager interface, you need to follow.
*/
type Storager interface {
	// String will implement Stringer.
	String() string

	// Metadata will return current storager's metadata.
	//
	// Implementer:
	//   - Metadata SHOULD only return static data without API call or with a cache.
	// Caller:
	//   - Metadata SHOULD be cheap.
	Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error)
	// MetadataWithContext will return current storager's metadata.
	MetadataWithContext(ctx context.Context, pairs ...*types.Pair) (m metadata.StorageMeta, err error)

	// Read will read the file's data.
	//
	// Caller:
	//   - MUST close reader while error happened or all data read.
	Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error)
	// ReadWithContext will read the file's data.
	ReadWithContext(ctx context.Context, path string, pairs ...*types.Pair) (r io.ReadCloser, err error)
	// Write will write data into a file.
	//
	// Caller:
	//   - MUST close reader while error happened or all data written.
	Write(path string, r io.Reader, pairs ...*types.Pair) (err error)
	// WriteWithContext will write data into a file.
	WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...*types.Pair) (err error)
	// Stat will stat a path to get info of an object.
	Stat(path string, pairs ...*types.Pair) (o *types.Object, err error)
	// StatWithContext will stat a path to get info of an object.
	StatWithContext(ctx context.Context, path string, pairs ...*types.Pair) (o *types.Object, err error)
	// Delete will delete an Object from service.
	Delete(path string, pairs ...*types.Pair) (err error)
	// DeleteWithContext will delete an Object from service.
	DeleteWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error)
}

// DirLister is used for directory based storage service to list objects under a dir.
type DirLister interface {
	// ListDir will return list a specific dir.
	ListDir(path string, pairs ...*types.Pair) (err error)
	// ListDirWithContext will return list a specific path.
	ListDirWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error)
}

// PrefixLister is used for prefix based storage service to list objects under a prefix.
type PrefixLister interface {
	// ListPrefix will return list a specific dir.
	//
	// Caller:
	//   - prefix SHOULD NOT start with /, and SHOULD relative to workdir.
	ListPrefix(prefix string, pairs ...*types.Pair) (err error)
	// ListPrefixWithContext will return list a specific path.
	ListPrefixWithContext(ctx context.Context, prefix string, pairs ...*types.Pair) (err error)
}

// Copier is the interface for Copy.
type Copier interface {
	// Copy will copy an Object or multiple object in the service.
	Copy(src, dst string, pairs ...*types.Pair) (err error)
	// CopyWithContext will copy an Object or multiple object in the service.
	CopyWithContext(ctx context.Context, src, dst string, pairs ...*types.Pair) (err error)
}

// Mover is the interface for Move.
type Mover interface {
	// Move will move an object or multiple object in the service.
	Move(src, dst string, pairs ...*types.Pair) (err error)
	// MoveWithContext will move an object or multiple object in the service.
	MoveWithContext(ctx context.Context, src, dst string, pairs ...*types.Pair) (err error)
}

// Reacher is the interface for Reach.
type Reacher interface {
	// Reach will provide a way, which can reach the object.
	//
	// Implementer:
	//   - SHOULD return a publicly reachable http url.
	Reach(path string, pairs ...*types.Pair) (url string, err error)
	// ReachWithContext will provide a way, which can reach the object.
	ReachWithContext(ctx context.Context, path string, pairs ...*types.Pair) (url string, err error)
}

// Statistician is the interface for Statistical.
type Statistician interface {
	// Statistical will count service's statistics, such as Size, Count.
	//
	// Implementer:
	//   - Statistical SHOULD only return dynamic data like Size, Count.
	// Caller:
	//   - Statistical call COULD be expensive.
	Statistical(pairs ...*types.Pair) (metadata.StorageStatistic, error)
	// StatisticalWithContext will count service's statistics, such as Size, Count.
	StatisticalWithContext(ctx context.Context, pairs ...*types.Pair) (metadata.StorageStatistic, error)
}
