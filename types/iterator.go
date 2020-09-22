package types

import (
	"errors"
	"fmt"
)

var IterateDone = errors.New("iterate is done")

/*
NextObjectFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type NextObjectFunc func(*ObjectPage) error

type ObjectPage struct {
	Token string
	Data  []*Object
}

type ObjectIterator struct {
	next NextObjectFunc

	index int
	done  bool

	o ObjectPage
}

func (it *ObjectIterator) Next() (object *Object, err error) {
	if it.index < len(it.o.Data) {
		it.index++
		return it.o.Data[it.index-1], nil
	}
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.o.Data = it.o.Data[:0]

	err = it.next(&it.o)
	if err == nil {
		it.index = 1
		return it.o.Data[0], nil
	}
	if !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}

	// Mark this iterator has been done, no more elem will be fetched.
	it.done = true
	if len(it.o.Data) == 0 {
		return nil, IterateDone
	}

	it.index = 1
	return it.o.Data[0], nil
}
