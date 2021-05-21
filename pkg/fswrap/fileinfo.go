package fswrap

import (
	"os"
	"time"

	"github.com/beyondstorage/go-storage/v4/types"
)

type fileInfoWrapper struct {
	object *types.Object
}

func (o fileInfoWrapper) Name() string {
	return o.object.Path
}

func (o fileInfoWrapper) Size() int64 {
	return o.object.MustGetContentLength()
}

func (o fileInfoWrapper) Mode() os.FileMode {
	return formatFileMode(o.object.Mode)
}

func (o fileInfoWrapper) ModTime() time.Time {
	return o.object.MustGetLastModified()
}

func (o fileInfoWrapper) IsDir() bool {
	return o.object.Mode.IsDir()
}

// Sys will return internal Object.
func (o fileInfoWrapper) Sys() interface{} {
	return o.object
}

func formatFileMode(om types.ObjectMode) os.FileMode {
	var m os.FileMode

	if om.IsDir() {
		m |= os.ModeDir
	}
	if om.IsAppend() {
		m |= os.ModeAppend
	}
	if om.IsLink() {
		m |= os.ModeSymlink
	}
	return m
}
