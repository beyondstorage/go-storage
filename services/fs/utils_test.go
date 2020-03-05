package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Xuanwo/storage/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	srv, c, err := New()
	assert.Nil(t, srv)
	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func TestClient_CreateDir(t *testing.T) {
	paths := make([]string, 10)
	for k := range paths {
		paths[k] = uuid.New().String() + "/a.doc"
	}
	tests := []struct {
		name string
		err  error
	}{
		{
			"error",
			&os.PathError{Op: "mkdir", Path: paths[0], Err: errors.New("mkdir fail")},
		},
		{
			"success",
			nil,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				osMkdirAll: func(path string, perm os.FileMode) error {
					assert.Equal(t, filepath.Dir(paths[k]), path)
					assert.Equal(t, os.FileMode(0755), perm)
					return v.err
				},
			}

			err := client.createDir(paths[k])
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
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
			testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatError(tt.input)
			assert.True(t, errors.Is(err, tt.expected))
		})
	}
}
