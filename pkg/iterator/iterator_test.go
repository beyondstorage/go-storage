package iterator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Xuanwo/storage/types"
)

func TestNewPrefixBasedIterator(t *testing.T) {
	fn := NextFunc(func(informer *[]types.Informer) error {
		return nil
	})

	got := NewPrefixBasedIterator(fn)

	assert.Equal(t, 0, got.index)
	assert.Equal(t, []types.Informer(nil), got.buf)
	assert.Equal(t, fmt.Sprintf("%v", fn), fmt.Sprintf("%v", got.next))
}

func TestPrefixBasedIterator_Next(t *testing.T) {
	testErr := errors.New("test error")

	fn := NextFunc(func(informer *[]types.Informer) error {
		x := make([]types.Informer, 1)
		x[0] = &types.Dir{Name: "test"}
		*informer = x
		return nil
	})
	it := NewPrefixBasedIterator(fn)
	// Every call will get an element.
	i, err := it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.NotNil(t, i.(*types.Dir))
	assert.Equal(t, "test", i.(*types.Dir).Name)
	assert.Equal(t, 1, len(it.buf))
	assert.Equal(t, 1, it.index)

	fn = func(informer *[]types.Informer) error {
		return testErr
	}
	it = NewPrefixBasedIterator(fn)
	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, testErr))

	fn = func(informer *[]types.Informer) error {
		x := make([]types.Informer, 2)
		x[0] = &types.Dir{Name: "test1"}
		x[1] = &types.Dir{Name: "test2"}
		*informer = x
		return ErrDone
	}
	it = NewPrefixBasedIterator(fn)
	// First call will get a valid item
	i, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.NotNil(t, i.(*types.Dir))
	assert.Equal(t, "test1", i.(*types.Dir).Name)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 1, it.index)
	// Second call will get remain value.
	i, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i.(*types.Dir))
	assert.Equal(t, "test2", i.(*types.Dir).Name)
	assert.Equal(t, 2, len(it.buf))
	assert.Equal(t, 2, it.index)
	// Third call will get Done.
	i, err = it.Next()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrDone))
	assert.Nil(t, i)
}
