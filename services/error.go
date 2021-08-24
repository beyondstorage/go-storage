package services

import (
	"fmt"

	"github.com/beyondstorage/go-storage/v4/types"
)

type InternalError interface {
	// IsInternalError SHOULD and SHOULD ONLY be implemented by error definitions in go-storage & go-service-*.
	// We depends on the InternalError interface to distinguish our errors.
	// There's no need for user code to implement or use this function and interface.
	IsInternalError()
}

// NewErrorCode creates a new error code.
//
// Developers SHOULD use this function to define error codes (sentinel errors), instead of `NewErrorCode`
//
// Users SHOULD NOT call this function. Use defined error codes instead.
func NewErrorCode(text string) error {
	return errorCode{text}
}

type errorCode struct {
	s string
}

func (e errorCode) Error() string {
	return e.s
}

// IsInternalError implements InternalError
func (e errorCode) IsInternalError() {}

var (
	// ErrUnexpected means this is an unexpected error which go-storage can't handle
	ErrUnexpected = NewErrorCode("unexpected")

	// ErrCapabilityInsufficient means this service doesn't have this capability
	ErrCapabilityInsufficient = NewErrorCode("capability insufficient")
	// ErrRestrictionDissatisfied means this operation doesn't meat service's restriction.
	ErrRestrictionDissatisfied = NewErrorCode("restriction dissatisfied")

	// ErrObjectNotExist means the object to be operated is not exist.
	ErrObjectNotExist = NewErrorCode("object not exist")
	// ErrObjectModeInvalid means the provided object mode is invalid.
	ErrObjectModeInvalid = NewErrorCode("object mode invalid")
	// ErrPermissionDenied means this operation doesn't have enough permission.
	ErrPermissionDenied = NewErrorCode("permission denied")
	// ErrListModeInvalid means the provided list mode is invalid.
	ErrListModeInvalid = NewErrorCode("list mode invalid")
	// ErrServiceNotRegistered means this service is not registered.
	ErrServiceNotRegistered = NewErrorCode("service not registered")
	// ErrServiceInternal means this service has an internal error.
	ErrServiceInternal = NewErrorCode("service internal")
	// ErrRequestThrottled means there are too many requests.
	ErrRequestThrottled = NewErrorCode("request throttled")
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

// IsInternalError implements InternalError
func (e MetadataUnrecognizedError) IsInternalError() {}

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

// IsInternalError implements InternalError
func (e PairUnsupportedError) IsInternalError() {}

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

// IsInternalError implements InternalError
func (e PairRequiredError) IsInternalError() {}

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

// IsInternalError implements InternalError
func (e ObjectModeInvalidError) IsInternalError() {}

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

// IsInternalError implements InternalError
func (e ListModeInvalidError) IsInternalError() {}
