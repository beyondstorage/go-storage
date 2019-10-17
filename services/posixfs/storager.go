package posixfs

import (
	"fmt"
	"io"
	"os"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
)

// Client is the posixfs client.
//
//go:generate go run ../../internal/cmd/meta_gen/main.go
type Client struct{}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, option ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "posixfs Stat path [%s] failed: %w"

	_, err = os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, err)
	}

	return nil, err
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// Reach implements Storager.Reach
func (c *Client) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	panic("implement me")
}

// CreateDir implements Storager.CreateDir
func (c *Client) CreateDir(path string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, option ...*types.Pair) iterator.ObjectIterator {
	panic("implement me")
}

// Read implements Storager.Read
func (c *Client) Read(path string, option ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

// WriteFile implements Storager.WriteFile
func (c *Client) WriteFile(path string, size int64, r io.Reader, option ...*types.Pair) (err error) {
	panic("implement me")
}

// WriteStream implements Storager.WriteStream
func (c *Client) WriteStream(path string, r io.Reader, option ...*types.Pair) (err error) {
	panic("implement me")
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// ReadSegment implements Storager.ReadSegment
func (c *Client) ReadSegment(path string, offset, size int64, option ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

// WriteSegment implements Storager.WriteSegment
func (c *Client) WriteSegment(path string, offset, size int64, r io.Reader, option ...*types.Pair) (err error) {
	panic("implement me")
}

// CompleteSegment implements Storager.CompleteSegment
func (c *Client) CompleteSegment(path string, option ...*types.Pair) (err error) {
	panic("implement me")
}

// AbortSegment implements Storager.AbortSegment
func (c *Client) AbortSegment(path string, option ...*types.Pair) (err error) {
	panic("implement me")
}
