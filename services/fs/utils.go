package fs

import (
	"errors"
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

func handleOsError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	// Add two conditions in case of os.IsNotExist not work with fmt.Errorf
	if errors.Is(err, os.ErrNotExist) || os.IsNotExist(err) {
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}
	// TODO: handle other osError here.
	return err
}
