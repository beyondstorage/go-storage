package memory

import (
	"strings"
	"sync"

	"github.com/beyondstorage/go-storage/v5/types"
)

type object struct {
	mode   types.ObjectMode
	length int64

	name   string
	parent *object
	mu     sync.Mutex
	child  map[string]*object
	data   []byte
}

func newObject(name string, parent *object, mode types.ObjectMode) *object {
	return &object{
		mode:   mode,
		name:   name,
		parent: parent,
		child:  make(map[string]*object),
	}
}

func (o *object) getChild(name string) *object {
	o.mu.Lock()
	defer o.mu.Unlock()

	x, ok := o.child[name]
	if !ok {
		return nil
	}
	return x
}

func (o *object) hasChild(name string) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	_, ok := o.child[name]
	return ok
}

func (o *object) removeChild(name string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.child, name)
}

func (o *object) insertChild(name string, c *object) *object {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.child[name] = c
	return c
}

func (o *object) getObjectByPath(path string) (ro *object) {
	ro = o
	ps := strings.Split(path, "/")

	for _, v := range ps {
		if v == "" {
			continue
		}
		ro = ro.getChild(v)
		if ro == nil {
			return nil
		}
	}
	return ro
}

func (o *object) insertChildByPath(path string) *object {
	p := o
	ps := strings.Split(path, "/")
	last := len(ps) - 1

	p = o.makeDirAll(ps[:last])
	if p == nil {
		return nil
	}

	return p.insertChild(ps[last], newObject(ps[last], p, types.ModeRead))
}

func (o *object) makeDirAll(ps []string) *object {
	p := o
	for _, v := range ps {
		if v == "" {
			continue
		}
		ro := p.getChild(v)
		// If child not exist, we can create a new dir object.
		if ro == nil {
			ro = p.insertChild(v, newObject(v, p, types.ModeDir))
		}
		// If child exist but not a dir, we should return false to indict failed.
		if !ro.mode.IsDir() {
			return nil
		}

		p = ro
	}
	return p
}
