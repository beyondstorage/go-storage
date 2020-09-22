package types

import (
	"errors"
	"fmt"
)

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

func NewObjectIterator(next NextObjectFunc) *ObjectIterator {
	return &ObjectIterator{
		next:  next,
		index: 0,
		done:  false,
		o:     ObjectPage{},
	}
}

func (it *ObjectIterator) Next() (object *Object, err error) {
	// Consume Data via index.
	if it.index < len(it.o.Data) {
		it.index++
		return it.o.Data[it.index-1], nil
	}
	// Return IterateDone if iterator is already done.
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.o.Data = it.o.Data[:0]

	err = it.next(&it.o)
	if err != nil && !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}
	// Make iterator to done so that we will not fetch from upstream anymore.
	if err != nil {
		it.done = true
	}
	// Return IterateDone directly if we don't have any data.
	if len(it.o.Data) == 0 {
		return nil, IterateDone
	}
	// Return the first object.
	it.index = 1
	return it.o.Data[0], nil
}

/*
NextObjectFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type NextSegmentFunc func(*SegmentPage) error

type SegmentPage struct {
	Token string
	Data  []Segment
}

type SegmentIterator struct {
	next NextSegmentFunc

	index int
	done  bool

	o SegmentPage
}

func NewSegmentIterator(next NextSegmentFunc) *SegmentIterator {
	return &SegmentIterator{
		next:  next,
		index: 0,
		done:  false,
		o:     SegmentPage{},
	}
}

func (it *SegmentIterator) Next() (object Segment, err error) {
	// Consume Data via index.
	if it.index < len(it.o.Data) {
		it.index++
		return it.o.Data[it.index-1], nil
	}
	// Return IterateDone if iterator is already done.
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.o.Data = it.o.Data[:0]

	err = it.next(&it.o)
	if err != nil && !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}
	// Make iterator to done so that we will not fetch from upstream anymore.
	if err != nil {
		it.done = true
	}
	// Return IterateDone directly if we don't have any data.
	if len(it.o.Data) == 0 {
		return nil, IterateDone
	}
	// Return the first object.
	it.index = 1
	return it.o.Data[0], nil
}

/*
NextObjectFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type NextStoragerFunc func(*StoragerPage) error

type StoragerPage struct {
	Token string
	Data  []Storager
}

type StoragerIterator struct {
	next NextStoragerFunc

	index int
	done  bool

	o StoragerPage
}

func NewStoragerIterator(next NextStoragerFunc) *StoragerIterator {
	return &StoragerIterator{
		next:  next,
		index: 0,
		done:  false,
		o:     StoragerPage{},
	}
}

func (it *StoragerIterator) Next() (object Storager, err error) {
	// Consume Data via index.
	if it.index < len(it.o.Data) {
		it.index++
		return it.o.Data[it.index-1], nil
	}
	// Return IterateDone if iterator is already done.
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.o.Data = it.o.Data[:0]

	err = it.next(&it.o)
	if err != nil && !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}
	// Make iterator to done so that we will not fetch from upstream anymore.
	if err != nil {
		it.done = true
	}
	// Return IterateDone directly if we don't have any data.
	if len(it.o.Data) == 0 {
		return nil, IterateDone
	}
	// Return the first object.
	it.index = 1
	return it.o.Data[0], nil
}
