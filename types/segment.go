package types

// Segment is carries basic segment info.
type Segment struct {
	Path string
	ID   string
}

type Part struct {
	Size       int64
	PartNumber int
}
