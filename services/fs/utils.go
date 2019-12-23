package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xuanwo/storage/types"
)

// ParseNamespace will parse namespace for posixfs.
func ParseNamespace(s string) string {
	if strings.Contains(s, ":") {
		return s
	}
	return filepath.Join("/", s)
}

func (c *Storage) createDir(path string) (err error) {
	errorMessage := "posixfs createDir [%s]: %w"

	rp := c.getDirPath(path)
	// Don't need to create work dir.
	if rp == c.workDir {
		return
	}

	err = c.osMkdirAll(rp, 0755)
	if err != nil {
		return fmt.Errorf(errorMessage, path, handleOsError(err))
	}
	return
}

func (c *Storage) getAbsPath(path string) string {
	return filepath.Join(c.workDir, path)
}

func (c *Storage) getDirPath(path string) string {
	if path == "" {
		return c.workDir
	}
	return filepath.Join(c.workDir, filepath.Dir(path))
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
