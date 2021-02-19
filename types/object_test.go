package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectMode_String(t *testing.T) {
	cases := []struct {
		name   string
		input  ObjectMode
		expect string
	}{
		{"simple case", ModeDir, "dir"},
		{"complex case", ModeDir | ModeRead | ModeLink, "dir|read|link"},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.expect, v.input.String())
		})
	}
}

func TestObjectMode_IsDir(t *testing.T) {
	cases := []struct {
		name   string
		input  ObjectMode
		expect bool
	}{
		{"simple case", ModeDir, true},
		{"complex case", ModeDir | ModeLink, true},
		{"negative case", ModeRead | ModeLink, false},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.expect, v.input.IsDir())
		})
	}
}
