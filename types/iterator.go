package types

import (
	"errors"
)

// IterateDone means this iterator has returned all data.
var IterateDone = errors.New("iterate is done")
