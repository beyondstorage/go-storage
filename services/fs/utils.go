package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
)

// New will create a fs client.
func New(pairs ...*types.Pair) (storage.Servicer, storage.Storager, error) {
	const errorMessage = "fs New: %w"

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	store := &Storage{
		ioCopyBuffer:  io.CopyBuffer,
		ioCopyN:       io.CopyN,
		ioutilReadDir: ioutil.ReadDir,
		osCreate:      os.Create,
		osMkdirAll:    os.MkdirAll,
		osOpen:        os.Open,
		osRemove:      os.Remove,
		osRename:      os.Rename,
		osStat:        os.Stat,

		workDir: opt.WorkDir,
	}
	return nil, store, nil
}

func formatOsError(err error) error {
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	case os.IsPermission(err):
		return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
	default:
		return err
	}
}
