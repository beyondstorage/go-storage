package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// StreamModeType is the stream mode type.
const StreamModeType = os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice

// Storage is the fs client.
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

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager fs {WorkDir: %s}", s.workDir)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.WorkDir = s.workDir
	return m, nil
}

// List implements Storager.List
func (s *Storage) List(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List [%s]: %w"

	opt, err := parseStoragePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	fi, err := s.ioutilReadDir(rp)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, handleOsError(err))
	}

	for _, v := range fi {
		o := &types.Object{
			ID:         filepath.Join(rp, v.Name()),
			Name:       filepath.Join(path, v.Name()),
			Size:       v.Size(),
			UpdatedAt:  v.ModTime(),
			ObjectMeta: metadata.NewObjectMeta(),
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
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	// If path is "-", return stdin directly.
	if path == "-" {
		f := os.Stdin
		if opt.HasSize {
			return iowrap.LimitReadCloser(f, opt.Size), nil
		}
		return f, nil
	}

	rp := s.getAbsPath(path)

	f, err := s.osOpen(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, handleOsError(err))
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
			return nil, fmt.Errorf(errorMessage, s, path, handleOsError(err))
		}
	}
	return f, nil
}

// WriteFile implements Storager.WriteFile
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	var f io.WriteCloser
	// If path is "-", use stdout directly.
	if path == "-" {
		f = os.Stdout
	} else {
		// Create dir for path.
		err = s.createDir(path)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}

		rp := s.getAbsPath(path)

		f, err = s.osCreate(rp)
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, handleOsError(err))
		}
	}

	if opt.HasSize {
		_, err = s.ioCopyN(f, r, opt.Size)
	} else {
		_, err = s.ioCopyBuffer(f, r, make([]byte, 1024*1024))
	}
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, handleOsError(err))
	}
	return
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	if path == "-" {
		return &types.Object{
			ID:         "-",
			Name:       "-",
			Type:       types.ObjectTypeStream,
			Size:       0,
			ObjectMeta: metadata.NewObjectMeta(),
		}, nil
	}

	rp := s.getAbsPath(path)

	fi, err := s.osStat(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, handleOsError(err))
	}

	o = &types.Object{
		ID:         rp,
		Name:       rp,
		Size:       fi.Size(),
		UpdatedAt:  fi.ModTime(),
		ObjectMeta: metadata.NewObjectMeta(),
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
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	err = s.osRemove(rp)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, handleOsError(err))
	}
	return nil
}

// Copy implements Storager.Copy
func (s *Storage) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Copy from [%s] to [%s]: %w"

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst.
	err = s.createDir(dst)
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, err)
	}

	srcFile, err := s.osOpen(rs)
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, handleOsError(err))
	}
	defer srcFile.Close()

	dstFile, err := s.osCreate(rd)
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, handleOsError(err))
	}
	defer dstFile.Close()

	_, err = s.ioCopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, handleOsError(err))
	}
	return
}

// Move implements Storager.Move
func (s *Storage) Move(src, dst string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Move from [%s] to [%s]: %w"

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst path.
	err = s.createDir(dst)
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, err)
	}

	err = s.osRename(rs, rd)
	if err != nil {
		return fmt.Errorf(errorMessage, s, src, dst, handleOsError(err))
	}
	return
}

func (s *Storage) createDir(path string) (err error) {
	errorMessage := "posixfs createDir [%s]: %w"

	rp := s.getDirPath(path)
	// Don't need to create work dir.
	if rp == s.workDir {
		return
	}

	err = s.osMkdirAll(rp, 0755)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return
}

func (s *Storage) getAbsPath(path string) string {
	return filepath.Join(s.workDir, path)
}

func (s *Storage) getDirPath(path string) string {
	if path == "" {
		return s.workDir
	}
	return filepath.Join(s.workDir, filepath.Dir(path))
}
