package segment

import (
	"fmt"
	"sort"
	"sync"
)

// IndexBasedPart is a part of segment.
type IndexBasedPart struct {
	index int
	size  int64
}

func (p *IndexBasedPart) String() string {
	return fmt.Sprintf("IndexBasedPart {Index: %d, Size: %d}", p.index, p.size)
}

func (p *IndexBasedPart) Index() int {
	return p.index
}

func (p *IndexBasedPart) Offset() int64 {
	panic(&Error{Op: "offset", Err: ErrPartTypeInvalid, Part: p})
}

func (p *IndexBasedPart) Size() int64 {
	return p.size
}

// IndexBasedSegment will hold the whole segment operations.
type IndexBasedSegment struct {
	id   string
	path string

	p map[int]*IndexBasedPart
	l sync.RWMutex
}

// NewIndexBasedSegment will init a new segment.
func NewIndexBasedSegment(path, id string) *IndexBasedSegment {
	return &IndexBasedSegment{
		id:   id,
		path: path,

		p: make(map[int]*IndexBasedPart),
	}
}

func (s *IndexBasedSegment) String() string {
	return fmt.Sprintf(
		"IndexBasedSegment {ID: %s, Path: %s}", s.id, s.path,
	)
}

func (s *IndexBasedSegment) ID() string {
	return s.id
}

func (s *IndexBasedSegment) Path() string {
	return s.path
}

// InsertPart will insert a part into a segment and return it's index.
// index will start from 0.
func (s *IndexBasedSegment) InsertPart(index, size int64) (_ Part, err error) {
	if size == 0 {
		panic(&Error{"insert part", ErrPartSizeInvalid, s, nil})
	}

	s.l.Lock()
	defer s.l.Unlock()

	p := &IndexBasedPart{
		size:  size,
		index: int(index),
	}

	s.p[p.index] = p
	return p, nil
}

// Parts will return sorted p.
func (s *IndexBasedSegment) Parts() []Part {
	s.l.RLock()
	defer s.l.RUnlock()

	x := make([]Part, 0, len(s.p))
	for _, v := range s.p {
		v := v
		x = append(x, v)
	}
	sort.Slice(x, func(i, j int) bool { return x[i].Index() < x[j].Index() })
	return x
}
