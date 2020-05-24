// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package s3

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/pkg/httpclient"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
	ps "github.com/Xuanwo/storage/types/pairs"
)

var _ credential.Provider
var _ endpoint.Provider
var _ segment.Segment
var _ storage.Storager
var _ services.ServiceError
var _ httpclient.Options

// Type is the type for s3
const Type = "s3"

// Service available pairs.
const (
	// StorageClass will // StorageClass
	PairStorageClass = "s3_storage_class"
)

// Service available infos.
const (
	InfoObjectMetaStorageClass = "s3-storage-class"
)

// WithStorageClass will apply storage_class value to Options
// This pair is used to // StorageClass
func WithStorageClass(v string) *types.Pair {
	return &types.Pair{
		Key:   PairStorageClass,
		Value: v,
	}
}

// GetStorageClass will get storage-class value from metadata.
func GetStorageClass(m info.ObjectMeta) (string, bool) {
	v, ok := m.Get(InfoObjectMetaStorageClass)
	if !ok {
		return "", false
	}
	return v.(string), true
}

// setstorage-class will set storage-class value into metadata.
func setStorageClass(m info.ObjectMeta, v string) info.ObjectMeta {
	return m.Set(InfoObjectMetaStorageClass, v)
}

// pairStorageCreateMap holds all available pairs
var pairStorageCreateMap = map[string]struct{}{
	// Required pairs
	ps.Location: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageCreate is the parsed struct
type pairStorageCreate struct {
	// Required pairs
	Location string
	// Optional pairs
	// Generated pairs
}

// parsePairStorageCreate will parse *types.Pair slice into *pairStorageCreate
func parsePairStorageCreate(opts []*types.Pair) (*pairStorageCreate, error) {
	result := &pairStorageCreate{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageCreateMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.Location]
	if !ok {
		return nil, services.NewPairRequiredError(ps.Location)
	}
	if ok {
		result.Location = v.(string)
	}
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageDeleteMap holds all available pairs
var pairStorageDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.Location: struct{}{},
	// Generated pairs
}

// pairStorageDelete is the parsed struct
type pairStorageDelete struct {
	// Required pairs
	// Optional pairs
	HasLocation bool
	Location    string
	// Generated pairs
}

// parsePairStorageDelete will parse *types.Pair slice into *pairStorageDelete
func parsePairStorageDelete(opts []*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageDeleteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	v, ok = values[ps.Location]
	if ok {
		result.HasLocation = true
		result.Location = v.(string)
	}
	// Handle generated pairs

	return result, nil
}

// pairStorageGetMap holds all available pairs
var pairStorageGetMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.Location: struct{}{},
	// Generated pairs
}

// pairStorageGet is the parsed struct
type pairStorageGet struct {
	// Required pairs
	// Optional pairs
	HasLocation bool
	Location    string
	// Generated pairs
}

