package segment

import (
	"errors"
	"fmt"
)

var IterateDone = errors.New("iterate is done")

/*
NextSegmentFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type NextSegmentFunc func(*Page) error

type Page struct {
	Token string
	Data  []Segment
}

type Iterator struct {
	next NextSegmentFunc

	index int
	done  bool

	p Page
}

func (it *Iterator) Next() (s Segment, err error) {
	if it.index < len(it.p.Data) {
		it.index++
		return it.p.Data[it.index-1], nil
	}
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.p.Data = it.p.Data[:0]

	err = it.next(&it.p)
	if err == nil {
		it.index = 1
		return it.p.Data[0], nil
	}
	if !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}

	// Mark this iterator has been done, no more elem will be fetched.
	it.done = true
	if len(it.p.Data) == 0 {
		return nil, IterateDone
	}

	it.index = 1
	return it.p.Data[0], nil
}
