// +build linux darwin

package fs

import (
	"testing"

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
			client := Storage{workDir: tt.base}

			gotPath := client.getAbsPath(tt.path)
			assert.Equal(t, tt.expectedPath, gotPath)
		})
	}
}
