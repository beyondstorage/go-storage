package types

import (
	"fmt"
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
	ObjectMeta
}

// Pair will store option for storage service.
type Pair struct {
	Key   string
	Value interface{}
}

func (p *Pair) String() string {
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}

// StorageMeta is the static metadata for StorageMeta.
type StorageMeta struct {
	Name    string
	WorkDir string

	m map[string]interface{}
}

// NewStorageMeta will create a new StorageMeta metadata.
func NewStorageMeta() StorageMeta {
	return StorageMeta{
		m: make(map[string]interface{}),
	}
}

// StorageStatistic is the statistic metadata for StorageMeta.
type StorageStatistic struct {
	m map[string]interface{}
}

// NewStorageStatistic will create a new StorageMeta statistic.
func NewStorageStatistic() StorageStatistic {
	return StorageStatistic{
		m: make(map[string]interface{}),
	}
}

// ObjectMeta is the metadata for ObjectMeta.
type ObjectMeta struct {
	m map[string]interface{}
}

// NewObjectMeta will create a new ObjectMeta metadata.
func NewObjectMeta() ObjectMeta {
	return ObjectMeta{
		m: make(map[string]interface{}),
	}
}

// Get will get meta from object meta.
func (m ObjectMeta) Get(key string) (interface{}, bool) {
	v, ok := m.m[key]
	if !ok {
		return nil, false
	}
	return v, true
}

// Set will get meta from object meta.
func (m ObjectMeta) Set(key string, value interface{}) ObjectMeta {
	m.m[key] = value
	return m
}
