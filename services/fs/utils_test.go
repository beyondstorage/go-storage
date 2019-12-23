package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Xuanwo/storage/types/pairs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Xuanwo/storage/types"
)

func TestGetAbsPath(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		expectedPath string
	}{
		{"under root", "/", "abc", "/abc"},
		{"under sub dir", "/root", "abc", "/root/abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := Storage{}
			err := client.Init(pairs.WithWorkDir(tt.base))
			if err != nil {
				t.Error(err)
			}

			gotPath := client.getAbsPath(tt.path)
			assert.Equal(t, tt.expectedPath, gotPath)
		})
	}
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

func TestHandleOsError(t *testing.T) {
	t.Run("nil error will panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = handleOsError(nil)
		})
	})

	{
		tests := []struct {
			name     string
			input    error
			expected error
		}{
			{
				"not found",
				os.ErrNotExist,
				types.ErrObjectNotExist,
			},
			{
				"wrapped not found",
				fmt.Errorf("%w: some other infos", os.ErrNotExist),
				types.ErrObjectNotExist,
			},
			{
				"other errors",
				errors.New("expect unhandled error"),
				types.ErrUnhandledError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.True(t, errors.Is(handleOsError(tt.input), tt.expected))
			})
		}
	}
}
