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
	*Segment

	p map[int]*IndexBasedPart
	l sync.RWMutex
}

// NewIndexBasedSegment will init a new segment.
func NewIndexBasedSegment(path, id string) *IndexBasedSegment {
	return &IndexBasedSegment{
		Segment: &Segment{
			ID:   id,
			Path: path,
		},

		p: make(map[int]*IndexBasedPart),
	}
}

func (s *IndexBasedSegment) String() string {
	return fmt.Sprintf(
		"IndexBasedSegment {ID: %s, Path: %s}", s.ID, s.Path,
	)
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

// IndexBasedSegments carrys all segments in a service.
type IndexBasedSegments struct {
	s map[string]*IndexBasedSegment
	l sync.RWMutex
}

// NewIndexBasedSegments will init a new segments.
func NewIndexBasedSegments() *IndexBasedSegments {
	return &IndexBasedSegments{
		s: make(map[string]*IndexBasedSegment),
	}
}

// Insert will insert a segment into segments
func (s *IndexBasedSegments) Insert(seg *IndexBasedSegment) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s[seg.ID] = seg
}

// Get will get a segment from segments via id.
func (s *IndexBasedSegments) Get(id string) (seg *IndexBasedSegment, err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	seg, ok := s.s[id]
	if !ok {
		return nil, &Error{"get segment", ErrSegmentNotFound, &Segment{ID: id}}
	}
	return seg, nil
}

// Delete will delete a segments.
func (s *IndexBasedSegments) Delete(id string) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.s, id)
}

// Len will return length of segments.
func (s *IndexBasedSegments) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.s)
}
