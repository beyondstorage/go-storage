// Code generated by go generate via internal/cmd/meta; DO NOT EDIT.
package kodo

import (
	"context"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
)

var _ credential.Provider
var _ endpoint.Provider
var _ segment.Segment
var _ storage.Storager

// Type is the type for kodo
const Type = "kodo"

type pairServiceCreate struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasLocation bool
	Location    string
}

func parseServicePairCreate(opts ...*types.Pair) (*pairServiceCreate, error) {
	result := &pairServiceCreate{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.Location]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.Location)
	}
	if ok {
		result.HasLocation = true
		result.Location = v.(string)
	}
	return result, nil
}

type pairServiceList struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasStoragerFunc bool
	StoragerFunc    storage.StoragerFunc
}

func parseServicePairList(opts ...*types.Pair) (*pairServiceList, error) {
	result := &pairServiceList{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.StoragerFunc]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.StoragerFunc)
	}
	if ok {
		result.HasStoragerFunc = true
		result.StoragerFunc = v.(storage.StoragerFunc)
	}
	return result, nil
}

type pairServiceNew struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasCredential bool
	Credential    *credential.Provider
}

func parseServicePairNew(opts ...*types.Pair) (*pairServiceNew, error) {
	result := &pairServiceNew{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.Credential]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.Credential)
	}
	if ok {
		result.HasCredential = true
		result.Credential = v.(*credential.Provider)
	}
	return result, nil
}

type pairStorageInit struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasWorkDir bool
	WorkDir    string
}

func parseStoragePairInit(opts ...*types.Pair) (*pairStorageInit, error) {
	result := &pairStorageInit{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.WorkDir]
	if ok {
		result.HasWorkDir = true
		result.WorkDir = v.(string)
	}
	return result, nil
}

type pairStorageInitSegment struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasPartSize bool
	PartSize    int64
}

func parseStoragePairInitSegment(opts ...*types.Pair) (*pairStorageInitSegment, error) {
	result := &pairStorageInitSegment{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.PartSize]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.PartSize)
	}
	if ok {
		result.HasPartSize = true
		result.PartSize = v.(int64)
	}
	return result, nil
}

type pairStorageListDir struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasFileFunc bool
	FileFunc    types.ObjectFunc
}

func parseStoragePairListDir(opts ...*types.Pair) (*pairStorageListDir, error) {
	result := &pairStorageListDir{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.FileFunc]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.FileFunc)
	}
	if ok {
		result.HasFileFunc = true
		result.FileFunc = v.(types.ObjectFunc)
	}
	return result, nil
}

type pairStorageListSegments struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasSegmentFunc bool
	SegmentFunc    segment.Func
}

func parseStoragePairListSegments(opts ...*types.Pair) (*pairStorageListSegments, error) {
	result := &pairStorageListSegments{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.SegmentFunc]
	if ok {
		result.HasSegmentFunc = true
		result.SegmentFunc = v.(segment.Func)
	}
	return result, nil
}

type pairStorageReach struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasExpire bool
	Expire    int
}

func parseStoragePairReach(opts ...*types.Pair) (*pairStorageReach, error) {
	result := &pairStorageReach{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.Expire]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.Expire)
	}
	if ok {
		result.HasExpire = true
		result.Expire = v.(int)
	}
	return result, nil
}

type pairStorageWrite struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasChecksum     bool
	Checksum        string
	HasSize         bool
	Size            int64
	HasStorageClass bool
	StorageClass    string
}

func parseStoragePairWrite(opts ...*types.Pair) (*pairStorageWrite, error) {
	result := &pairStorageWrite{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[pairs.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[pairs.Checksum]
	if ok {
		result.HasChecksum = true
		result.Checksum = v.(string)
	}
	v, ok = values[pairs.Size]
	if !ok {
		return nil, types.NewErrPairRequired(pairs.Size)
	}
	if ok {
		result.HasSize = true
		result.Size = v.(int64)
	}
	v, ok = values[pairs.StorageClass]
	if ok {
		result.HasStorageClass = true
		result.StorageClass = v.(string)
	}
	return result, nil
}
