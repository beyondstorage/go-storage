package iterator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/stretchr/testify/assert"

	"github.com/Xuanwo/storage/types"
)

func TestNewGenericIterator(t *testing.T) {
	fn := NextObjectFunc(func(informer *[]*types.Object) error {
		return nil
	})

	got := NewObjectIterator(fn)

	assert.Equal(t, 0, got.index)
	assert.Equal(t, []*types.Object(nil), got.buf)
	assert.Equal(t, fmt.Sprintf("%v", fn), fmt.Sprintf("%v", got.next))
}

func TestGenericIterator_Next(t *testing.T) {
	testErr := errors.New("test error")

	fn := NextObjectFunc(func(informer *[]*types.Object) error {
		x := make([]*types.Object, 1)
		x[0] = &types.Object{Name: "test"}
		*informer = x
		return nil
	})
	it := NewObjectIterator(fn)
	// Every call will get an element.
	i, err := it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.Equal(t, "test", i.Name)
	assert.Equal(t, 1, len(it.buf))
	assert.Equal(t, 1, it.index)

	fn = func(informer *[]*types.Object) error {
		return testErr
	}
	it = NewObjectIterator(fn)
	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, testErr))

	fn = func(informer *[]*types.Object) error {
		return ErrDone
	}
	it = NewObjectIterator(fn)
	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, ErrDone))

	fn = func(informer *[]*types.Object) error {
		x := make([]*types.Object, 2)
		x[0] = &types.Object{Name: "test1"}
		x[1] = &types.Object{Name: "test2"}
		*informer = x
		return ErrDone
	}
	it = NewObjectIterator(fn)
	// First call will get a valid item
	i, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.Equal(t, "test1", i.Name)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 1, it.index)
	// Second call will get remain value.
	i, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "test2", i.Name)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 2, it.index)
	// Third call will get Done.
	i, err = it.Next()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrDone))
	assert.Nil(t, i)
}

func TestSegmentIterator_Next(t *testing.T) {
	testErr := errors.New("test error")

	fn := NextSegmentFunc(func(informer *[]*segment.Segment) error {
		x := make([]*segment.Segment, 1)
		x[0] = &segment.Segment{ID: "test"}
		*informer = x
		return nil
	})
	it := NewSegmentIterator(fn)
	// Every call will get an element.
	i, err := it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.Equal(t, "test", i.ID)
	assert.Equal(t, 1, len(it.buf))
	assert.Equal(t, 1, it.index)

	fn = func(informer *[]*segment.Segment) error {
		return testErr
	}
	it = NewSegmentIterator(fn)
	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, testErr))

	fn = func(informer *[]*segment.Segment) error {
		return ErrDone
	}
	it = NewSegmentIterator(fn)
	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, ErrDone))

	fn = func(informer *[]*segment.Segment) error {
		x := make([]*segment.Segment, 2)
		x[0] = &segment.Segment{ID: "test1"}
		x[1] = &segment.Segment{ID: "test2"}
		*informer = x
		return ErrDone
	}
	it = NewSegmentIterator(fn)
	// First call will get a valid item
	i, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.Equal(t, "test1", i.ID)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 1, it.index)
	// Second call will get remain value.
	i, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "test2", i.ID)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 2, it.index)
	// Third call will get Done.
	i, err = it.Next()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrDone))
	assert.Nil(t, i)
}
