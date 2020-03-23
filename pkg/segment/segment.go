package segment

// Segment is carries basic segment info.
type Segment interface {
	String() string

	ID() string
	Path() string
}

// Func will handle a Segment.
type Func func(segment Segment)
