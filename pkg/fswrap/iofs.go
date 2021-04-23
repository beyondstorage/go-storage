//+build go1.16

package fswrap

import (
	"bytes"
	"errors"
	"io/fs"
	"path"

	"github.com/aos-dev/go-storage/v3/pairs"
	"github.com/aos-dev/go-storage/v3/types"
)

var (
	_ fs.FS         = fsWrapper{}
	_ fs.GlobFS     = fsWrapper{}
	_ fs.ReadDirFS  = fsWrapper{}
	_ fs.ReadFileFS = fsWrapper{}
	_ fs.StatFS     = fsWrapper{}
	_ fs.SubFS      = fsWrapper{}

	_ fs.File = &fileWrapper{}

	_ fs.FileInfo = &fileInfoWrapper{}

	_ fs.DirEntry = &dirEntryWrapper{}
)

// Fs convert a Storager to fs.FS
func Fs(s types.Storager) fs.FS {
	return fsWrapper{s}
}

type fsWrapper struct {
	store types.Storager
}

func (w fsWrapper) Open(name string) (fs.File, error) {
	o, err := w.store.Stat(name)
	if err != nil {
		return nil, err
	}
	return &fileWrapper{store: w.store, object: o}, nil
}

func (w fsWrapper) Glob(name string) ([]string, error) {
	it, err := w.store.List("", pairs.WithListMode(types.ListModePrefix))
	if err != nil {
		return nil, err
	}

	s := make([]string, 0)
	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, types.IterateDone) {
			break
		}
		if err != nil {
			return nil, err
		}

		ok, err := path.Match(name, o.Path)
		if err != nil {
			return nil, err
		}
		if ok {
			s = append(s, o.Path)
		}
	}
	return s, nil
}

func (w fsWrapper) ReadDir(name string) ([]fs.DirEntry, error) {
	it, err := w.store.List(name, pairs.WithListMode(types.ListModeDir))
	if err != nil {
		return nil, err
	}

	ds := make([]fs.DirEntry, 0)
	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, types.IterateDone) {
			break
		}
		if err != nil {
			return nil, err
		}

		ds = append(ds, dirEntryWrapper{o})
	}
	return ds, nil
}

func (w fsWrapper) ReadFile(name string) ([]byte, error) {
	var buf bytes.Buffer

	_, err := w.store.Read(name, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (w fsWrapper) Stat(name string) (fs.FileInfo, error) {
	o, err := w.store.Stat(name)
	if err != nil {
		return nil, err
	}
	return &fileInfoWrapper{object: o}, nil
}

func (w fsWrapper) Sub(dir string) (fs.FS, error) {
	panic("implement me")
}

type fileWrapper struct {
	store  types.Storager
	object *types.Object

	offset int64
	buf    bytes.Buffer
}

func (o fileWrapper) Stat() (fs.FileInfo, error) {
	return &fileInfoWrapper{o.object}, nil
}

func (o fileWrapper) Read(bs []byte) (int, error) {
	size := int64(len(bs))

	n, err := o.store.Read(o.object.Path, &o.buf, pairs.WithSize(size), pairs.WithOffset(o.offset))
	if err != nil {
		return int(n), err
	}
	o.offset += n

	nn := copy(bs, o.buf.Bytes())
	o.buf.Reset() // Reset internal buffer after copy
	return nn, nil
}

func (o fileWrapper) Close() error {
	o.store = nil
	o.object = nil
	o.offset = 0
	return nil
}

type dirEntryWrapper struct {
	object *types.Object
}

func (d dirEntryWrapper) Name() string {
	return d.object.Path
}

func (d dirEntryWrapper) IsDir() bool {
	return d.object.Mode.IsDir()
}

func (d dirEntryWrapper) Type() fs.FileMode {
	return formatFileMode(d.object.Mode)
}

func (d dirEntryWrapper) Info() (fs.FileInfo, error) {
	return &fileInfoWrapper{d.object}, nil
}
