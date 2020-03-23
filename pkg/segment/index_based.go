package segment

import (
	"fmt"
	"sort"
	"sync"
)

// IndexBasedPart is a part of segment.
type IndexBasedPart struct {
	Index int
	Size  int64
}

func (p *IndexBasedPart) String() string {
	return fmt.Sprintf("IndexBasedPart {Index: %d, Size: %d}", p.Index, p.Size)
}

// IndexBasedSegment will hold the whole segment operations.
type IndexBasedSegment struct {
	path string
	id   string

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

// ID implements Segment.ID()
func (s *IndexBasedSegment) ID() string {
	return s.id
}

// Path implements Segment.Path()
func (s *IndexBasedSegment) Path() string {
	return s.path
}

// InsertPart will insert a part into a segment and return it's Index.
// Index will start from 0.
func (s *IndexBasedSegment) InsertPart(index int, size int64) (p *IndexBasedPart, err error) {
	s.l.Lock()
	defer s.l.Unlock()

	p = &IndexBasedPart{
		Size:  size,
		Index: index,
	}

	s.p[p.Index] = p
	return p, nil
}

// Parts will return sorted p.
func (s *IndexBasedSegment) Parts() []*IndexBasedPart {
	s.l.RLock()
	defer s.l.RUnlock()

	x := make([]*IndexBasedPart, 0, len(s.p))
	for _, v := range s.p {
		v := v
		x = append(x, v)
	}
	sort.Slice(x, func(i, j int) bool { return x[i].Index < x[j].Index })
	return x
}
