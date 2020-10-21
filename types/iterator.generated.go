package types

import (
	"context"
	"errors"
	"fmt"
)

/*
NextObjectFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type NextObjectFunc func(ctx context.Context, page *ObjectPage) error

type ObjectPage struct {
	Status interface{}
	Data   []*Object
}

type ObjectIterator struct {
	ctx  context.Context
	next NextObjectFunc

	index int
	done  bool

	o ObjectPage
}

func NewObjectIterator(ctx context.Context, next NextObjectFunc, status interface{}) *ObjectIterator {
	return &ObjectIterator{
		ctx:   ctx,
		next:  next,
		index: 0,
		done:  false,
		o: ObjectPage{
			Status: status,
		},
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

	err = it.next(it.ctx, &it.o)
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
type NextSegmentFunc func(ctx context.Context, page *SegmentPage) error

type SegmentPage struct {
	Status interface{}
	Data   []Segment
}

type SegmentIterator struct {
	ctx  context.Context
	next NextSegmentFunc

	index int
	done  bool

	o SegmentPage
}

func NewSegmentIterator(ctx context.Context, next NextSegmentFunc, status interface{}) *SegmentIterator {
	return &SegmentIterator{
		ctx:   ctx,
		next:  next,
		index: 0,
		done:  false,
		o: SegmentPage{
			Status: status,
		},
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

	err = it.next(it.ctx, &it.o)
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
type NextStoragerFunc func(ctx context.Context, page *StoragerPage) error

type StoragerPage struct {
	Status interface{}
	Data   []Storager
}

type StoragerIterator struct {
	ctx  context.Context
	next NextStoragerFunc

	index int
	done  bool

	o StoragerPage
}

func NewStoragerIterator(ctx context.Context, next NextStoragerFunc, status interface{}) *StoragerIterator {
	return &StoragerIterator{
		ctx:   ctx,
		next:  next,
		index: 0,
		done:  false,
		o: StoragerPage{
			Status: status,
		},
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

	err = it.next(it.ctx, &it.o)
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