// parsePairStorageGet will parse *types.Pair slice into *pairStorageGet
func parsePairStorageGet(opts []*types.Pair) (*pairStorageGet, error) {
	result := &pairStorageGet{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageGetMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	v, ok = values[ps.Location]
	if ok {
		result.HasLocation = true
		result.Location = v.(string)
	}
	// Handle generated pairs

	return result, nil
}

// pairStorageListMap holds all available pairs
var pairStorageListMap = map[string]struct{}{
	// Required pairs
	ps.StoragerFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageList is the parsed struct
type pairStorageList struct {
	// Required pairs
	StoragerFunc storage.StoragerFunc
	// Optional pairs
	// Generated pairs
}

// parsePairStorageList will parse *types.Pair slice into *pairStorageList
func parsePairStorageList(opts []*types.Pair) (*pairStorageList, error) {
	result := &pairStorageList{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageListMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.StoragerFunc]
	if !ok {
		return nil, services.NewPairRequiredError(ps.StoragerFunc)
	}
	if ok {
		result.StoragerFunc = v.(storage.StoragerFunc)
	}
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageDeleteMap holds all available pairs
var pairStorageDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageDelete is the parsed struct
type pairStorageDelete struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageDelete will parse *types.Pair slice into *pairStorageDelete
func parsePairStorageDelete(opts []*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageDeleteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageMetadataMap holds all available pairs
var pairStorageMetadataMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageMetadata is the parsed struct
type pairStorageMetadata struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMetadata will parse *types.Pair slice into *pairStorageMetadata
func parsePairStorageMetadata(opts []*types.Pair) (*pairStorageMetadata, error) {
	result := &pairStorageMetadata{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageMetadataMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageReadMap holds all available pairs
var pairStorageReadMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageRead is the parsed struct
type pairStorageRead struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageRead will parse *types.Pair slice into *pairStorageRead
func parsePairStorageRead(opts []*types.Pair) (*pairStorageRead, error) {
	result := &pairStorageRead{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageReadMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageStatMap holds all available pairs
var pairStorageStatMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageStat is the parsed struct
type pairStorageStat struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageStat will parse *types.Pair slice into *pairStorageStat
func parsePairStorageStat(opts []*types.Pair) (*pairStorageStat, error) {
	result := &pairStorageStat{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageStatMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageWriteMap holds all available pairs
var pairStorageWriteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageWrite is the parsed struct
type pairStorageWrite struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageWrite will parse *types.Pair slice into *pairStorageWrite
func parsePairStorageWrite(opts []*types.Pair) (*pairStorageWrite, error) {
	result := &pairStorageWrite{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageWriteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageListPrefixMap holds all available pairs
var pairStorageListPrefixMap = map[string]struct{}{
	// Required pairs
	ps.ObjectFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageListPrefix is the parsed struct
type pairStorageListPrefix struct {
	// Required pairs
	ObjectFunc types.ObjectFunc
	// Optional pairs
	// Generated pairs
}

// parsePairStorageListPrefix will parse *types.Pair slice into *pairStorageListPrefix
func parsePairStorageListPrefix(opts []*types.Pair) (*pairStorageListPrefix, error) {
	result := &pairStorageListPrefix{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageListPrefixMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.ObjectFunc]
	if !ok {
		return nil, services.NewPairRequiredError(ps.ObjectFunc)
	}
	if ok {
		result.ObjectFunc = v.(types.ObjectFunc)
	}
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageListDirMap holds all available pairs
var pairStorageListDirMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.DirFunc:  struct{}{},
	ps.FileFunc: struct{}{},
	// Generated pairs
}

// pairStorageListDir is the parsed struct
type pairStorageListDir struct {
	// Required pairs
	// Optional pairs
	HasDirFunc  bool
	DirFunc     types.ObjectFunc
	HasFileFunc bool
	FileFunc    types.ObjectFunc
	// Generated pairs
}

// parsePairStorageListDir will parse *types.Pair slice into *pairStorageListDir
func parsePairStorageListDir(opts []*types.Pair) (*pairStorageListDir, error) {
	result := &pairStorageListDir{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageListDirMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	v, ok = values[ps.DirFunc]
	if ok {
		result.HasDirFunc = true
		result.DirFunc = v.(types.ObjectFunc)
	}
	v, ok = values[ps.FileFunc]
	if ok {
		result.HasFileFunc = true
		result.FileFunc = v.(types.ObjectFunc)
	}
	// Handle generated pairs

	return result, nil
}

// pairStorageCreateMap holds all available pairs
var pairStorageCreateMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageCreate is the parsed struct
type pairStorageCreate struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageCreate will parse *types.Pair slice into *pairStorageCreate
func parsePairStorageCreate(opts []*types.Pair) (*pairStorageCreate, error) {
	result := &pairStorageCreate{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageCreateMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageDeleteMap holds all available pairs
var pairStorageDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageDelete is the parsed struct
type pairStorageDelete struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageDelete will parse *types.Pair slice into *pairStorageDelete
func parsePairStorageDelete(opts []*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageDeleteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageGetMap holds all available pairs
var pairStorageGetMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageGet is the parsed struct
type pairStorageGet struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageGet will parse *types.Pair slice into *pairStorageGet
func parsePairStorageGet(opts []*types.Pair) (*pairStorageGet, error) {
	result := &pairStorageGet{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageGetMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageListMap holds all available pairs
var pairStorageListMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageList is the parsed struct
type pairStorageList struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageList will parse *types.Pair slice into *pairStorageList
func parsePairStorageList(opts []*types.Pair) (*pairStorageList, error) {
	result := &pairStorageList{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageListMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageStatMap holds all available pairs
var pairStorageStatMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageStat is the parsed struct
type pairStorageStat struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageStat will parse *types.Pair slice into *pairStorageStat
func parsePairStorageStat(opts []*types.Pair) (*pairStorageStat, error) {
	result := &pairStorageStat{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageStatMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageWriteMap holds all available pairs
var pairStorageWriteMap = map[string]struct{}{
	// Required pairs
	ps.Size: struct{}{},
	// Optional pairs
	ps.Checksum:      struct{}{},
	PairStorageClass: struct{}{},
	// Generated pairs
}

// pairStorageWrite is the parsed struct
type pairStorageWrite struct {
	// Required pairs
	Size int64
	// Optional pairs
	HasChecksum     bool
	Checksum        string
	HasStorageClass bool
	StorageClass    string
	// Generated pairs
}

// parsePairStorageWrite will parse *types.Pair slice into *pairStorageWrite
func parsePairStorageWrite(opts []*types.Pair) (*pairStorageWrite, error) {
	result := &pairStorageWrite{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageWriteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.Size]
	if !ok {
		return nil, services.NewPairRequiredError(ps.Size)
	}
	if ok {
		result.Size = v.(int64)
	}
	// Handle optional pairs
	v, ok = values[ps.Checksum]
	if ok {
		result.HasChecksum = true
		result.Checksum = v.(string)
	}
	v, ok = values[PairStorageClass]
	if ok {
		result.HasStorageClass = true
		result.StorageClass = v.(string)
	}
	// Handle generated pairs

	return result, nil
}

// pairStorageDeleteMap holds all available pairs
var pairStorageDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageDelete is the parsed struct
type pairStorageDelete struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageDelete will parse *types.Pair slice into *pairStorageDelete
func parsePairStorageDelete(opts []*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageDeleteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageMetadataMap holds all available pairs
var pairStorageMetadataMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageMetadata is the parsed struct
type pairStorageMetadata struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMetadata will parse *types.Pair slice into *pairStorageMetadata
func parsePairStorageMetadata(opts []*types.Pair) (*pairStorageMetadata, error) {
	result := &pairStorageMetadata{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageMetadataMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageReadMap holds all available pairs
var pairStorageReadMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageRead is the parsed struct
type pairStorageRead struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageRead will parse *types.Pair slice into *pairStorageRead
func parsePairStorageRead(opts []*types.Pair) (*pairStorageRead, error) {
	result := &pairStorageRead{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageReadMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageInitIndexSegmentMap holds all available pairs
var pairStorageInitIndexSegmentMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageInitIndexSegment is the parsed struct
type pairStorageInitIndexSegment struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageInitIndexSegment will parse *types.Pair slice into *pairStorageInitIndexSegment
func parsePairStorageInitIndexSegment(opts []*types.Pair) (*pairStorageInitIndexSegment, error) {
	result := &pairStorageInitIndexSegment{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageInitIndexSegmentMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageWriteIndexSegmentMap holds all available pairs
var pairStorageWriteIndexSegmentMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageWriteIndexSegment is the parsed struct
type pairStorageWriteIndexSegment struct {
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageWriteIndexSegment will parse *types.Pair slice into *pairStorageWriteIndexSegment
func parsePairStorageWriteIndexSegment(opts []*types.Pair) (*pairStorageWriteIndexSegment, error) {
	result := &pairStorageWriteIndexSegment{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageWriteIndexSegmentMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageListPrefixSegmentsMap holds all available pairs
var pairStorageListPrefixSegmentsMap = map[string]struct{}{
	// Required pairs
	ps.SegmentFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageListPrefixSegments is the parsed struct
type pairStorageListPrefixSegments struct {
	// Required pairs
	SegmentFunc segment.Func
	// Optional pairs
	// Generated pairs
}

// parsePairStorageListPrefixSegments will parse *types.Pair slice into *pairStorageListPrefixSegments
func parsePairStorageListPrefixSegments(opts []*types.Pair) (*pairStorageListPrefixSegments, error) {
	result := &pairStorageListPrefixSegments{}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageListPrefixSegmentsMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.SegmentFunc]
	if !ok {
		return nil, services.NewPairRequiredError(ps.SegmentFunc)
	}
	if ok {
		result.SegmentFunc = v.(segment.Func)
	}
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}
