package types

import (
	"context"
	"io"

	"github.com/aos-dev/go-storage/v2/pkg/segment"
)

// Copier is the interface for Copy.
type Copier interface {

	// Copy will copy an Object or multiple object in the service.
	Copy(src string, dst string, pairs ...*Pair) (err error)
	// CopyWithContext will copy an Object or multiple object in the service.
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...*Pair) (err error)
}

// DirLister is used for directory based storage service to list objects under a dir.
type DirLister interface {

	// ListDir will return list a specific dir.
	ListDir(dir string, pairs ...*Pair) (oi *ObjectIterator, err error)
	// ListDirWithContext will return list a specific dir.
	ListDirWithContext(ctx context.Context, dir string, pairs ...*Pair) (oi *ObjectIterator, err error)
}

// DirSegmentsLister is used for directory based storage service to list segments under a dir.
type DirSegmentsLister interface {
	segmenter

	// ListDirSegments will list segments via dir.
	ListDirSegments(dir string, pairs ...*Pair) (si *segment.Iterator, err error)
	// ListDirSegmentsWithContext will list segments via dir.
	ListDirSegmentsWithContext(ctx context.Context, dir string, pairs ...*Pair) (si *segment.Iterator, err error)
}

// IndexSegmenter is the interface for index based segment.
type IndexSegmenter interface {
	segmenter

	// InitIndexSegment will init an index based segment.
	InitIndexSegment(path string, pairs ...*Pair) (seg segment.Segment, err error)
	// InitIndexSegmentWithContext will init an index based segment.
	InitIndexSegmentWithContext(ctx context.Context, path string, pairs ...*Pair) (seg segment.Segment, err error)

	// WriteIndexSegment will write a part into an index based segment.
	WriteIndexSegment(seg segment.Segment, r io.Reader, index int, size int64, pairs ...*Pair) (err error)
	// WriteIndexSegmentWithContext will write a part into an index based segment.
	WriteIndexSegmentWithContext(ctx context.Context, seg segment.Segment, r io.Reader, index int, size int64, pairs ...*Pair) (err error)
}

// Mover is the interface for Move.
type Mover interface {

	// Move will move an object in the service.
	Move(src string, dst string, pairs ...*Pair) (err error)
	// MoveWithContext will move an object in the service.
	MoveWithContext(ctx context.Context, src string, dst string, pairs ...*Pair) (err error)
}

// PrefixLister is used for prefix based storage service to list objects under a prefix.
type PrefixLister interface {

	// ListPrefix will return list a specific dir.
	ListPrefix(prefix string, pairs ...*Pair) (oi *ObjectIterator, err error)
	// ListPrefixWithContext will return list a specific dir.
	ListPrefixWithContext(ctx context.Context, prefix string, pairs ...*Pair) (oi *ObjectIterator, err error)
}

// PrefixSegmentsLister is used for prefix based storage service to list segments under a prefix.
type PrefixSegmentsLister interface {
	segmenter

	// ListPrefixSegments will list segments.
	ListPrefixSegments(prefix string, pairs ...*Pair) (si *segment.Iterator, err error)
	// ListPrefixSegmentsWithContext will list segments.
	ListPrefixSegmentsWithContext(ctx context.Context, prefix string, pairs ...*Pair) (si *segment.Iterator, err error)
}

// Reacher is the interface for Reach.
type Reacher interface {

	// Reach will provide a way, which can reach the object.
	Reach(path string, pairs ...*Pair) (url string, err error)
	// ReachWithContext will provide a way, which can reach the object.
	ReachWithContext(ctx context.Context, path string, pairs ...*Pair) (url string, err error)
}

// Segmenter
type segmenter interface {

	// AbortSegment will abort a segment.
	AbortSegment(seg segment.Segment, pairs ...*Pair) (err error)
	// AbortSegmentWithContext will abort a segment.
	AbortSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*Pair) (err error)

	// CompleteSegment will complete a segment and merge them into a File.
	CompleteSegment(seg segment.Segment, pairs ...*Pair) (err error)
	// CompleteSegmentWithContext will complete a segment and merge them into a File.
	CompleteSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*Pair) (err error)
}

// Servicer can maintain multipart storage services.
type Servicer interface {
	String() string

	// Create will create a new storager instance.
	Create(name string, pairs ...*Pair) (store Storager, err error)
	// CreateWithContext will create a new storager instance.
	CreateWithContext(ctx context.Context, name string, pairs ...*Pair) (store Storager, err error)

	// Delete will delete a storager instance.
	Delete(name string, pairs ...*Pair) (err error)
	// DeleteWithContext will delete a storager instance.
	DeleteWithContext(ctx context.Context, name string, pairs ...*Pair) (err error)

	// Get will get a valid storager instance for service.
	Get(name string, pairs ...*Pair) (store Storager, err error)
	// GetWithContext will get a valid storager instance for service.
	GetWithContext(ctx context.Context, name string, pairs ...*Pair) (store Storager, err error)

	// List will list all storager instances under this service.
	List(pairs ...*Pair) (err error)
	// ListWithContext will list all storager instances under this service.
	ListWithContext(ctx context.Context, pairs ...*Pair) (err error)
}

// Statistician is the interface for Statistical.
type Statistician interface {

	// Statistical will count service's statistics, such as Size, Count.
	Statistical(pairs ...*Pair) (statistic StorageStatistic, err error)
	// StatisticalWithContext will count service's statistics, such as Size, Count.
	StatisticalWithContext(ctx context.Context, pairs ...*Pair) (statistic StorageStatistic, err error)
}

// Storager is the interface for storage service.
type Storager interface {
	String() string

	// Delete will delete an Object from service.
	Delete(path string, pairs ...*Pair) (err error)
	// DeleteWithContext will delete an Object from service.
	DeleteWithContext(ctx context.Context, path string, pairs ...*Pair) (err error)

	// Metadata will return current storager's metadata.
	Metadata(pairs ...*Pair) (meta StorageMeta, err error)
	// MetadataWithContext will return current storager's metadata.
	MetadataWithContext(ctx context.Context, pairs ...*Pair) (meta StorageMeta, err error)

	// Read will read the file's data.
	Read(path string, w io.Writer, pairs ...*Pair) (err error)
	// ReadWithContext will read the file's data.
	ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...*Pair) (err error)

	// Stat will stat a path to get info of an object.
	Stat(path string, pairs ...*Pair) (o *Object, err error)
	// StatWithContext will stat a path to get info of an object.
	StatWithContext(ctx context.Context, path string, pairs ...*Pair) (o *Object, err error)

	// Write will write data into a file.
	Write(path string, r io.Reader, pairs ...*Pair) (err error)
	// WriteWithContext will write data into a file.
	WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...*Pair) (err error)
}
