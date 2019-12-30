package iowrap

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLimitedReadCloser_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		hasCall bool
		initN   int64
		n       int
		err     error
	}{
		{
			"EOF wile input -1",
			false,
			-1,
			0,
			io.EOF,
		},
		{
			"success",
			true,
			10,
			10,
			nil,
		},
		{
			"success with large buf",
			true,
			5,
			5,
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := NewMockReader(ctrl)
			c := NewMockCloser(ctrl)
			if v.hasCall {
				r.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (n int, err error) {
					return v.n, v.err
				})
			}
			lr := LimitReadCloser(struct {
				io.Reader
				io.Closer
			}{r, c}, v.initN)
			rn, err := lr.Read(make([]byte, 10))
			assert.Equal(t, v.err == nil, err == nil)
			assert.Equal(t, v.n, rn)
		})
	}
}

func TestLimitedReadCloser_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
	}{
		{
			"error",
			io.EOF,
		},
		{
			"scuccess",
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := NewMockReader(ctrl)
			c := NewMockCloser(ctrl)
			c.EXPECT().Close().DoAndReturn(func() error {
				return v.err
			})
			lr := LimitReadCloser(struct {
				io.Reader
				io.Closer
			}{r, c}, 1)
			err := lr.Close()
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestSectionedReadCloser_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		hasCall bool
		initOff int64
		initN   int64
		bufSize int
		n       int
		err     error
	}{
		{
			"EOF wile input -1",
			false,
			0,
			-1,
			10,
			0,
			io.EOF,
		},
		{
			"success",
			true,
			100,
			10,
			10,
			10,
			nil,
		},
		{
			"success with large buf",
			true,
			100,
			10,
			100,
			10,
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := NewMockReaderAt(ctrl)
			c := NewMockCloser(ctrl)
			if v.hasCall {
				r.EXPECT().ReadAt(gomock.Any(), gomock.Any()).DoAndReturn(func(p []byte, off int64) (n int, err error) {
					assert.Equal(t, v.initOff, off)
					return v.n, v.err
				})
			}
			lr := SectionReadCloser(struct {
				io.ReaderAt
				io.Closer
			}{r, c}, v.initOff, v.initN)
			rn, err := lr.Read(make([]byte, v.bufSize))
			assert.Equal(t, v.err == nil, err == nil)
			assert.Equal(t, v.n, rn)
		})
	}
}

func TestSectionedReadCloser_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
	}{
		{
			"error",
			io.EOF,
		},
		{
			"success",
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := NewMockReaderAt(ctrl)
			c := NewMockCloser(ctrl)
			c.EXPECT().Close().DoAndReturn(func() error {
				return v.err
			})
			lr := SectionReadCloser(struct {
				io.ReaderAt
				io.Closer
			}{r, c}, 1, 1)
			err := lr.Close()
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestReaderSeekerCloser_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("real reader", func(t *testing.T) {
		r := NewMockReader(ctrl)

		r.EXPECT().Read(gomock.Any()).Times(1)

		x := ReadSeekCloser(r)
		b := make([]byte, 100)
		x.Read(b)
	})
}

func TestReaderSeekerCloser_Seek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("real seeker", func(t *testing.T) {
		reader := NewMockReader(ctrl)
		seeker := NewMockSeeker(ctrl)
		r := struct {
			io.Reader
			io.Seeker
		}{
			reader,
			seeker,
		}

		seeker.EXPECT().Seek(gomock.Any(), gomock.Any()).Times(1)

		x := ReadSeekCloser(r)
		x.Seek(0, 0)
	})

	t.Run("not a seeker", func(t *testing.T) {
		reader := NewMockReader(ctrl)

		x := ReadSeekCloser(reader)
		pos, err := x.Seek(100, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), pos)
	})
}

func TestReaderCloseerCloser_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("real closer", func(t *testing.T) {
		reader := NewMockReader(ctrl)
		Closeer := NewMockCloser(ctrl)
		r := struct {
			io.Reader
			io.Closer
		}{
			reader,
			Closeer,
		}

		Closeer.EXPECT().Close().Times(1)

		x := ReadSeekCloser(r)
		x.Close()
	})

	t.Run("not a Closeer", func(t *testing.T) {
		reader := NewMockReader(ctrl)

		x := ReadSeekCloser(reader)
		err := x.Close()
		assert.NoError(t, err)
	})
}
