package types

import (
	"context"
	"io"
)

// Appender
type Appender interface {

	// CreateAppend
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	// CreateAppendWithContext
	CreateAppendWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WriteAppend
	WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteAppendWithContext
	WriteAppendWithContext(ctx context.Context, o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
}

// Blocker
type Blocker interface {

	// CombineBlock
	CombineBlock(o *Object, bids []string, pairs ...Pair) (err error)
	// CombineBlockWithContext
	CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error)

	// CreateBlock
	CreateBlock(path string, pairs ...Pair) (o *Object, err error)
	// CreateBlockWithContext
	CreateBlockWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListBlock
	ListBlock(o *Object, pairs ...Pair) (bi *BlockIterator, err error)
	// ListBlockWithContext
	ListBlockWithContext(ctx context.Context, o *Object, pairs ...Pair) (bi *BlockIterator, err error)

	// WriteBlock
	WriteBlock(o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)
	// WriteBlockWithContext
	WriteBlockWithContext(ctx context.Context, o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)
}

// Copier is the interface for Copy.
type Copier interface {

	// Copy will copy an Object or multiple object in the service.
	Copy(src string, dst string, pairs ...Pair) (err error)
	// CopyWithContext will copy an Object or multiple object in the service.
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)
}

// Fetcher is the interface for Fetch.
type Fetcher interface {

	// Fetch will fetch from a given url to path.
	Fetch(path string, url string, pairs ...Pair) (err error)
	// FetchWithContext will fetch from a given url to path.
	FetchWithContext(ctx context.Context, path string, url string, pairs ...Pair) (err error)
}

// Mover is the interface for Move.
type Mover interface {

	// Move will move an object in the service.
	Move(src string, dst string, pairs ...Pair) (err error)
	// MoveWithContext will move an object in the service.
	MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)
}

// Pager
type Pager interface {

	// CreatePage
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	// CreatePageWithContext
	CreatePageWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WritePage
	WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
	// WritePageWithContext
	WritePageWithContext(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
}

// Parter
type Parter interface {

	// CompletePart
	CompletePart(o *Object, parts []*Part, pairs ...Pair) (err error)
	// CompletePartWithContext
	CompletePartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error)

	// CreatePart
	CreatePart(path string, pairs ...Pair) (o *Object, err error)
	// CreatePartWithContext
	CreatePartWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListPart
	ListPart(o *Object, pairs ...Pair) (pi *PartIterator, err error)
	// ListPartWithContext
	ListPartWithContext(ctx context.Context, o *Object, pairs ...Pair) (pi *PartIterator, err error)

	// WritePart
	WritePart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, err error)
	// WritePartWithContext
	WritePartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, err error)
}

// Reacher is the interface for Reach.
type Reacher interface {

	// Reach will provide a way, which can reach the object.
	Reach(path string, pairs ...Pair) (url string, err error)
	// ReachWithContext will provide a way, which can reach the object.
	ReachWithContext(ctx context.Context, path string, pairs ...Pair) (url string, err error)
}

// Servicer can maintain multipart storage services.
type Servicer interface {
	String() string

	// Create will create a new storager instance.
	Create(name string, pairs ...Pair) (store Storager, err error)
	// CreateWithContext will create a new storager instance.
	CreateWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error)

	// Delete will delete a storager instance.
	Delete(name string, pairs ...Pair) (err error)
	// DeleteWithContext will delete a storager instance.
	DeleteWithContext(ctx context.Context, name string, pairs ...Pair) (err error)

	// Get will get a valid storager instance for service.
	Get(name string, pairs ...Pair) (store Storager, err error)
	// GetWithContext will get a valid storager instance for service.
	GetWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error)

	// List will list all storager instances under this service.
	List(pairs ...Pair) (sti *StoragerIterator, err error)
	// ListWithContext will list all storager instances under this service.
	ListWithContext(ctx context.Context, pairs ...Pair) (sti *StoragerIterator, err error)
}

// Statistician is the interface for Statistical.
type Statistician interface {

	// Statistical will count service's statistics, such as Size, Count.
	Statistical(pairs ...Pair) (statistic *StorageStatistic, err error)
	// StatisticalWithContext will count service's statistics, such as Size, Count.
	StatisticalWithContext(ctx context.Context, pairs ...Pair) (statistic *StorageStatistic, err error)
}

// Storager is the interface for storage service.
type Storager interface {
	String() string

	// Delete will delete an Object from service.
	Delete(path string, pairs ...Pair) (err error)
	// DeleteWithContext will delete an Object from service.
	DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error)

	// List will return list a specific path.
	List(path string, pairs ...Pair) (oi *ObjectIterator, err error)
	// ListWithContext will return list a specific path.
	ListWithContext(ctx context.Context, path string, pairs ...Pair) (oi *ObjectIterator, err error)

	// Metadata will return current storager metadata.
	Metadata(pairs ...Pair) (meta *StorageMeta, err error)
	// MetadataWithContext will return current storager metadata.
	MetadataWithContext(ctx context.Context, pairs ...Pair) (meta *StorageMeta, err error)

	// Read will read the file's data.
	Read(path string, w io.Writer, pairs ...Pair) (n int64, err error)
	// ReadWithContext will read the file's data.
	ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...Pair) (n int64, err error)

	// Stat will stat a path to get info of an object.
	Stat(path string, pairs ...Pair) (o *Object, err error)
	// StatWithContext will stat a path to get info of an object.
	StatWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// Write will write data into a file.
	Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteWithContext will write data into a file.
	WriteWithContext(ctx context.Context, path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
}

type PairPolicy struct {
	All bool

	// pairs for interface Appender
	CreateAppend bool
	WriteAppend  bool

	// pairs for interface Blocker
	CombineBlock bool
	CreateBlock  bool
	ListBlock    bool
	WriteBlock   bool

	// pairs for interface Copier
	Copy bool

	// pairs for interface Fetcher
	Fetch bool

	// pairs for interface Mover
	Move bool

	// pairs for interface Pager
	CreatePage bool
	WritePage  bool

	// pairs for interface Parter
	CompletePart bool
	CreatePart   bool
	ListPart     bool
	WritePart    bool

	// pairs for interface Reacher
	Reach bool

	// pairs for interface Statistician
	Statistical bool

	// pairs for interface Storager
	Delete                bool
	List                  bool
	Metadata              bool
	Read                  bool
	ReadSize              bool
	ReadOffset            bool
	ReadReadCallbackFunc  bool
	Stat                  bool
	Write                 bool
	WriteStorageClass     bool
	WriteContentType      bool
	WriteContentMd5       bool
	WriteReadCallbackFunc bool
}
