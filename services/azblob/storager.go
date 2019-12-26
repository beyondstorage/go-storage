package azblob

import (
	"io"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the azblob service client.
//
//go:generate ../../internal/bin/meta
type Storage struct {
	name    string
	workDir string
}

// String implements Storager.String
func (s Storage) String() string {
	panic("implement me")
}

// Init implements Storager.Init
func (s Storage) Init(pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Metadata implements Storager.Metadata
func (s Storage) Metadata() (m metadata.Storage, err error) {
	panic("implement me")
}

// ListDir implements Storager.ListDir
func (s Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Read implements Storager.Read
func (s Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

// Write implements Storager.Write
func (s Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Stat implements Storager.Stat
func (s Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	panic("implement me")
}

// Delete implements Storager.Delete
func (s Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
