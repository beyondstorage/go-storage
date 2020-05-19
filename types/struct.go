package types

import (
	"fmt"
	"time"

	"github.com/Xuanwo/storage/types/info"
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
	// ID is the unique key in service.
	ID string
	// name is the relative path towards service's WorkDir.
	Name string
	// type should be one of "file", "stream", "dir" or "invalid".
	Type ObjectType

	// Size is the size of this object.
	// If the ObjectMeta do not have a Size, it will be 0.
	Size int64
	// UpdatedAt is the update time of this object.
	// If the object do not have a UpdatedAt, it will be time.Time{wall:0x0, ext:0, loc:(*time.Location)(nil)}
	UpdatedAt time.Time

	// metadata is the metadata of the object.
	//
	// The difference between `struct value` and `metadata` is:
	//
	// - All value in `struct` are required, caller can use them safely.
	// - All value in `metadata` are optional, caller need to check them before using.
	//
	// Two requirement must be satisfied in order to add struct value, or they need to be
	// a Metadata.
	//
	// - All storage services can provide this value in same way.
	// - User need to access this value.
	info.ObjectMeta
}

// ObjectFunc will handle an ObjectMeta.
type ObjectFunc func(object *Object)

// Pair will store option for storage service.
type Pair struct {
	Key   string
	Value interface{}
}

func (p *Pair) String() string {
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}
