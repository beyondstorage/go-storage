package hdfs

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.beyondstorage.io/v5/services"
)

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

func TestGetAbsPath(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		expectedPath string
	}{
		{"direct path", "", "abc", "abc"},
		{"under direct path", "", "root/abc", "root/abc"},
		{"under direct path", "", "root/abc/", "root/abc"},
		{"under root", "/", "abc", "/abc"},
		{"under exist dir", "/root", "abc", "/root/abc"},
		{"under new dir", "/root", "abc/", "/root/abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := Storage{workDir: tt.base}

			getPath := client.getAbsPath(tt.path)
			assert.Equal(t, tt.expectedPath, getPath)
		})
	}
}
