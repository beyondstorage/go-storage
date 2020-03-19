package segment

import (
	"fmt"
	"sort"
	"sync"
)

// Part is a part of segment.
type Part struct {
	Index int
	Size  int64
}

func (p *Part) String() string {
	return fmt.Sprintf("Part {Index: %d, Size: %d}", p.Index, p.Size)
}

// Segment will hold the whole segment operations.
type Segment struct {
	ID   string
	Path string

	p map[int]*Part
	l sync.RWMutex
}

// Func will handle a Segment.
type Func func(segment *Segment)

// NewSegment will init a new segment.
func NewSegment(path, id string) *Segment {
	return &Segment{
		ID:   id,
		Path: path,

		p: make(map[int]*Part),
	}
}

func (s *Segment) String() string {
	return fmt.Sprintf(
		"Segment {ID: %s, Path: %s}", s.ID, s.Path,
	)
}

// InsertPart will insert a part into a segment and return it's Index.
// Index will start from 0.
func (s *Segment) InsertPart(idx int, size int64) (p *Part, err error) {
	if size == 0 {
		panic(&Error{"insert part", ErrPartSizeInvalid, s, nil})
	}

	s.l.Lock()
	defer s.l.Unlock()

	p = &Part{
		Size:  size,
		Index: idx,
	}

	s.p[idx] = p
	return p, nil
}

// Parts will return sorted p.
func (s *Segment) Parts() []*Part {
	s.l.RLock()
	defer s.l.RUnlock()

	x := make([]*Part, 0, len(s.p))
	for _, v := range s.p {
		v := v
		x = append(x, v)
	}
	sort.Slice(x, func(i, j int) bool { return x[i].Index < x[j].Index })
	return x
}

// Segments carrys all segments in a service.
type Segments struct {
	s map[string]*Segment
	l sync.RWMutex
}

// NewSegments will init a new segments.
func NewSegments() *Segments {
	return &Segments{
		s: make(map[string]*Segment),
	}
}

// Insert will insert a segment into segments
func (s *Segments) Insert(seg *Segment) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s[seg.ID] = seg
}

// Get will get a segment from segments via id.
func (s *Segments) Get(id string) (seg *Segment, err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	seg, ok := s.s[id]
	if !ok {
		return nil, &Error{"get segment", ErrSegmentNotFound, &Segment{ID: id}, nil}
	}
	return seg, nil
}

// Delete will delete a segments.
func (s *Segments) Delete(id string) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.s, id)
}

// Len will return length of segments.
func (s *Segments) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.s)
}
