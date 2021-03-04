package types

import (
	"context"
	"io"
)

// Appender is the interface for Append related operations.
type Appender interface {

	// CreateAppend will create an append object.
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	// CreateAppendWithContext will create an append object.
	CreateAppendWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WriteAppend will append content to an append object.
	WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteAppendWithContext will append content to an append object.
	WriteAppendWithContext(ctx context.Context, o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
}

// Blocker is the interface for Block related operations.
type Blocker interface {

	// CombineBlock will combine blocks into an object.
	CombineBlock(o *Object, bids []string, pairs ...Pair) (err error)
	// CombineBlockWithContext will combine blocks into an object.
	CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error)

	// CreateBlock will create a new block object.
	CreateBlock(path string, pairs ...Pair) (o *Object, err error)
	// CreateBlockWithContext will create a new block object.
	CreateBlockWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListBlock will list blocks belong to this object.
	ListBlock(o *Object, pairs ...Pair) (bi *BlockIterator, err error)
	// ListBlockWithContext will list blocks belong to this object.
	ListBlockWithContext(ctx context.Context, o *Object, pairs ...Pair) (bi *BlockIterator, err error)

	// WriteBlock will write content to a block.
	WriteBlock(o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)
	// WriteBlockWithContext will write content to a block.
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

// Multiparter is the interface for Multipart related operations.
type Multiparter interface {

	// CompleteMultipart will complete a multipart upload and construct an Object.
	CompleteMultipart(o *Object, parts []*Part, pairs ...Pair) (err error)
	// CompleteMultipartWithContext will complete a multipart upload and construct an Object.
	CompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error)

	// CreateMultipart will create a new multipart.
	CreateMultipart(path string, pairs ...Pair) (o *Object, err error)
	// CreateMultipartWithContext will create a new multipart.
	CreateMultipartWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListMultipart will list parts belong to this multipart.
	ListMultipart(o *Object, pairs ...Pair) (pi *PartIterator, err error)
	// ListMultipartWithContext will list parts belong to this multipart.
	ListMultipartWithContext(ctx context.Context, o *Object, pairs ...Pair) (pi *PartIterator, err error)

	// WriteMultipart will write content to a multipart.
	WriteMultipart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, err error)
	// WriteMultipartWithContext will write content to a multipart.
	WriteMultipartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, err error)
}

// Pager is the interface for Page related operations which support random write.
type Pager interface {

	// CreatePage will create a new page object.
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	// CreatePageWithContext will create a new page object.
	CreatePageWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WritePage will write content to specific offset.
	WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
	// WritePageWithContext will write content to specific offset.
	WritePageWithContext(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
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

// Storager is the interface for storage service.
type Storager interface {
	String() string

	// Create will create a new object without any api call.
	Create(path string, pairs ...Pair) (o *Object)

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

	// pairs for interface Multiparter
	CompleteMultipart bool
	CreateMultipart   bool
	ListMultipart     bool
	WriteMultipart    bool

	// pairs for interface Pager
	CreatePage bool
	WritePage  bool

	// pairs for interface Reacher
	Reach bool

	// pairs for interface Storager
	Create           bool
	Delete           bool
	List             bool
	ListListMode     bool
	Metadata         bool
	Read             bool
	ReadSize         bool
	ReadOffset       bool
	ReadIoCallback   bool
	Stat             bool
	Write            bool
	WriteContentType bool
	WriteContentMd5  bool
	WriteIoCallback  bool
}
