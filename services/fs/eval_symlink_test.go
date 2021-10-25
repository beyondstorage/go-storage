package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalSymlinks(t *testing.T) {
	tmpDir := t.TempDir()
	t.Log(tmpDir)

	// Make sure the test base dir is not a symlink.
	tmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		name        string
		path        string
		lpath       string
		target      string
		targetExist bool
		expected    string
	}{
		{
			"symlink point to an existing target",
			filepath.Join(tmpDir, "lna"),
			"",
			filepath.Join(tmpDir, "a"),
			true,
			filepath.Join(tmpDir, "a"),
		},
		{
			"symlink point to a non-existent target",
			filepath.Join(tmpDir, "lnb"),
			"",
			filepath.Join(tmpDir, "b"),
			false,
			filepath.Join(tmpDir, "b"),
		},
		{
			"symlink point to another symlink",
			filepath.Join(tmpDir, "lnd"),
			filepath.Join(tmpDir, "lnc"),
			filepath.Join(tmpDir, "c"),
			false,
			filepath.Join(tmpDir, "c"),
		},
		{
			"symlink point to a relative target",
			filepath.Join(tmpDir, "lle"),
			"",
			tmpDir + "/./e",
			false,
			filepath.Join(tmpDir, "e"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.targetExist {
				err := os.MkdirAll(tt.target, 0755)
				if err != nil {
					t.Log(err.Error())
				}
			}

			if tt.lpath != "" {
				err := os.Symlink(tt.target, tt.lpath)
				if err != nil {
					t.Log(err.Error())
				}
				err = os.Symlink(tt.lpath, tt.path)
				if err != nil {
					t.Log(err.Error())
				}
			} else {
				err := os.Symlink(tt.target, tt.path)
				if err != nil {
					t.Log(err.Error())
				}
			}

			actualTarget, err := evalSymlinks(tt.path)
			if err != nil {
				t.Log(err.Error())
			}

			assert.Equal(t, tt.expected, actualTarget)
		})
	}
}
