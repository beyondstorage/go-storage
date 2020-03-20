package segment

import (
	"fmt"
	"sort"
	"sync"
)

// OffsetBasedPart is a part of segment.
type OffsetBasedPart struct {
	offset int64
	size   int64
}

func (p *OffsetBasedPart) String() string {
	return fmt.Sprintf("OffsetBasedPart {offset: %d, Size: %d}", p.offset, p.size)
}

func (p *OffsetBasedPart) Index() int {
	panic(&Error{Op: "index", Err: ErrPartTypeInvalid, Part: p})
}

func (p *OffsetBasedPart) Offset() int64 {
	return p.offset
}

func (p *OffsetBasedPart) Size() int64 {
	return p.size
}

// OffsetBasedSegment will hold the whole segment operations.
type OffsetBasedSegment struct {
	id   string
	path string

	p map[int64]*OffsetBasedPart
	l sync.RWMutex
}

// NewOffsetBasedSegment will init a new segment.
func NewOffsetBasedSegment(path, id string) *OffsetBasedSegment {
	return &OffsetBasedSegment{
		id:   id,
		path: path,

		p: make(map[int64]*OffsetBasedPart),
	}
}

func (s *OffsetBasedSegment) String() string {
	return fmt.Sprintf(
		"OffsetBasedSegment {ID: %s, Path: %s}", s.id, s.path,
	)
}

func (s *OffsetBasedSegment) ID() string {
	return s.id
}

func (s *OffsetBasedSegment) Path() string {
	return s.path
}

// InsertPart will insert a part into a segment and return it's offset.
// offset will start from 0.
func (s *OffsetBasedSegment) InsertPart(offset, size int64) (_ Part, err error) {
	if size == 0 {
		panic(&Error{"insert part", ErrPartSizeInvalid, s, nil})
	}

	s.l.Lock()
	defer s.l.Unlock()

	p := &OffsetBasedPart{
		offset: offset,
		size:   size,
	}

	s.p[p.offset] = p
	return p, nil
}

// Parts will return sorted p.
func (s *OffsetBasedSegment) Parts() []Part {
	s.l.RLock()
	defer s.l.RUnlock()

	x := make([]Part, 0, len(s.p))
	for _, v := range s.p {
		v := v
		x = append(x, v)
	}
	sort.Slice(x, func(i, j int) bool { return x[i].Offset() < x[j].Offset() })
	return x
}
