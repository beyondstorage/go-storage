package types

import (
	"sync/atomic"
)

type ObjectMode uint32

const (
	ModeIrregular ObjectMode = 0
	ModeDir                  = 1 << iota
	ModeRead
	ModeLink
	ModePart
	ModeBlock
	ModePage
	ModeAppend
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

	ob, err := o.client.Stat(o.Path)
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
