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
	opt, err := parseStoragePairNew(pairs...)
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
