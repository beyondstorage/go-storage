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
func New(pairs ...*types.Pair) (_ storage.Servicer, _ storage.Storager, err error) {
	defer func() {
		if err != nil {
			err = &services.InitError{Type: Type, Err: err, Pairs: pairs}
		}
	}()

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, nil, err
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
		loose:   opt.Loose,
	}
	return nil, store, nil
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
