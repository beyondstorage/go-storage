package services

import (
	"errors"
	"fmt"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
)

var (
	// Minor pair error that could be ignored in loose mode.

	// ErrPairNotSupported mains this operation doesn't this pair.
	ErrPairNotSupported = errors.New("pair not supported")
	// ErrStorageClassNotSupported means this service doesn't support this storage class
	ErrStorageClassNotSupported = errors.New("storage class not supported")

	// Serious pair error that affects subsequent operations

	// ErrCredentialProtocolNotSupported means this service doesn't support this credential protocol
	ErrCredentialProtocolNotSupported = errors.New("credential protocol not supported")
	// ErrPairRequired means this operation missing required pairs.
	ErrPairRequired = errors.New("pair required")
	// ErrPairConflict means this operation has conflict pairs.
	ErrPairConflict = errors.New("pair conflict")

	// Service business logic related errors.

	// ErrObjectNotExist means the object to be operated is not exist.
	ErrObjectNotExist = errors.New("object not exist")
	// ErrPermissionDenied means this operation doesn't have enough permission.
	ErrPermissionDenied = errors.New("permission denied")
)

// MinorPairError means this error will not affect subsequent operations, and could be ignored in loose mode
type MinorPairError struct {
	Op  string
	Err error

	Pairs []*types.Pair
}

func (e *MinorPairError) Error() string {
	return fmt.Sprintf("%s: %v: %s", e.Op, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *MinorPairError) Unwrap() error {
	return e.Err
}

// SeriousPairError means this error affects subsequent operations.
type SeriousPairError struct {
	Op  string
	Err error

	Pairs []*types.Pair
}

func (e *SeriousPairError) Error() string {
	return fmt.Sprintf("%s: %v: %s", e.Op, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *SeriousPairError) Unwrap() error {
	return e.Err
}

// InitError means this service init failed.
type InitError struct {
	Type string
	Err  error

	Pairs []*types.Pair
}

func (e *InitError) Error() string {
	return fmt.Sprintf("new %s: %v: %s", e.Type, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *InitError) Unwrap() error {
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
