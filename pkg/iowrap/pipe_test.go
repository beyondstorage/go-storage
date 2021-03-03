package iowrap

import (
	"github.com/aos-dev/go-storage/v3/pkg/randbytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestPipe(t *testing.T) {
	expected := make([]byte, 16*1024*1024)
	_, err := randbytes.NewRand().Read(expected)
	if err != nil {
		t.Error(err)
	}

	r, w := Pipe()
	io.Pipe()

	go func() {
		defer w.Close()

		_, _ = w.Write(expected)
	}()

	actual, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, expected, actual)
}
