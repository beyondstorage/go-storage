package posixfs

import (
	"os"

	"github.com/Xuanwo/storage/types"
)

func handleOsError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	if os.IsNotExist(err) {
		return types.ErrObjectNotExist
	}
	// TODO: handle other osError here.
	return err
}
