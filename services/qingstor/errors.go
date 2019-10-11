package qingstor

import (
	"errors"
	"fmt"

	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"

	"github.com/Xuanwo/storage/types"
)

func handleQingStorError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	var e qserror.QingStorError
	if !errors.As(err, &e) {
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}

	switch e.Code {
	case "permission_denied":
		return fmt.Errorf("%w: %v", types.ErrPermissionDenied, err)
	case "object_not_exists":
		return fmt.Errorf("%w: %v", types.ErrObjectNotExist, err)
	case "invalid_access_key_id":
		return fmt.Errorf("%w: %v", types.ErrConfigIncorrect, err)
	default:
		return fmt.Errorf("%w: %v", types.ErrUnhandledError, err)
	}
}
