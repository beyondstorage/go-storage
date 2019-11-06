package posixfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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
	// options for this storager.
	workDir string // workDir dir for all operation.

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

// String implements Storager.String
func (c *Client) String() string {
	return fmt.Sprintf("posixfs Storager {WorkDir %s}", c.workDir)
}

// Init implements Storager.Init
func (c *Client) Init(pairs ...*types.Pair) (err error) {
	errorMessage := "posixfs Init: %w"

	opt, err := parseStoragePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, err)
	}

	if opt.HasWorkDir {
		// TODO: validate workDir.
		c.workDir = opt.WorkDir
	}
	return nil
}

// Metadata implements Storager.Metadata
//
// Currently, there is no useful metadata for posixfs, just keep it empty.
func (c *Client) Metadata() (m types.Metadata, err error) {
	m = make(types.Metadata)
	// WorkDir must be set.
	m.SetWorkDir(c.workDir)
	return m, nil
}

// Stat implements Storager.Stat
func (c *Client) Stat(path string, option ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "posixfs Stat path [%s]: %w"

	rp := c.getAbsPath(path)

	fi, err := c.osStat(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, handleOsError(err))
	}

	o = &types.Object{
		Name:     rp,
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

	rp := c.getAbsPath(path)

	if opt.HasRecursive && opt.Recursive {
		err = c.osRemoveAll(rp)
	} else {
		err = c.osRemove(rp)
	}
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Client) Copy(src, dst string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs Copy from [%s] to [%s]: %w"

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	srcFile, err := c.osOpen(rs)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	defer srcFile.Close()

	dstFile, err := c.osCreate(rd)
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

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	err = c.osRename(rs, rd)
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

	rp := filepath.Join(c.workDir, path)

	err = c.osMkdirAll(rp, 0755)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return
}

// ListDir implements Storager.ListDir
func (c *Client) ListDir(path string, pairs ...*types.Pair) (it iterator.ObjectIterator) {
	errorMessage := "posixfs ListDir [%s]: %w"

	opt, _ := parseStoragePairListDir(pairs...)

	recursive := false
	if opt.HasRecursive && opt.Recursive {
		recursive = true
	}

	var fn iterator.NextObjectFunc
	if !recursive {
		rp := c.getAbsPath(path)

		fn = func(objects *[]*types.Object) error {
			fi, err := c.ioutilReadDir(rp)
			if err != nil {
				return fmt.Errorf(errorMessage, path, handleOsError(err))
			}

			buf := make([]*types.Object, 0, len(fi))

			for _, v := range fi {
				o := &types.Object{
					Name:     filepath.Join(path, v.Name()),
					Metadata: make(types.Metadata),
				}

				if v.IsDir() {
					o.Type = types.ObjectTypeDir
				} else {
					o.Type = types.ObjectTypeFile
				}

				o.SetSize(v.Size())
				o.SetUpdatedAt(v.ModTime())

				buf = append(buf, o)
			}

			// Set input objects
			*objects = buf
			return iterator.ErrDone
		}
	} else {
		paths := []string{path}

		fn = func(objects *[]*types.Object) error {
			if len(paths) == 0 {
				return iterator.ErrDone
			}
			p := c.getAbsPath(paths[0])
			cp := paths[0]

			fi, err := c.ioutilReadDir(p)
			if err != nil {
				return fmt.Errorf(errorMessage, path, handleOsError(err))
			}

			// Remove the first path.
			paths = paths[1:]

			buf := make([]*types.Object, 0, len(fi))

			for _, v := range fi {
				if v.IsDir() {
					paths = append(paths, filepath.Join(cp, v.Name()))
					continue
				}

				o := &types.Object{
					Name:     filepath.Join(cp, v.Name()),
					Metadata: make(types.Metadata),
					Type:     types.ObjectTypeFile,
				}

				o.SetSize(v.Size())
				o.SetUpdatedAt(v.ModTime())

				buf = append(buf, o)
			}

			// Set input objects
			*objects = buf
			return nil
		}
	}

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

	rp := c.getAbsPath(path)

	f, err := c.osOpen(rp)
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
		rp := c.getAbsPath(path)

		f, err = c.osCreate(rp)
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

func (c *Client) getAbsPath(path string) string {
	return filepath.Join(c.workDir, path)
}
