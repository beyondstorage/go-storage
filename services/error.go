package services

import (
	"errors"
	"fmt"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/types"
)

var (
	// ErrCapabilityInsufficient means this service doesn't have this capability
	ErrCapabilityInsufficient = errors.New("capability insufficient")
	// ErrRestrictionDissatisfied means this operation doesn't meat service's restriction.
	ErrRestrictionDissatisfied = errors.New("restriction dissatisfied")

	// ErrObjectNotExist means the object to be operated is not exist.
	ErrObjectNotExist = errors.New("object not exist")
	// ErrPermissionDenied means this operation doesn't have enough permission.
	ErrPermissionDenied = errors.New("permission denied")
)

// InitError means this service init failed.
//
// Only returned in New
type InitError struct {
	Op   string
	Type string
	Err  error

	Pairs []*types.Pair
}

func (e *InitError) Error() string {
	return fmt.Sprintf("%s %s: %v: %s", e.Type, e.Op, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *InitError) Unwrap() error {
	return e.Err
}

// ServiceError represent errors related to service.
//
// Only returned in Servicer related operations
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
//
// Only returned in Storager related operations
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

// NewMetadataNotRecognizedError will create a new MetadataUnrecognizedError.
func NewMetadataNotRecognizedError(key string, value interface{}) *MetadataUnrecognizedError {
	return &MetadataUnrecognizedError{
		Err:   ErrCapabilityInsufficient,
		Key:   key,
		Value: value,
	}
}

// MetadataUnrecognizedError means this operation meets unrecognized metadata.
type MetadataUnrecognizedError struct {
	Err error

	Key   string
	Value interface{}
}

func (e *MetadataUnrecognizedError) Error() string {
	return fmt.Sprintf("metadata unrecognized, %s, %v: %s", e.Key, e.Value, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *MetadataUnrecognizedError) Unwrap() error {
	return e.Err
}

// NewPairUnsupportedError will create a new PairUnsupportedError.
func NewPairUnsupportedError(pair *types.Pair) *PairUnsupportedError {
	return &PairUnsupportedError{
		Err:  ErrCapabilityInsufficient,
		Pair: pair,
	}
}

// PairUnsupportedError means this operation has unsupported pair.
type PairUnsupportedError struct {
	Err error

	Pair *types.Pair
}

func (e *PairUnsupportedError) Error() string {
	return fmt.Sprintf("pair unsupported, %s: %s", e.Pair, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *PairUnsupportedError) Unwrap() error {
	return e.Err
}

// NewPairRequiredError will create a new PairRequiredError.
func NewPairRequiredError(keys ...string) *PairRequiredError {
	return &PairRequiredError{
		Err:  ErrRestrictionDissatisfied,
		Keys: keys,
	}
}

// PairRequiredError means this operation has required pair but missing.
type PairRequiredError struct {
	Err error

	Keys []string
}

func (e *PairRequiredError) Error() string {
	return fmt.Sprintf("pair required, %v: %s", e.Keys, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e *PairRequiredError) Unwrap() error {
	return e.Err
}
