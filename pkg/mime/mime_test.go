package mime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeByFileName(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		expected string
	}{
		{"normal case", "test.pdf", "application/pdf"},
		{"no ext", "testxxx", ""},
		{"not a valid type", "test.xxx", ""},
		{"multiple dots", "test.xx.a.pdf", "application/pdf"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := TypeByFileName(tt.filename)
			assert.Equal(t, tt.expected, got)
		})
	}
}
