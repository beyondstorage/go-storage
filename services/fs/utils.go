package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage/types"
)

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

func handleOsError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	// Add two conditions in case of os.IsNotExist not work with fmt.Errorf
	if errors.Is(err, os.ErrNotExist) || os.IsNotExist(err) {
		return fmt.Errorf("%w: %v", types.ErrObjectNotExist, err)
	}
	// TODO: handle other osError here.
	return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
}
