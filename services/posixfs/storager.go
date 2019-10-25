package posixfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
)

// StreamModeType is the stream mode type.
const StreamModeType = os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice

// Client is the posixfs client.
//
//go:generate go run ../../internal/cmd/meta_gen/main.go
type Client struct {
	// All stdlib call will be added here for better unit test.
	ioCopyBuffer  func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)
	ioCopyN       func(dst io.Writer, src io.Reader, n int64) (written int64, err error)
	ioutilReadDir func(dirname string) ([]os.FileInfo, error)
	osCreate      func(name string) (*os.File, error)
	osMkdirAll    func(path string, perm os.FileMode) error
	osOpen        func(name string) (*os.File, error)
	osRemove      func(name string) error
	osRemoveAll   func(name string) error
	osRename      func(oldpath, newpath string) error
	osStat        func(name string) (os.FileInfo, error)
}

// NewClient will create a posix client.
func NewClient() *Client {
	return &Client{
		ioCopyBuffer:  io.CopyBuffer,
		ioCopyN:       io.CopyN,
		ioutilReadDir: ioutil.ReadDir,
		osCreate:      os.Create,
		osMkdirAll:    os.MkdirAll,
		osOpen:        os.Open,
		osRemove:      os.Remove,
		osRemoveAll:   os.RemoveAll,
		osRename:      os.Rename,
		osStat:        os.Stat,
	}
}

// Metadata implements Storager.Metadata
//
// Currently, there is no useful metadata for posixfs, just keep it empty.
func (c *Client) Metadata() (m types.Metadata, err error) {
	m = make(types.Metadata)
	return m, nil
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, option ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "posixfs Stat path [%s]: %w"

	fi, err := c.osStat(path)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, handleOsError(err))
	}

	o = &types.Object{
		Name:     path,
		Metadata: make(types.Metadata),
	}

	if fi.IsDir() {
		o.Type = types.ObjectTypeDir
		return
	}
	if fi.Mode().IsRegular() {
		o.Type = types.ObjectTypeFile
		o.SetSize(fi.Size())
		return
	}
	if fi.Mode()&StreamModeType != 0 {
		o.Type = types.ObjectTypeStream
		return
	}

	o.Type = types.ObjectTypeInvalid
	return o, nil
}

// Delete implements Storager.Delete
func (c *Client) Delete(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "posixfs Delete path [%s]: %w"

	opt, err := parseStoragePairDelete(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, path, err)
	}

	if opt.HasRecursive && opt.Recursive {
		err = c.osRemoveAll(path)
	} else {
		err = c.osRemove(path)
	}
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs Copy from [%s] to [%s]: %w"

	srcFile, err := c.osOpen(src)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	defer srcFile.Close()

	dstFile, err := c.osCreate(dst)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	defer dstFile.Close()

	_, err = c.ioCopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	return
}

// Move implements Storager.Move
func (c *Client) Move(src, dst string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs Move from [%s] to [%s]: %w"

	err = c.osRename(src, dst)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	return
}

// Reach implements Storager.Reach
func (c *Client) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	panic("not supported")
}

// CreateDir implements Storager.CreateDir
func (c *Client) CreateDir(path string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs CreateDir [%s]: %w"

	err = c.osMkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, option ...*types.Pair) (it iterator.ObjectIterator) {
	errorMessage := "posixfs ListDir [%s]: %w"

	fn := iterator.NextObjectFunc(func(objects *[]*types.Object) error {
		fi, err := c.ioutilReadDir(path)
		if err != nil {
			return fmt.Errorf(errorMessage, path, handleOsError(err))
		}

		idx := 0
		buf := make([]*types.Object, len(fi))

		for _, v := range fi {
			o := &types.Object{
				Name:     v.Name(),
				Metadata: make(types.Metadata),
			}

			if v.IsDir() {
				o.Type = types.ObjectTypeDir
			} else {
				o.Type = types.ObjectTypeFile
			}

			o.SetSize(v.Size())
			o.SetUpdatedAt(v.ModTime())

			buf[idx] = o
			idx++
		}

		// Set input objects
		*objects = buf[:idx]
		return iterator.ErrDone
	})

	it = iterator.NewObjectIterator(fn)
	return
}

// Read implements Storager.Read
func (c *Client) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	errorMessage := "posixfs Read [%s]: %w"

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, err)
	}

	// If path is "-", return stdin directly.
	if path == "-" {
		f := os.Stdin
		if opt.HasSize {
			return iowrap.LimitReadCloser(f, opt.Size), nil
		}
		return f, nil
	}

	f, err := c.osOpen(path)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	if opt.HasSize && opt.HasOffset {
		return iowrap.SectionReadCloser(f, opt.Offset, opt.Size), nil
	}
	if opt.HasSize {
		return iowrap.LimitReadCloser(f, opt.Size), nil
	}
	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return nil, fmt.Errorf(errorMessage, path, handleOsError(err))
		}
	}
	return f, nil
}

// WriteFile implements Storager.WriteFile
func (c *Client) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	errorMessage := "posixfs WriteFile [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, path, err)
	}

	var f io.WriteCloser
	// If path is "-", use stdout directly.
	if path == "-" {
		f = os.Stdout
	} else {
		f, err = c.osCreate(path)
		if err != nil {
			return fmt.Errorf(errorMessage, path, handleOsError(err))
		}
	}

	if opt.HasSize {
		_, err = c.ioCopyN(f, r, opt.Size)
	} else {
		_, err = c.ioCopyBuffer(f, r, make([]byte, 1024*1024))
	}
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return
}

// ListSegments implements Storager.ListSegments
func (c *Client) ListSegments(path string, option ...*types.Pair) iterator.SegmentIterator {
	panic("implement me")
}

// InitSegment implements Storager.InitSegment
func (c *Client) InitSegment(path string, option ...*types.Pair) (id string, err error) {
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
