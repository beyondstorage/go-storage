package types

import (
	"sync/atomic"
)

// ObjectType is the type for object, under layer type is string.
type ObjectType string

// All available type for object.
const (
	// ObjectTypeFile means this object is file-alike object which owns content.
	ObjectTypeFile ObjectType = "file"
	// ObjectTypeDir means this object is a dir object which doesn't have content.
	ObjectTypeDir ObjectType = "dir"
	// ObjectTypeLink means this object is a link object which contains a link target to
	// another object.
	ObjectTypeLink ObjectType = "link"
	// ObjectTypeUnknown means storager can't recognize which type of this object.
	ObjectTypeUnknown ObjectType = "unknown"
)

// NewObject will create a new object with client.
func NewObject(client Storager, done bool) *Object {
	o := &Object{
		client: client,
		meta:   make(map[string]interface{}),
	}

	if done {
		// Done means this object already stated, we don't need to stat it anymore.
		o.done = 1
	}
	return o
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

	// No matter stat success or not, we only execute once.
	defer atomic.StoreUint32(&o.done, 1)

	ob, err := o.client.Stat(o.Name)
	if err != nil {
		// Ignore all errors while object stat, just keep them empty
		return
	}

	o.clone(ob)
}

// Get will get meta from object meta.
func (o *Object) Get(key string) (interface{}, bool) {
	o.stat()

	v, ok := o.meta[key]
	if !ok {
		return nil, false
	}
	return v, true
}

// Set will get meta from object meta.
func (o *Object) Set(key string, value interface{}) *Object {
	o.meta[key] = value
	return o
}
