package segment

// Segment is carries basic segment info.
type Segment struct {
	ID   string
	Path string
}

// NewSegment will create a new segment.
func NewSegment(path, id string) *Segment {
	return &Segment{Path: path, ID: id}
}

// Func will handle a Segment.
type Func func(segment *Segment)
