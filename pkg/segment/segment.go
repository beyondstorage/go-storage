package segment

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// All errors that segment could return.
var (
	ErrPartSizeInvalid         = errors.New("part size invalid")
	ErrPartIntersected         = errors.New("part intersected")
	ErrSegmentAlreadyInitiated = errors.New("segment already initiated")
	ErrSegmentNotInitiated     = errors.New("segment not initiated")
	ErrSegmentPartsEmpty       = errors.New("segment Parts are empty")
	ErrSegmentNotFulfilled     = errors.New("segment not fulfilled")
)

// Part is a part of segment.
type Part struct {
	Offset int64
	Size   int64

	Index int
}

func (p *Part) String() string {
	return fmt.Sprintf("Part {Offset: %d, Size: %d}", p.Offset, p.Size)
}

// Segment will hold the whole segment operations.
type Segment struct {
	ID       string
	Path     string
	PartSize int64
	Parts    map[int64]*Part

	l sync.RWMutex
}

// NewSegment will init a new segment.
func NewSegment(path, id string, partSize int64) *Segment {
	return &Segment{
		ID:       id,
		Path:     path,
		PartSize: partSize,
		Parts:    make(map[int64]*Part),
	}
}

func (s *Segment) String() string {
	return fmt.Sprintf(
		"Segment {ID: %s, Path: %s, PartSize: %d}",
		s.ID, s.Path, s.PartSize,
	)
}

// InsertPart will insert a part into a segment and return it's Index.
// Index will start from 0.
func (s *Segment) InsertPart(offset, size int64) (p *Part, err error) {
	if size == 0 {
		panic(ErrPartSizeInvalid)
	}
	if s.PartSize == 0 {
		panic(ErrPartSizeInvalid)
	}

	s.l.Lock()
	defer s.l.Unlock()

	p = &Part{
		Offset: offset,
		Size:   size,
		Index:  int(offset / s.PartSize),
	}

	s.Parts[offset] = p
	return p, nil
}

// SortedParts will return sorted Parts.
func (s *Segment) SortedParts() []*Part {
	s.l.RLock()
	defer s.l.RUnlock()

	x := make([]*Part, 0, len(s.Parts))
	for _, v := range s.Parts {
		v := v
		x = append(x, v)
	}
	sort.Slice(x, func(i, j int) bool { return x[i].Offset < x[j].Offset })
	return x
}

// ValidateParts will validate a segment's Parts.
func (s *Segment) ValidateParts() (err error) {
	errorMessage := "%s validate Parts failed: %w"

	s.l.RLock()
	defer s.l.RUnlock()

	// Zero Parts are not allowed, cause they can't be completed.
	if len(s.Parts) == 0 {
		return fmt.Errorf(errorMessage, s, ErrSegmentPartsEmpty)
	}

	p := s.SortedParts()

	// First part offset must be 0
	if p[0].Offset != 0 {
		return fmt.Errorf(errorMessage, s, ErrSegmentNotFulfilled)
	}

	for idx := 1; idx < len(s.Parts); idx++ {
		last := p[idx-1]
		cur := p[idx]
		if last.Offset+last.Size != cur.Offset {
			return fmt.Errorf(errorMessage, s, ErrSegmentNotFulfilled)
		}
	}

	return nil
}
