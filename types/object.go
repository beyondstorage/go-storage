package types

import (
	"sync"
	"sync/atomic"
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

func NewObject(client Storager, id, name string, ty ObjectType) *Object {
	return &Object{
		ID:   id,
		Name: name,
		Type: ty,

		client: client,
		meta: objectMeta{
			m: make(map[string]interface{}),
		},
	}
}

func NewStatedObject(client Storager, id, name string, ty ObjectType) *Object {
	return &Object{
		ID:   id,
		Name: name,
		Type: ty,

		client: client,
		meta: objectMeta{
			m: make(map[string]interface{}),
		},

		// done == 1 means this object already stated, we don't need to stat it anymore.
		done: 1,
	}
}

// Object may be a *File, *Dir or a *Stream.
type Object struct {
	// ID is the unique key in service.
	ID string
	// name is the relative path towards service's WorkDir.
	Name string
	// type should be one of "file", "stream", "dir" or "invalid".
	Type ObjectType

	// client is the client in which Object is alive.
	client Storager
	// metadata is the metadata of the object.
	meta objectMeta

	done uint32
	m    sync.Mutex
}

// Borrowed from sync.Once
func (o *Object) stat() {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.statSlow()
	}
}

func (o *Object) statSlow() {
	o.m.Lock()
	defer o.m.Unlock()

	defer atomic.StoreUint32(&o.done, 1)

	if o.done == 0 {
		ob, err := o.client.Stat(o.Name)
		if err != nil {
			// Ignore all errors while object stat, just keep them empty
			return
		}
		o.meta = ob.meta
	}
}

// Get will get meta from object meta.
func (o *Object) Get(key string) (interface{}, bool) {
	o.stat()

	v, ok := o.meta.m[key]
	if !ok {
		return nil, false
	}
	return v, true
}

// Set will get meta from object meta.
func (o *Object) Set(key string, value interface{}) *Object {
	o.meta.m[key] = value
	return o
}
