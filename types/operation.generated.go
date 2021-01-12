package types

import (
	"context"
	"io"
)

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

// IndexSegmenter is the interface for index based segment.
type IndexSegmenter interface {
	Segmenter

	// CompleteIndexSegment will complete a segment and merge them into a File.
	CompleteIndexSegment(seg Segment, parts []*Part, pairs ...Pair) (err error)
	// CompleteIndexSegmentWithContext will complete a segment and merge them into a File.
	CompleteIndexSegmentWithContext(ctx context.Context, seg Segment, parts []*Part, pairs ...Pair) (err error)

	// ListIndexSegment
	ListIndexSegment(seg Segment, pairs ...Pair) (pi *PartIterator, err error)
	// ListIndexSegmentWithContext
	ListIndexSegmentWithContext(ctx context.Context, seg Segment, pairs ...Pair) (pi *PartIterator, err error)

	// WriteIndexSegment will write a part into an index based segment.
	WriteIndexSegment(seg Segment, r io.Reader, index int, size int64, pairs ...Pair) (err error)
	// WriteIndexSegmentWithContext will write a part into an index based segment.
	WriteIndexSegmentWithContext(ctx context.Context, seg Segment, r io.Reader, index int, size int64, pairs ...Pair) (err error)
}

// Mover is the interface for Move.
type Mover interface {

	// Move will move an object in the service.
	Move(src string, dst string, pairs ...Pair) (err error)
	// MoveWithContext will move an object in the service.
	MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)
}

// OffsetSegmenter is the interface for offset based segment.
type OffsetSegmenter interface {
	Segmenter

	// WriteOffsetSegment will write a part into an index based segment.
	WriteOffsetSegment(seg Segment, r io.Reader, offset int64, size int64, pairs ...Pair) (err error)
	// WriteOffsetSegmentWithContext will write a part into an index based segment.
	WriteOffsetSegmentWithContext(ctx context.Context, seg Segment, r io.Reader, offset int64, size int64, pairs ...Pair) (err error)
}

// Reacher is the interface for Reach.
type Reacher interface {

	// Reach will provide a way, which can reach the object.
	Reach(path string, pairs ...Pair) (url string, err error)
	// ReachWithContext will provide a way, which can reach the object.
	ReachWithContext(ctx context.Context, path string, pairs ...Pair) (url string, err error)
}

// Segmenter
type Segmenter interface {

	// AbortSegment will abort a segment.
	AbortSegment(seg Segment, pairs ...Pair) (err error)
	// AbortSegmentWithContext will abort a segment.
	AbortSegmentWithContext(ctx context.Context, seg Segment, pairs ...Pair) (err error)

	// InitSegment will init a segment.
	InitSegment(path string, pairs ...Pair) (seg Segment, err error)
	// InitSegmentWithContext will init a segment.
	InitSegmentWithContext(ctx context.Context, path string, pairs ...Pair) (seg Segment, err error)

	// ListSegments will list segments.
	ListSegments(path string, pairs ...Pair) (si *SegmentIterator, err error)
	// ListSegmentsWithContext will list segments.
	ListSegmentsWithContext(ctx context.Context, path string, pairs ...Pair) (si *SegmentIterator, err error)
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
	Write(path string, r io.Reader, pairs ...Pair) (n int64, err error)
	// WriteWithContext will write data into a file.
	WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...Pair) (n int64, err error)
}

type PairPolicy struct {
	All bool

	// pairs for interface Copier
	Copy bool

	// pairs for interface Fetcher
	Fetch bool

	// pairs for interface IndexSegmenter
	CompleteIndexSegment bool
	ListIndexSegment     bool
	WriteIndexSegment    bool

	// pairs for interface Mover
	Move bool

	// pairs for interface OffsetSegmenter
	WriteOffsetSegment bool

	// pairs for interface Reacher
	Reach bool

	// pairs for interface Segmenter
	AbortSegment bool
	InitSegment  bool
	ListSegments bool

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
	WriteSize             bool
	WriteOffset           bool
	WriteStorageClass     bool
	WriteContentType      bool
	WriteContentMd5       bool
	WriteReadCallbackFunc bool
}
