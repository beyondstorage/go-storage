package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
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

// NewStorager will create Storager only.
func NewStorager(pairs ...*types.Pair) (storage.Storager, error) {
	return newStorager(pairs...)
}

// newStorager will create a fs client.
func newStorager(pairs ...*types.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Op: services.OpNewStorager, Type: Type, Err: err, Pairs: pairs}
		}
	}()
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	store = &Storage{
		ioCopyBuffer:  io.CopyBuffer,
		ioCopyN:       io.CopyN,
		ioutilReadDir: ioutil.ReadDir,
		osCreate:      os.Create,
		osMkdirAll:    os.MkdirAll,
		osOpen:        os.Open,
		osRemove:      os.Remove,
		osRename:      os.Rename,
		osStat:        os.Stat,

		workDir: "/",
	}

	if opt.HasWorkDir {
		store.workDir = opt.WorkDir
	}

	// Check and create work dir
	err = store.osMkdirAll(store.workDir, 0755)
	if err != nil {
		return nil, err
	}
	return
}

func formatError(err error) error {
	// Handle error returned by os package.
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case os.IsPermission(err):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
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
