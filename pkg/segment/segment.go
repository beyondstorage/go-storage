package segment

// Segment is carries basic segment info.
type Segment interface {
	String() string

	ID() string
	Path() string
}
