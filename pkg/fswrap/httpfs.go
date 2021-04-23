package fswrap

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/types"
)

var (
	_ http.FileSystem = httpFsWrapper{}

	_ http.File = httpFileWrapper{}
)

// HttpFs convert a Storager to http.FileSystem
func HttpFs(s types.Storager) http.FileSystem {
	return httpFsWrapper{s}
}

type httpFsWrapper struct {
	store types.Storager
}

func (h httpFsWrapper) Open(name string) (http.File, error) {
	o, err := h.store.Stat(name)
	if err != nil {
		return nil, err
	}
	return &httpFileWrapper{store: h.store, object: o}, nil
}

type httpFileWrapper struct {
	store  types.Storager
	object *types.Object

	offset int64
	buf    bytes.Buffer
}

func (h httpFileWrapper) Close() error {
	h.store = nil
	h.object = nil
	h.offset = 0
	return nil
}

func (h httpFileWrapper) Read(bs []byte) (int, error) {
	size := int64(len(bs))

	n, err := h.store.Read(h.object.Path, &h.buf, pairs.WithSize(size), pairs.WithOffset(h.offset))
	if err != nil {
		return int(n), err
	}
	h.offset += n

	nn := copy(bs, h.buf.Bytes())
	h.buf.Reset() // Reset internal buffer after copy
	return nn, nil
}

func (h httpFileWrapper) Seek(offset int64, whence int) (int64, error) {
	size := h.object.MustGetContentLength()

	switch whence {
	case io.SeekStart:
		h.offset = offset
	case io.SeekCurrent:
		h.offset += offset
	case io.SeekEnd:
		// TODO: Do we need to check value here?
		h.offset = size - offset
	}
	return h.offset, nil
}

func (h httpFileWrapper) Readdir(count int) ([]os.FileInfo, error) {
	if !h.object.Mode.IsDir() {
		return nil, os.ErrInvalid
	}

	it, err := h.store.List(h.object.Path, pairs.WithListMode(types.ListModeDir))
	if err != nil {
		return nil, err
	}

	// Change the meaning of n for the implementation below.
	//
	// The n above was for the public interface of "if n <= 0,
	// Readdir returns all the FileInfo from the directory in a
	// single slice".
	//
	// But below, we use only negative to mean looping until the
	// end and positive to mean bounded, with positive
	// terminating at 0.
	if count == 0 {
		count = -1
	}

	fi := make([]os.FileInfo, 0)
	for count != 0 {
		o, err := it.Next()
		if err != nil && errors.Is(err, types.IterateDone) {
			break
		}
		if err != nil {
			return nil, err
		}

		fi = append(fi, fileInfoWrapper{o})

		if count > 0 {
			count--
		}
	}
	return fi, nil
}

func (h httpFileWrapper) Stat() (os.FileInfo, error) {
	return &fileInfoWrapper{object: h.object}, nil
}
