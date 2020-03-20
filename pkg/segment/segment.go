package segment

import (
	"sync"
)

type Part interface {
	String() string

	Index() int
	Offset() int64

	Size() int64
}

// Func will handle a Segment.
type Func func(segment Segment)

type Segment interface {
	String() string

	ID() string
	Path() string

	InsertPart(int64, int64) (Part, error)
	Parts() []Part
}

// Segments carrys all segments in a service.
type Segments struct {
	s map[string]Segment
	l sync.RWMutex
}

// NewSegments will init a new segments.
func NewSegments() *Segments {
	return &Segments{
		s: make(map[string]Segment),
	}
}

// Insert will insert a segment into segments
func (s *Segments) Insert(seg Segment) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s[seg.ID()] = seg
}

// Get will get a segment from segments via id.
func (s *Segments) Get(id string) (seg Segment, err error) {
	s.l.RLock()
	defer s.l.RUnlock()

	seg, ok := s.s[id]
	if !ok {
		return nil, &Error{"get segment", ErrSegmentNotFound, &IndexBasedSegment{id: id}, nil}
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
