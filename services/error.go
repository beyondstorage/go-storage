package services

import (
	"errors"
	"fmt"

	"github.com/aos-dev/go-storage/v3/types"
)

var (
	// ErrUnexpected means this is an unexpected error which go-storage can't handle
	ErrUnexpected = errors.New("unexpected")

	// ErrCapabilityInsufficient means this service doesn't have this capability
	ErrCapabilityInsufficient = errors.New("capability insufficient")
	// ErrRestrictionDissatisfied means this operation doesn't meat service's restriction.
	ErrRestrictionDissatisfied = errors.New("restriction dissatisfied")

	// ErrObjectNotExist means the object to be operated is not exist.
	ErrObjectNotExist = errors.New("object not exist")
	// ErrObjectModeInvalid means the provided object mode is invalid.
	ErrObjectModeInvalid = errors.New("object mode invalid")
	// ErrPermissionDenied means this operation doesn't have enough permission.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrListModeInvalid means the provided list mode is invalid.
	ErrListModeInvalid = errors.New("list mode invalid")
	// ErrServiceNotRegistered means this service is not registered.
	ErrServiceNotRegistered = errors.New("service not registered")
)

// InitError means this service init failed.
//
// Only returned in New
type InitError struct {
	Op   string
	Type string
	Err  error

	Pairs []types.Pair
}

func (e InitError) Error() string {
	return fmt.Sprintf("%s %s: %v: %s", e.Type, e.Op, e.Pairs, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e InitError) Unwrap() error {
	return e.Err
}

// ServiceError represent errors related to service.
//
// Only returned in Servicer related operations
type ServiceError struct {
	Op  string
	Err error

	types.Servicer
	Name string
}

func (e ServiceError) Error() string {
	if e.Name == "" {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Servicer, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Servicer, e.Name, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e ServiceError) Unwrap() error {
	return e.Err
}

// StorageError represent errors related to storage.
//
// Only returned in Storager related operations
type StorageError struct {
	Op  string
	Err error

	types.Storager
	Path []string
}

func (e StorageError) Error() string {
	if e.Path == nil {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Storager, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s, %s: %s", e.Op, e.Storager, e.Path, e.Err.Error())
}

// Unwrap implements xerrors.Wrapper
func (e StorageError) Unwrap() error {
	return e.Err
}

// MetadataUnrecognizedError means this operation meets unrecognized metadata.
type MetadataUnrecognizedError struct {
	Key   string
	Value interface{}
}

func (e MetadataUnrecognizedError) Error() string {
	return fmt.Sprintf("metadata unrecognized, %s, %v: %s", e.Key, e.Value, ErrCapabilityInsufficient.Error())
}

// Unwrap implements xerrors.Wrapper
func (e MetadataUnrecognizedError) Unwrap() error {
	return ErrCapabilityInsufficient
}

// PairUnsupportedError means this operation has unsupported pair.
type PairUnsupportedError struct {
	Pair types.Pair
}

func (e PairUnsupportedError) Error() string {
	return fmt.Sprintf("pair unsupported, %s: %s", e.Pair, ErrCapabilityInsufficient.Error())
}

// Unwrap implements xerrors.Wrapper
func (e PairUnsupportedError) Unwrap() error {
	return ErrCapabilityInsufficient
}

// PairRequiredError means this operation has required pair but missing.
type PairRequiredError struct {
	Keys []string
}

func (e PairRequiredError) Error() string {
	return fmt.Sprintf("pair required, %v: %s", e.Keys, ErrRestrictionDissatisfied.Error())
}

// Unwrap implements xerrors.Wrapper
func (e PairRequiredError) Unwrap() error {
	return ErrRestrictionDissatisfied
}

// ObjectModeInvalidError means the provided object mode is invalid.
type ObjectModeInvalidError struct {
	Expected types.ObjectMode
	Actual   types.ObjectMode
}

func (e ObjectModeInvalidError) Error() string {
	return fmt.Sprintf("object mode invalid, expected %v, actual %v", e.Expected, e.Actual)
}

func (e ObjectModeInvalidError) Unwrap() error {
	return ErrObjectModeInvalid
}

// ListModeInvalidError means the provided list mode is invalid.
type ListModeInvalidError struct {
	Actual types.ListMode
}

func (e ListModeInvalidError) Error() string {
	return fmt.Sprintf("list mode invalid, actual %v", e.Actual)
}

func (e ListModeInvalidError) Unwrap() error {
	return ErrListModeInvalid
}
