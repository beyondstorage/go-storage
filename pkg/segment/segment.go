package segment

import (
	"errors"
	"fmt"
	"sort"
)

// All errors that segment could return.
var (
	ErrPartSizeInvalid         = errors.New("part size invalid")
	ErrPartIntersected         = errors.New("part intersected")
	ErrSegmentAlreadyInitiated = errors.New("segment already initiated")
	ErrSegmentNotInitiated     = errors.New("segment not initiated")
	ErrSegmentPartsEmpty       = errors.New("segment parts are empty")
	ErrSegmentNotFulfilled     = errors.New("segment not fulfilled")
)

// Part is a part of segment.
type Part struct {
	Offset int64
	Size   int64
}

func (p *Part) String() string {
	return fmt.Sprintf("Part {Offset: %d, Size: %d}", p.Offset, p.Size)
}

// Segment will hold the whole segment operations.
type Segment struct {
	ID    string
	Parts []*Part
}

func (s *Segment) String() string {
	return fmt.Sprintf("Segment [%s]", s.ID)
}

// GetPartIndex will get a part's insert index in a segment.
func (s *Segment) GetPartIndex(p *Part) (cur int, err error) {
	errorMessage := "%s get part index with %s failed: %w"

	if p.Size == 0 {
		panic(ErrPartSizeInvalid)
	}

	length := len(s.Parts)
	// Get the index for insert.
	cur = sort.Search(length, func(i int) bool {
		return s.Parts[i].Offset >= p.Offset
	})

	// The current Part is the only Part.
	if length == 0 {
		return 0, nil
	}

	// The current Part is the first Part, and it should not have intersecting area with next Part.
	if cur == 0 {
		nextPart := s.Parts[cur]
		if p.Offset+p.Size > nextPart.Offset {
			return 0, fmt.Errorf(errorMessage, s, p, ErrPartIntersected)
		}
		return
	}

	// The current Part is the last Part, and it should not have intersecting area with last Part.
	if cur == length {
		lastPart := s.Parts[cur-1]
		if lastPart.Offset+lastPart.Size > p.Offset {
			return 0, fmt.Errorf(errorMessage, s, p, ErrPartIntersected)
		}
		return
	}

	// The current Part is the middle Part, and it should satisfy following rules:
	// 1. No intersecting area with last Part.
	// 2. No intersecting area with next Part.
	lastPart := s.Parts[cur-1]
	nextPart := s.Parts[cur]
	if lastPart.Offset+lastPart.Size > p.Offset {
		return 0, fmt.Errorf(errorMessage, s, p, ErrPartIntersected)
	}

	if p.Offset+p.Size > nextPart.Offset {
		return 0, fmt.Errorf(errorMessage, s, p, ErrPartIntersected)
	}
	return
}

// InsertPart will insert a part into a segment.
func (s *Segment) InsertPart(p *Part) (err error) {
	cur, err := s.GetPartIndex(p)
	if err != nil {
		return err
	}
	s.Parts = append(s.Parts, &Part{})
	copy(s.Parts[cur+1:], s.Parts[cur:])
	s.Parts[cur] = p
	return
}

// ValidateParts will validate a segment's parts.
func (s *Segment) ValidateParts() (err error) {
	errorMessage := "%s validate parts failed: %w"

	// Zero parts are not allowed, cause they can't be completed.
	if len(s.Parts) == 0 {
		return fmt.Errorf(errorMessage, s, ErrSegmentPartsEmpty)
	}

	// Check parts continuity
	prePart := s.Parts[0]
	for _, v := range s.Parts[1:] {
		if prePart.Offset+prePart.Size != v.Offset {
			return fmt.Errorf(errorMessage, s, ErrSegmentNotFulfilled)
		}
	}
	return nil
}
