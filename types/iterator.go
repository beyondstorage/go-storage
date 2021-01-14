package types

import (
	"errors"
)

// IterateDone means this iterator has returned all data.
var IterateDone = errors.New("iterate is done")

// Part is the index segment parts.
type Part struct {
	Index int
	Size  int64
	ETag  string
}

type Block struct {
}
