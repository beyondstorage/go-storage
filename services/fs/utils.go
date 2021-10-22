package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Std{in/out/err} support
const (
	Stdin  = "/dev/stdin"
	Stdout = "/dev/stdout"
	Stderr = "/dev/stderr"
)

// Storage is the fs client.
type Storage struct {
	// options for this storager.
	workDir string // workDir dir for all operation.

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedCopier
	typ.UnimplementedMover
	typ.UnimplementedFetcher
	typ.UnimplementedAppender
	typ.UnimplementedDirer
	typ.UnimplementedLinker
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager fs {WorkDir: %s}", s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	return newStorager(pairs...)
}

// newStorager will create a fs client.
func newStorager(pairs ...typ.Pair) (store *Storage, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return
	}

	store = &Storage{
		workDir: "/",
	}

	if opt.HasDefaultStoragePairs {
		store.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		store.features = opt.StorageFeatures
	}
	if opt.HasWorkDir {
		workDir, err := evalSymlinks(opt.WorkDir)
		if err != nil {
			return nil, err
		}
		store.workDir = workDir
	}

	// Check and create work dir
	err = os.MkdirAll(store.workDir, 0755)
	if err != nil {
		return nil, err
	}
	return
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	// Handle error returned by os package.
	switch {
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

func (s *Storage) openFile(absPath string, mode int) (f *os.File, needClose bool, err error) {
	switch absPath {
	case Stdin:
		f = os.Stdin
	case Stdout:
		f = os.Stdout
	case Stderr:
		f = os.Stderr
	default:
		needClose = true
		f, err = os.OpenFile(absPath, mode, 0664)
	}

	return
}

func (s *Storage) createFile(absPath string) (f *os.File, needClose bool, err error) {
	return s.createFileWithFlag(absPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
}

func (s *Storage) createFileWithFlag(absPath string, flag int) (f *os.File, needClose bool, err error) {
	switch absPath {
	case Stdin:
		return os.Stdin, false, nil
	case Stdout:
		return os.Stdout, false, nil
	case Stderr:
		return os.Stderr, false, nil
	}

	fi, err := os.Lstat(absPath)
	if err == nil {
		// File is exist, let's check if the file is a dir or a symlink.
		if fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 {
			return nil, false, services.ErrObjectModeInvalid
		}
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Something error other than ErrNotExist happened, return directly.
		return
	}
	// Set stat error to nil.
	err = nil

	// The file is not exist, we should create the dir and create the file.
	if fi == nil {
		err = os.MkdirAll(filepath.Dir(absPath), 0755)
		if err != nil {
			return nil, false, err
		}
	}

	// There are two situations we handled here:
	// - The file is exist and not a dir
	// - The file is not exist
	f, err = os.OpenFile(absPath, flag, 0666)
	if err != nil {
		return nil, false, err
	}
	return f, true, nil
}

func (s *Storage) statFile(absPath string) (fi os.FileInfo, err error) {
	switch absPath {
	case Stdin:
		fi, err = os.Stdin.Stat()
	case Stdout:
		fi, err = os.Stdout.Stat()
	case Stderr:
		fi, err = os.Stderr.Stat()
	default:
		// Use Lstat here to not follow symlinks.
		// We will resolve symlinks target while this object's type is link.
		fi, err = os.Lstat(absPath)
	}

	return
}

func (s *Storage) getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	absPath := filepath.Join(s.workDir, path)

	return absPath
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
