package types

import (
	"strings"
	"sync/atomic"
)

// ObjectMode describes what users can operate on this object.
//
// Different object modes are orthogonal, and an object could have different
// object modes at the same time.
//
// For example:
//
// - ModeDir means we can do list on this object('s path)
// - ModeRead means we can read it as a normal file
// - ModeLink means we can read this object's link target
//
// And we can compose them together:
//
// - ModeRead | ModeLink: Think about a symlink, we can still read it.
// - ModeDir | ModeLink: Think about a symlink to dir.
// - ModeAppend | ModeRead: file in fs could be both read and append.
//
// Reference
//
// - GSP-25: https://github.com/beyondstorage/specs/blob/master/rfcs/25-object-mode.md
// - Core Concept Object: https://beyondstorage.io/docs/go-storage/internal/core-concept#object
type ObjectMode uint32

// All available object mode
const (
	// ModeDir means this Object represents a dir which can be used to list with dir mode.
	ModeDir ObjectMode = 1 << iota
	// ModeRead means this Object can be used to read content.
	ModeRead
	// ModeLink means this Object is a link which targets to another Object.
	ModeLink
	// ModePart means this Object is a Multipart Object which can be used for multipart operations.
	ModePart
	// ModeBlock means this Object is a Block Object which can be used for block operations.
	ModeBlock
	// ModePage means this Object is a Page Object which can be used for random write with offset.
	ModePage
	// ModeAppend means this Object is a Append Object which can be used for append.
	ModeAppend
)

// String implement Stringer for ObjectMode.
//
// An object with Read,Append will print like "read|append"
func (o ObjectMode) String() string {
	s := make([]string, 0)
	if o.IsDir() {
		s = append(s, "dir")
	}
	if o.IsRead() {
		s = append(s, "read")
	}
	if o.IsLink() {
		s = append(s, "link")
	}
	if o.IsPart() {
		s = append(s, "part")
	}
	if o.IsBlock() {
		s = append(s, "block")
	}
	if o.IsPage() {
		s = append(s, "page")
	}
	if o.IsAppend() {
		s = append(s, "append")
	}

	return strings.Join(s, "|")
}

// Add support add ObjectMode into current one.
func (o *ObjectMode) Add(mode ObjectMode) {
	*o |= mode
}

// Del support delete ObjectMode from current one.
func (o *ObjectMode) Del(mode ObjectMode) {
	*o &= ^mode
}
func (o ObjectMode) IsDir() bool {
	return o&ModeDir != 0
}
func (o ObjectMode) IsRead() bool {
	return o&ModeRead != 0
}
func (o ObjectMode) IsLink() bool {
	return o&ModeLink != 0
}
func (o ObjectMode) IsPart() bool {
	return o&ModePart != 0
}
func (o ObjectMode) IsBlock() bool {
	return o&ModeBlock != 0
}
func (o ObjectMode) IsPage() bool {
	return o&ModePage != 0
}
func (o ObjectMode) IsAppend() bool {
	return o&ModeAppend != 0
}

// Part is the part of Multipart Object.
type Part struct {
	Index int
	Size  int64
	ETag  string
}

// Block is the block of Block Object.
type Block struct {
	ID   string
	Size int64
}

// NewObject will create a new object with client.
func NewObject(client Storager, done bool) *Object {
	o := &Object{
		client: client,
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
