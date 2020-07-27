package types

import (
	"fmt"

	"github.com/aos-dev/go-storage/v2/types/info"
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

	// metadata is the metadata of the object.
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
