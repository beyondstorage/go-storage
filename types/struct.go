package types

import (
	"github.com/Xuanwo/storage/types/metadata"
)

// ServicerType is the type for service, under layer type is string.
type ServicerType string

// StoragerType is the type for storager, under layer type is string.
type StoragerType string

// ObjectType is the type for object, under layer type is string.
type ObjectType string

// All available type for object.
const (
	ObjectTypeFile    ObjectType = "file"
	ObjectTypeStream  ObjectType = "stream"
	ObjectTypeDir     ObjectType = "dir"
	ObjectTypeInvalid ObjectType = "invalid"
)

// Object may be a *File, *Dir or a *Stream.
type Object struct {
	// name must a complete path instead of basename in POSIX.
	Name string
	// type should be one of "file", "stream", "dir" or "invalid".
	Type ObjectType

	// metadata is the metadata of the object.
	metadata.Metadata
}

// ObjectFunc will handle an Object.
type ObjectFunc func(object *Object)

// Pair will store option for storage service.
type Pair struct {
	Key   string
	Value interface{}
}
