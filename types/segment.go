package types

// Segment is carries basic segment info.
type Segment struct {
	Path string
	ID   string
}

// Part is the index segment parts.
type Part struct {
	Index int
	Size  int64
	ETag  string
}
