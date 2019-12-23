package types

import (
	"time"

	"github.com/Xuanwo/storage/types/metadata"
)

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

	// Size is the size of this Object.
	// If the Object do not have a Size, it will be 0.
	Size int64
	// UpdatedAt is the update time of this Object.
	// If the Object do not have a UpdatedAt, it will be time.Time{wall:0x0, ext:0, loc:(*time.Location)(nil)}
	UpdatedAt time.Time

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
