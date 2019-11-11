package s3

import (
	"io"
	"sync"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/Xuanwo/storage/pkg/segment"
)

// Client is the s3 object storage client.
//
//go:generate ../../internal/bin/meta
type Client struct {
	s3 s3iface.S3API
	// options for this storager.
	workDir string // workDir dir for all operation.

	segments    map[string]*segment.Segment
	segmentLock sync.RWMutex
}

// String implements Storager.String
func (c *Client) String() string {
	panic("implement me")
}

// Init implements Storager.Init
func (c *Client) Init(pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Metadata implements Storager.Metadata
func (c *Client) Metadata() (types.Metadata, error) {
	panic("implement me")
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, pairs ...*types.Pair) iterator.ObjectIterator {
	panic("implement me")
}

// Read implements Storager.Read
func (c *Client) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	panic("implement me")
}

// Write implements Storager.Write
func (c *Client) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	panic("implement me")
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// Reach implements Storager.Reach
func (c *Client) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	panic("implement me")
}

// ListSegments implements Storager.ListSegments
func (c *Client) ListSegments(path string, pairs ...*types.Pair) iterator.SegmentIterator {
	panic("implement me")
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, pairs ...*types.Pair) (id string, err error) {
	panic("implement me")
}

// WriteSegment implements Storager.WriteSegment
func (c *Client) WriteSegment(id string, offset, size int64, r io.Reader, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// CompleteSegment implements Storager.CompleteSegment
func (c *Client) CompleteSegment(id string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}

// AbortSegment implements Storager.AbortSegment
func (c *Client) AbortSegment(id string, pairs ...*types.Pair) (err error) {
	panic("implement me")
}
