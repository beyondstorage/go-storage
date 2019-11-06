package posixfs

import (
	"testing"

	"github.com/Xuanwo/storage/types"
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
