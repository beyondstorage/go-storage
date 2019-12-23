package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// StreamModeType is the stream mode type.
const StreamModeType = os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice

// Storage is the fs client.
//
//go:generate ../../internal/bin/meta
type Storage struct {
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
	osRename      func(oldpath, newpath string) error
	osStat        func(name string) (os.FileInfo, error)
}

// New will create a fs client.
func New() *Storage {
	return &Storage{
		ioCopyBuffer:  io.CopyBuffer,
		ioCopyN:       io.CopyN,
		ioutilReadDir: ioutil.ReadDir,
		osCreate:      os.Create,
		osMkdirAll:    os.MkdirAll,
		osOpen:        os.Open,
		osRemove:      os.Remove,
		osRename:      os.Rename,
		osStat:        os.Stat,
	}
}

// String implements Storager.String
func (c *Storage) String() string {
	return fmt.Sprintf("posixfs Storager {WorkDir: %s}", c.workDir)
}

// Init implements Storager.Init
func (c *Storage) Init(pairs ...*types.Pair) (err error) {
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
// Currently, there is no useful metadata for fs, just keep it empty.
func (c *Storage) Metadata() (m metadata.Metadata, err error) {
	m = make(metadata.Metadata)
	// WorkDir must be set.
	m.SetWorkDir(c.workDir)
	return m, nil
}

// Stat implements Storager.Stat
func (c *Storage) Stat(path string, option ...*types.Pair) (o *types.Object, err error) {
	errorMessage := "posixfs Stat path [%s]: %w"

	if path == "-" {
		return &types.Object{
			Name:     "-",
			Type:     types.ObjectTypeStream,
			Size:     0,
			Metadata: make(metadata.Metadata),
		}, nil
	}

	rp := c.getAbsPath(path)

	fi, err := c.osStat(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, path, handleOsError(err))
	}

	o = &types.Object{
		Name:      rp,
		Size:      fi.Size(),
		UpdatedAt: fi.ModTime(),
		Metadata:  make(metadata.Metadata),
	}

	if fi.IsDir() {
		o.Type = types.ObjectTypeDir
		return
	}
	if fi.Mode().IsRegular() {
		o.Type = types.ObjectTypeFile
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
func (c *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "posixfs Delete path [%s]: %w"

	rp := c.getAbsPath(path)

	err = c.osRemove(rp)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return nil
}

// Copy implements Storager.Copy
func (c *Storage) Copy(src, dst string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs Copy from [%s] to [%s]: %w"

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	// Create dir for dst.
	err = c.createDir(dst)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, err)
	}

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
func (c *Storage) Move(src, dst string, option ...*types.Pair) (err error) {
	errorMessage := "posixfs Move from [%s] to [%s]: %w"

	rs := c.getAbsPath(src)
	rd := c.getAbsPath(dst)

	// Create dir for dst path.
	err = c.createDir(dst)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, err)
	}

	err = c.osRename(rs, rd)
	if err != nil {
		return fmt.Errorf(errorMessage, src, dst, handleOsError(err))
	}
	return
}

// Reach implements Storager.Reach
func (c *Storage) Reach(path string, pairs ...*types.Pair) (url string, err error) {
	panic("not supported")
}

// ListDir implements Storager.ListDir
func (c *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	errorMessage := "posixfs ListDir [%s]: %w"

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, path, err)
	}

	rp := c.getAbsPath(path)

	fi, err := c.ioutilReadDir(rp)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}

	for _, v := range fi {
		o := &types.Object{
			Name:      filepath.Join(path, v.Name()),
			Size:      v.Size(),
			UpdatedAt: v.ModTime(),
			Metadata:  make(metadata.Metadata),
		}

		if v.IsDir() {
			o.Type = types.ObjectTypeDir
			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
			continue
		}

		o.Type = types.ObjectTypeFile
		if opt.HasFileFunc {
			opt.FileFunc(o)
		}
	}
	return
}

// Read implements Storager.Read
func (c *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
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
func (c *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
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
		// Create dir for path.
		err = c.createDir(path)
		if err != nil {
			return fmt.Errorf(errorMessage, path, err)
		}

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
