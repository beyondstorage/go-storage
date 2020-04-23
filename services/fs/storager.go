package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/mime"
	"github.com/Xuanwo/storage/services"
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

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListDir, err, path)
	}()

	opt, err := s.parsePairListDir(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	fi, err := s.ioutilReadDir(rp)
	if err != nil {
		return err
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

		if v := mime.TypeByFileName(v.Name()); v != "" {
			o.SetContentType(v)
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
	defer func() {
		err = s.formatError(services.OpRead, err, path)
	}()

	opt, err := s.parsePairRead(pairs...)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return nil, err
		}
	}

	r = f
	if opt.HasSize {
		r = iowrap.LimitReadCloser(r, opt.Size)
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}
	return r, nil
}

// WriteFile implements Storager.WriteFile
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWrite, err, path)
	}()

	opt, err := s.parsePairWrite(pairs...)
	if err != nil {
		return err
	}

	var f io.WriteCloser
	// If path is "-", use stdout directly.
	if path == "-" {
		f = os.Stdout
	} else {
		// Create dir for path.
		err = s.createDir(path)
		if err != nil {
			return err
		}

		rp := s.getAbsPath(path)

		f, err = s.osCreate(rp)
		if err != nil {
			return err
		}
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	if opt.HasSize {
		_, err = s.ioCopyN(f, r, opt.Size)
	} else {
		_, err = s.ioCopyBuffer(f, r, make([]byte, 1024*1024))
	}
	if err != nil {
		return err
	}
	return
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError(services.OpStat, err, path)
	}()

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
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Size:       fi.Size(),
		UpdatedAt:  fi.ModTime(),
		ObjectMeta: metadata.NewObjectMeta(),
	}

	if fi.IsDir() {
		o.Type = types.ObjectTypeDir
		return
	}
	if fi.Mode().IsRegular() {
		if v := mime.TypeByFileName(path); v != "" {
			o.SetContentType(v)
		}

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
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
	}()

	rp := s.getAbsPath(path)

	err = s.osRemove(rp)
	if err != nil {
		return err
	}
	return nil
}

// Copy implements Storager.Copy
func (s *Storage) Copy(src, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpCopy, err, src, dst)
	}()

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	srcFile, err := s.osOpen(rs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.osCreate(rd)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = s.ioCopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return err
	}
	return
}

// Move implements Storager.Move
func (s *Storage) Move(src, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpMove, err, src, dst)
	}()

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst path.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	err = s.osRename(rs, rd)
	if err != nil {
		return err
	}
	return
}

func (s *Storage) createDir(path string) (err error) {
	defer func() {
		err = s.formatError("create_dir", err, path)
	}()

	rp := s.getDirPath(path)
	// Don't need to create work dir.
	if rp == s.workDir {
		return
	}

	err = s.osMkdirAll(rp, 0755)
	if err != nil {
		return err
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

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
