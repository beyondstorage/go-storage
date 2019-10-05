package segment

import (
	"fmt"
	"sort"
)

// Segment will hold the whole segment operations.
type Segment struct {
	TotalSize int64
	ID        string
	Parts     []*Part
}

// Part is a part of segment.
type Part struct {
	Offset int64
	Size   int64
}

// GetPartIndex will get a part's insert index in a segment.
func (s *Segment) GetPartIndex(p *Part) (cur int, err error) {
	if p.Size == 0 {
		panic("zero size path is invalid")
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
			return 0, fmt.Errorf("part is intersected")
		}
		return
	}

	// The current Part is the last Part, and it should not have intersecting area with last Part.
	if cur == length {
		lastPart := s.Parts[cur-1]
		if lastPart.Offset+lastPart.Size > p.Offset {
			return 0, fmt.Errorf("part is intersected")
		}
		return
	}

	// The current Part is the middle Part, and it should satisfy following rules:
	// 1. No intersecting area with last Part.
	// 2. No intersecting area with next Part.
	lastPart := s.Parts[cur-1]
	nextPart := s.Parts[cur]
	if lastPart.Offset+lastPart.Size > p.Offset {
		return 0, fmt.Errorf("part is intersected")
	}

	if p.Offset+p.Size > nextPart.Offset {
		return 0, fmt.Errorf("part is intersected")
	}
	return
}

// InsertPart will insert a part into a segment.
func (s *Segment) InsertPart(p *Part) (err error) {
	cur, err := s.GetPartIndex(p)
	if err != nil {
		// TODO: format with error
		return err
	}
	s.Parts = append(s.Parts, &Part{})
	copy(s.Parts[cur+1:], s.Parts[cur:])
	s.Parts[cur] = p
	return
}

// ValidateParts will validate a segment's parts.
func (s *Segment) ValidateParts() (err error) {
	totalSize := int64(0)
	if len(s.Parts) == 0 {
		return fmt.Errorf("segment %s parts are empty", s.ID)
	}

	prePart := s.Parts[0]
	totalSize += prePart.Size
	for k, v := range s.Parts[1:] {
		if prePart.Offset+prePart.Size != v.Offset {
			return fmt.Errorf("segment is not fullfilled between part %d and %d", k-1, k)
		}
		totalSize += v.Size
	}

	if totalSize != s.TotalSize {
		return fmt.Errorf("segment size is not match with calculated parts total size")
	}
	return nil
}
