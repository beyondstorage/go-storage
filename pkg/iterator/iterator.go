/*
Package iterator provided iterator support for storage.

This package handles iteration details, and Storager implementer just need to write a NextObjectFunc.
*/
package iterator

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

// ErrDone means this iterator has been done.
var ErrDone = errors.New("iterator is done")

// ObjectIterator is the iterator interface used to do iterate
// TODO: waiting for golang generic support.
type ObjectIterator interface {
	Next() (*types.Object, error)
}

/*
NextObjectFunc is the func used in iterator.

Notes

- ErrDone should be return while there are no items any more.

- Input objects slice should be set every time.
*/
type NextObjectFunc func(*[]*types.Object) error

// GenericObjectIterator is the prefix based iterator.
type GenericObjectIterator struct {
	buf   []*types.Object
	index int
	next  NextObjectFunc
	done  bool
}

// NewObjectIterator will return a new prefix based iterator.
func NewObjectIterator(fn NextObjectFunc) *GenericObjectIterator {
	return &GenericObjectIterator{
		buf:   nil,
		index: 0,
		next:  fn,
	}
}

// Next implements ObjectIterator interface.
//
// Next call is not thread safe, do not call it in multi goroutine.
func (it *GenericObjectIterator) Next() (i *types.Object, err error) {
	if it.index < len(it.buf) {
		it.index++
		return it.buf[it.index-1], nil
	}
	if it.done {
		return nil, ErrDone
	}

	// Reset buf before call next.
	it.buf = nil
	// next may return empty buf without done, we should keep going.
	for err == nil && len(it.buf) == 0 {
		err = it.next(&it.buf)
	}
	if err == nil {
		it.index = 1
		return it.buf[0], nil
	}
	if !errors.Is(err, ErrDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}

	// Mark this iterator has been done, no more elem will be fetched.
	it.done = true
	if len(it.buf) == 0 {
		return nil, ErrDone
	}

	it.index = 1
	return it.buf[0], nil
}

// SegmentIterator is the iterator interface used to do iterate
// TODO: waiting for golang generic support.
type SegmentIterator interface {
	Next() (*segment.Segment, error)
}

/*
NextSegmentFunc is the func used in iterator.

Notes

- ErrDone should be return while there are no items any more.

- Input objects slice should be set every time.
*/
type NextSegmentFunc func(*[]*segment.Segment) error

// GenericSegmentIterator is the prefix based iterator.
type GenericSegmentIterator struct {
	buf   []*segment.Segment
	index int
	next  NextSegmentFunc
	done  bool
}

// NewSegmentIterator will return a new prefix based iterator.
func NewSegmentIterator(fn NextSegmentFunc) *GenericSegmentIterator {
	return &GenericSegmentIterator{
		buf:   nil,
		index: 0,
		next:  fn,
	}
}

// Next implements SegmentIterator interface.
//
// Next call is not thread safe, do not call it in multi goroutine.
func (it *GenericSegmentIterator) Next() (i *segment.Segment, err error) {
	if it.index < len(it.buf) {
		it.index++
		return it.buf[it.index-1], nil
	}
	if it.done {
		return nil, ErrDone
	}

	// Reset buf before call next.
	it.buf = nil
	err = it.next(&it.buf)
	if err == nil {
		it.index = 1
		return it.buf[0], nil
	}
	if !errors.Is(err, ErrDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}

	// Mark this iterator has been done, no more elem will be fetched.
	it.done = true
	if len(it.buf) == 0 {
		return nil, ErrDone
	}

	it.index = 1
	return it.buf[0], nil
}
