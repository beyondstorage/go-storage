package iowrap

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCallbackifyWriter_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer, reader := NewMockWriter(ctrl), NewMockReader(ctrl)

	writer.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		return 10, io.EOF
	}).AnyTimes()

	reader.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		return 10, io.EOF
	}).AnyTimes()

	called := false
	x := CallbackWriter(writer, func(bytes []byte) {
		called = true
	})

	_, _ = io.Copy(x, reader)

	assert.True(t, called)
}
