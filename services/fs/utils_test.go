package fs

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.beyondstorage.io/v5/services"
)

func TestNewClient(t *testing.T) {
	c, err := NewStorager()
	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func TestFormatOsError(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			"not found",
			os.ErrNotExist,
			services.ErrObjectNotExist,
		},
		{
			"not supported error",
			testErr,
			services.ErrUnexpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatError(tt.input)
			assert.True(t, errors.Is(err, tt.expected))
		})
	}
}

func BenchmarkStorage_getAbsPath(b *testing.B) {
	store := &Storage{
		workDir: "/abc/def",
	}

	b.StartTimer()
	for i := 0; i < b.N; i += 1 {
		_ = store.getAbsPath("xyz/nml")
	}
	b.StopTimer()
}
