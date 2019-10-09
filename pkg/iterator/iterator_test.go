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
	idx := 0
	testErr := errors.New("test error")
	fn := NextFunc(func(informer *[]types.Informer) error {
		// First call will find buf and return a valid value.
		if idx == 0 {
			x := make([]types.Informer, 1)
			x[0] = &types.Dir{Name: "test"}
			*informer = x
			idx++
			return nil
		}
		// Second call will trigger an error
		if idx == 1 {
			idx++
			return testErr
		}
		// Third call will trigger a done.
		if idx == 2 {
			idx++
			return ErrDone
		}
		panic("should not reach")
	})

	it := NewPrefixBasedIterator(fn)

	i, err := it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, i)
	assert.NotNil(t, i.(*types.Dir))
	assert.Equal(t, "test", i.(*types.Dir).Name)
	assert.Equal(t, 1, len(it.buf))
	assert.Equal(t, 1, it.index)

	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, testErr))

	i, err = it.Next()
	assert.Error(t, err)
	assert.Nil(t, i)
	assert.True(t, errors.Is(err, ErrDone))
}
