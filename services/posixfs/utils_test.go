package posixfs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
			client := Client{}
			err := client.Init(types.WithWorkDir(tt.base))
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
			client := Client{
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
