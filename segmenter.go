package storage

import (
	"context"
	"io"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

type segmenter interface {
	// CompleteSegment will complete a segment and merge them into a File.
	//
	// Implementer:
	//   - SHOULD return error while caller call CompleteSegment without init.
	// Caller:
	//   - SHOULD call InitIndexSegment before CompleteSegment.
	CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error)
	// CompleteSegmentWithContext will complete a segment and merge them into a File.
	CompleteSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*types.Pair) (err error)
	// AbortSegment will abort a segment.
	//
	// Implementer:
	//   - SHOULD return error while caller call AbortSegment without init.
	// Caller:
	//   - SHOULD call InitIndexSegment before AbortSegment.
	AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error)
	// AbortSegmentWithContext will abort a segment.
	AbortSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*types.Pair) (err error)
}

// DirSegmentsLister is used for directory based storage service to list segments under a dir.
type DirSegmentsLister interface {
	segmenter

	// ListDirSegments will list segments via dir.
	ListDirSegments(path string, pairs ...*types.Pair) (err error)
	// ListDirSegmentsWithContext will list segments via dir.
	ListDirSegmentsWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error)
}

// PrefixSegmentsLister is used for prefix based storage service to list segments under a prefix.
type PrefixSegmentsLister interface {
	segmenter

	// ListSegments will list segments.
	//
	// Implementer:
	//   - If prefix == "", services should return all segments.
	ListPrefixSegments(prefix string, pairs ...*types.Pair) (err error)
	// ListSegmentsWithContext will list segments.
	ListPrefixSegmentsWithContext(ctx context.Context, prefix string, pairs ...*types.Pair) (err error)
}

// IndexSegmenter is the interface for index based segment.
type IndexSegmenter interface {
	segmenter

	// InitIndexSegment will init an index based segment
	InitIndexSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error)
	// InitIndexSegmentWithContext will init an index based segment
	InitIndexSegmentWithContext(ctx context.Context, path string, pairs ...*types.Pair) (seg segment.Segment, err error)

	// WriteIndexSegment will write a part into an index based segment.
	WriteIndexSegment(seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error)
	// WriteIndexSegmentWithContext will write a part into an index based segment.
	WriteIndexSegmentWithContext(ctx context.Context, seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error)
}
