package services

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
)

var (
	// Errors related to service capability and restriction.

	// ErrCredentialProtocolNotSupported means this service doesn't support this credential protocol
	ErrCredentialProtocolNotSupported = errors.New("credential protocol not supported")
	// ErrPairRequired means this operation missing required pairs.
	ErrPairRequired = errors.New("pair required")
	// ErrPairConflict means this operation has conflict pairs.
	ErrPairConflict = errors.New("pair conflict")
	// ErrStorageClassNotSupported means this service doesn't support this storage class
	ErrStorageClassNotSupported = errors.New("storage class not supported")

	// Errors related to service business logic.

	// ErrObjectNotExist means the object to be operated is not exist.
	ErrObjectNotExist = errors.New("object not exist")
	// ErrPermissionDenied means this operation doesn't have enough permission.
	ErrPermissionDenied = errors.New("permission denied")
)

// PairError represent errors related to pair.
type PairError struct {
	Op  string
	Err error

	Key   string
	Value interface{}
}

func (e *PairError) Error() string {
	if e.Value == nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Key, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Key, e.Value, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *PairError) Unwrap() error {
	return e.Err
}

// ServiceError represent errors related to service.
type ServiceError struct {
	Op  string
	Err error

	storage.Servicer
	Name string
}

func (e *ServiceError) Error() string {
	if e.Name == "" {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Servicer, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Servicer, e.Name, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// StorageError represent errors related to storage.
type StorageError struct {
	Op  string
	Err error

	storage.Storager
	Path []string
}

func (e *StorageError) Error() string {
	if e.Path == nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Storager, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Storager, e.Path, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *StorageError) Unwrap() error {
	return e.Err
}
