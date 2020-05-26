// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package s3

import (
	"context"
	"io"

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

// pairServiceNewMap holds all available pairs
var pairServiceNewMap = map[string]struct{}{
	// Required pairs
	ps.Credential: struct{}{},
	// Optional pairs
	ps.Endpoint: struct{}{},
	// Generated pairs
	ps.HTTPClientOptions: struct{}{},
}

// pairServiceNew is the parsed struct
type pairServiceNew struct {
	pairs []*types.Pair

	// Required pairs
	Credential *credential.Provider
	// Optional pairs
	HasEndpoint bool
	Endpoint    endpoint.Provider
	// Generated pairs
	HasHTTPClientOptions bool
	HTTPClientOptions    *httpclient.Options
}

// parsePairServiceNew will parse *types.Pair slice into *pairServiceNew
func parsePairServiceNew(opts []*types.Pair) (*pairServiceNew, error) {
	result := &pairServiceNew{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}
	var v interface{}
	var ok bool

	// Handle required pairs
	v, ok = values[ps.Credential]
	if !ok {
		return nil, services.NewPairRequiredError(ps.Credential)
	}
	if ok {
		result.Credential = v.(*credential.Provider)
	}
	// Handle optional pairs
	v, ok = values[ps.Endpoint]
	if ok {
		result.HasEndpoint = true
		result.Endpoint = v.(endpoint.Provider)
	}
	// Handle generated pairs
	v, ok = values[ps.HTTPClientOptions]
	if ok {
		result.HasHTTPClientOptions = true
		result.HTTPClientOptions = v.(*httpclient.Options)
	}

	return result, nil
}

// pairServiceCreateMap holds all available pairs
var pairServiceCreateMap = map[string]struct{}{
	// Required pairs
	ps.Location: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairServiceCreate is the parsed struct
type pairServiceCreate struct {
	pairs []*types.Pair

	// Required pairs
	Location string
	// Optional pairs
	// Generated pairs
}

// parsePairServiceCreate will parse *types.Pair slice into *pairServiceCreate
func parsePairServiceCreate(opts []*types.Pair) (*pairServiceCreate, error) {
	result := &pairServiceCreate{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairServiceCreateMap[v.Key]; !ok {
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

// pairServiceDeleteMap holds all available pairs
var pairServiceDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.Location: struct{}{},
	// Generated pairs
}

// pairServiceDelete is the parsed struct
type pairServiceDelete struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	HasLocation bool
	Location    string
	// Generated pairs
}

// parsePairServiceDelete will parse *types.Pair slice into *pairServiceDelete
func parsePairServiceDelete(opts []*types.Pair) (*pairServiceDelete, error) {
	result := &pairServiceDelete{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairServiceDeleteMap[v.Key]; !ok {
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

// pairServiceGetMap holds all available pairs
var pairServiceGetMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.Location: struct{}{},
	// Generated pairs
}

// pairServiceGet is the parsed struct
type pairServiceGet struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	HasLocation bool
	Location    string
	// Generated pairs
}

// parsePairServiceGet will parse *types.Pair slice into *pairServiceGet
func parsePairServiceGet(opts []*types.Pair) (*pairServiceGet, error) {
	result := &pairServiceGet{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairServiceGetMap[v.Key]; !ok {
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

// pairServiceListMap holds all available pairs
var pairServiceListMap = map[string]struct{}{
	// Required pairs
	ps.StoragerFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairServiceList is the parsed struct
type pairServiceList struct {
	pairs []*types.Pair

	// Required pairs
	StoragerFunc storage.StoragerFunc
	// Optional pairs
	// Generated pairs
}

// parsePairServiceList will parse *types.Pair slice into *pairServiceList
func parsePairServiceList(opts []*types.Pair) (*pairServiceList, error) {
	result := &pairServiceList{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairServiceListMap[v.Key]; !ok {
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

// Create will create a new storager instance.
//
// This function will create a context by default.
func (s *Service) Create(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	ctx := context.Background()
	return s.CreateWithContext(ctx, name, pairs...)
}

// CreateWithContext will create a new storager instance.
func (s *Service) CreateWithContext(ctx context.Context, name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError(services.OpCreate, err, name)
	}()
	var opt *pairServiceCreate
	opt, err = parsePairServiceCreate(pairs)
	if err != nil {
		return
	}

	return s.create(ctx, name, opt)
}

// Delete will delete a storager instance.
//
// This function will create a context by default.
func (s *Service) Delete(name string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.DeleteWithContext(ctx, name, pairs...)
}

// DeleteWithContext will delete a storager instance.
func (s *Service) DeleteWithContext(ctx context.Context, name string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, name)
	}()
	var opt *pairServiceDelete
	opt, err = parsePairServiceDelete(pairs)
	if err != nil {
		return
	}

	return s.delete(ctx, name, opt)
}

// Get will get a valid storager instance for service.
//
// This function will create a context by default.
func (s *Service) Get(name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	ctx := context.Background()
	return s.GetWithContext(ctx, name, pairs...)
}

// GetWithContext will get a valid storager instance for service.
func (s *Service) GetWithContext(ctx context.Context, name string, pairs ...*types.Pair) (store storage.Storager, err error) {
	defer func() {
		err = s.formatError(services.OpGet, err, name)
	}()
	var opt *pairServiceGet
	opt, err = parsePairServiceGet(pairs)
	if err != nil {
		return
	}

	return s.get(ctx, name, opt)
}

// List will list all storager instances under this service.
//
// This function will create a context by default.
func (s *Service) List(pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.ListWithContext(ctx, pairs...)
}

// ListWithContext will list all storager instances under this service.
func (s *Service) ListWithContext(ctx context.Context, pairs ...*types.Pair) (err error) {
	defer func() {

		err = s.formatError(services.OpList, err, "")
	}()
	var opt *pairServiceList
	opt, err = parsePairServiceList(pairs)
	if err != nil {
		return
	}

	return s.list(ctx, opt)
}

// pairStorageNewMap holds all available pairs
var pairStorageNewMap = map[string]struct{}{
	// Required pairs
	ps.Location: struct{}{},
	ps.Name:     struct{}{},
	// Optional pairs
	ps.WorkDir: struct{}{},
	// Generated pairs
	ps.HTTPClientOptions: struct{}{},
}

// pairStorageNew is the parsed struct
type pairStorageNew struct {
	pairs []*types.Pair

	// Required pairs
	Location string
	Name     string
	// Optional pairs
	HasWorkDir bool
	WorkDir    string
	// Generated pairs
	HasHTTPClientOptions bool
	HTTPClientOptions    *httpclient.Options
}

// parsePairStorageNew will parse *types.Pair slice into *pairStorageNew
func parsePairStorageNew(opts []*types.Pair) (*pairStorageNew, error) {
	result := &pairStorageNew{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
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
	v, ok = values[ps.Name]
	if !ok {
		return nil, services.NewPairRequiredError(ps.Name)
	}
	if ok {
		result.Name = v.(string)
	}
	// Handle optional pairs
	v, ok = values[ps.WorkDir]
	if ok {
		result.HasWorkDir = true
		result.WorkDir = v.(string)
	}
	// Handle generated pairs
	v, ok = values[ps.HTTPClientOptions]
	if ok {
		result.HasHTTPClientOptions = true
		result.HTTPClientOptions = v.(*httpclient.Options)
	}

	return result, nil
}

// pairStorageAbortSegmentMap holds all available pairs
var pairStorageAbortSegmentMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageAbortSegment is the parsed struct
type pairStorageAbortSegment struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageAbortSegment will parse *types.Pair slice into *pairStorageAbortSegment
func parsePairStorageAbortSegment(opts []*types.Pair) (*pairStorageAbortSegment, error) {
	result := &pairStorageAbortSegment{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageAbortSegmentMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairStorageCompleteSegmentMap holds all available pairs
var pairStorageCompleteSegmentMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageCompleteSegment is the parsed struct
type pairStorageCompleteSegment struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageCompleteSegment will parse *types.Pair slice into *pairStorageCompleteSegment
func parsePairStorageCompleteSegment(opts []*types.Pair) (*pairStorageCompleteSegment, error) {
	result := &pairStorageCompleteSegment{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageCompleteSegmentMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

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
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageDelete will parse *types.Pair slice into *pairStorageDelete
func parsePairStorageDelete(opts []*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageDeleteMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

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
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageInitIndexSegment will parse *types.Pair slice into *pairStorageInitIndexSegment
func parsePairStorageInitIndexSegment(opts []*types.Pair) (*pairStorageInitIndexSegment, error) {
	result := &pairStorageInitIndexSegment{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageInitIndexSegmentMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

	// Handle required pairs
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
	pairs []*types.Pair

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
	result := &pairStorageListDir{
		pairs: opts,
	}

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

// pairStorageListPrefixMap holds all available pairs
var pairStorageListPrefixMap = map[string]struct{}{
	// Required pairs
	ps.ObjectFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageListPrefix is the parsed struct
type pairStorageListPrefix struct {
	pairs []*types.Pair

	// Required pairs
	ObjectFunc types.ObjectFunc
	// Optional pairs
	// Generated pairs
}

// parsePairStorageListPrefix will parse *types.Pair slice into *pairStorageListPrefix
func parsePairStorageListPrefix(opts []*types.Pair) (*pairStorageListPrefix, error) {
	result := &pairStorageListPrefix{
		pairs: opts,
	}

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

// pairStorageListPrefixSegmentsMap holds all available pairs
var pairStorageListPrefixSegmentsMap = map[string]struct{}{
	// Required pairs
	ps.SegmentFunc: struct{}{},
	// Optional pairs
	// Generated pairs
}

// pairStorageListPrefixSegments is the parsed struct
type pairStorageListPrefixSegments struct {
	pairs []*types.Pair

	// Required pairs
	SegmentFunc segment.Func
	// Optional pairs
	// Generated pairs
}

// parsePairStorageListPrefixSegments will parse *types.Pair slice into *pairStorageListPrefixSegments
func parsePairStorageListPrefixSegments(opts []*types.Pair) (*pairStorageListPrefixSegments, error) {
	result := &pairStorageListPrefixSegments{
		pairs: opts,
	}

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

// pairStorageMetadataMap holds all available pairs
var pairStorageMetadataMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageMetadata is the parsed struct
type pairStorageMetadata struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMetadata will parse *types.Pair slice into *pairStorageMetadata
func parsePairStorageMetadata(opts []*types.Pair) (*pairStorageMetadata, error) {
	result := &pairStorageMetadata{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageMetadataMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

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
	ps.ReadCallbackFunc: struct{}{},
}

// pairStorageRead is the parsed struct
type pairStorageRead struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
}

// parsePairStorageRead will parse *types.Pair slice into *pairStorageRead
func parsePairStorageRead(opts []*types.Pair) (*pairStorageRead, error) {
	result := &pairStorageRead{
		pairs: opts,
	}

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
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}

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
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageStat will parse *types.Pair slice into *pairStorageStat
func parsePairStorageStat(opts []*types.Pair) (*pairStorageStat, error) {
	result := &pairStorageStat{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageStatMap[v.Key]; !ok {
			return nil, services.NewPairUnsupportedError(v)
		}
		values[v.Key] = v.Value
	}

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
	ps.ReadCallbackFunc: struct{}{},
}

// pairStorageWrite is the parsed struct
type pairStorageWrite struct {
	pairs []*types.Pair

	// Required pairs
	Size int64
	// Optional pairs
	HasChecksum     bool
	Checksum        string
	HasStorageClass bool
	StorageClass    string
	// Generated pairs
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
}

// parsePairStorageWrite will parse *types.Pair slice into *pairStorageWrite
func parsePairStorageWrite(opts []*types.Pair) (*pairStorageWrite, error) {
	result := &pairStorageWrite{
		pairs: opts,
	}

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
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}

	return result, nil
}

// pairStorageWriteIndexSegmentMap holds all available pairs
var pairStorageWriteIndexSegmentMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
	ps.ReadCallbackFunc: struct{}{},
}

// pairStorageWriteIndexSegment is the parsed struct
type pairStorageWriteIndexSegment struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
}

// parsePairStorageWriteIndexSegment will parse *types.Pair slice into *pairStorageWriteIndexSegment
func parsePairStorageWriteIndexSegment(opts []*types.Pair) (*pairStorageWriteIndexSegment, error) {
	result := &pairStorageWriteIndexSegment{
		pairs: opts,
	}

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
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}

	return result, nil
}

// AbortSegment will abort a segment.
//
// This function will create a context by default.
func (s *Storage) AbortSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.AbortSegmentWithContext(ctx, seg, pairs...)
}

// AbortSegmentWithContext will abort a segment.
func (s *Storage) AbortSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpAbortSegment, err, seg.Path(), seg.ID())
	}()
	var opt *pairStorageAbortSegment
	opt, err = parsePairStorageAbortSegment(pairs)
	if err != nil {
		return
	}

	return s.abortSegment(ctx, seg, opt)
}

// CompleteSegment will complete a segment and merge them into a File.
//
// This function will create a context by default.
func (s *Storage) CompleteSegment(seg segment.Segment, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.CompleteSegmentWithContext(ctx, seg, pairs...)
}

// CompleteSegmentWithContext will complete a segment and merge them into a File.
func (s *Storage) CompleteSegmentWithContext(ctx context.Context, seg segment.Segment, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpCompleteSegment, err, seg.Path(), seg.ID())
	}()
	var opt *pairStorageCompleteSegment
	opt, err = parsePairStorageCompleteSegment(pairs)
	if err != nil {
		return
	}

	return s.completeSegment(ctx, seg, opt)
}

// Delete will delete an Object from service.
//
// This function will create a context by default.
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.DeleteWithContext(ctx, path, pairs...)
}

// DeleteWithContext will delete an Object from service.
func (s *Storage) DeleteWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpDelete, err, path)
	}()
	var opt *pairStorageDelete
	opt, err = parsePairStorageDelete(pairs)
	if err != nil {
		return
	}

	return s.delete(ctx, path, opt)
}

// InitIndexSegment will init an index based segment.
//
// This function will create a context by default.
func (s *Storage) InitIndexSegment(path string, pairs ...*types.Pair) (seg segment.Segment, err error) {
	ctx := context.Background()
	return s.InitIndexSegmentWithContext(ctx, path, pairs...)
}

// InitIndexSegmentWithContext will init an index based segment.
func (s *Storage) InitIndexSegmentWithContext(ctx context.Context, path string, pairs ...*types.Pair) (seg segment.Segment, err error) {
	defer func() {
		err = s.formatError(services.OpInitIndexSegment, err, path)
	}()
	var opt *pairStorageInitIndexSegment
	opt, err = parsePairStorageInitIndexSegment(pairs)
	if err != nil {
		return
	}

	return s.initIndexSegment(ctx, path, opt)
}

// ListDir will return list a specific dir.
//
// This function will create a context by default.
func (s *Storage) ListDir(dir string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.ListDirWithContext(ctx, dir, pairs...)
}

// ListDirWithContext will return list a specific dir.
func (s *Storage) ListDirWithContext(ctx context.Context, dir string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListDir, err, dir)
	}()
	var opt *pairStorageListDir
	opt, err = parsePairStorageListDir(pairs)
	if err != nil {
		return
	}

	return s.listDir(ctx, dir, opt)
}

// ListPrefix will return list a specific dir.
//
// This function will create a context by default.
func (s *Storage) ListPrefix(prefix string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.ListPrefixWithContext(ctx, prefix, pairs...)
}

// ListPrefixWithContext will return list a specific dir.
func (s *Storage) ListPrefixWithContext(ctx context.Context, prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefix, err, prefix)
	}()
	var opt *pairStorageListPrefix
	opt, err = parsePairStorageListPrefix(pairs)
	if err != nil {
		return
	}

	return s.listPrefix(ctx, prefix, opt)
}

// ListPrefixSegments will list segments.
//
// This function will create a context by default.
func (s *Storage) ListPrefixSegments(prefix string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.ListPrefixSegmentsWithContext(ctx, prefix, pairs...)
}

// ListPrefixSegmentsWithContext will list segments.
func (s *Storage) ListPrefixSegmentsWithContext(ctx context.Context, prefix string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpListPrefixSegments, err, prefix)
	}()
	var opt *pairStorageListPrefixSegments
	opt, err = parsePairStorageListPrefixSegments(pairs)
	if err != nil {
		return
	}

	return s.listPrefixSegments(ctx, prefix, opt)
}

// Metadata will return current storager's metadata.
//
// This function will create a context by default.
func (s *Storage) Metadata(pairs ...*types.Pair) (meta info.StorageMeta, err error) {
	ctx := context.Background()
	return s.MetadataWithContext(ctx, pairs...)
}

// MetadataWithContext will return current storager's metadata.
func (s *Storage) MetadataWithContext(ctx context.Context, pairs ...*types.Pair) (meta info.StorageMeta, err error) {
	defer func() {
		err = s.formatError(services.OpMetadata, err)
	}()
	var opt *pairStorageMetadata
	opt, err = parsePairStorageMetadata(pairs)
	if err != nil {
		return
	}

	return s.metadata(ctx, opt)
}

// Read will read the file's data.
//
// This function will create a context by default.
func (s *Storage) Read(path string, pairs ...*types.Pair) (rc io.ReadCloser, err error) {
	ctx := context.Background()
	return s.ReadWithContext(ctx, path, pairs...)
}

// ReadWithContext will read the file's data.
func (s *Storage) ReadWithContext(ctx context.Context, path string, pairs ...*types.Pair) (rc io.ReadCloser, err error) {
	defer func() {
		err = s.formatError(services.OpRead, err, path)
	}()
	var opt *pairStorageRead
	opt, err = parsePairStorageRead(pairs)
	if err != nil {
		return
	}

	return s.read(ctx, path, opt)
}

// Stat will stat a path to get info of an object.
//
// This function will create a context by default.
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	ctx := context.Background()
	return s.StatWithContext(ctx, path, pairs...)
}

// StatWithContext will stat a path to get info of an object.
func (s *Storage) StatWithContext(ctx context.Context, path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError(services.OpStat, err, path)
	}()
	var opt *pairStorageStat
	opt, err = parsePairStorageStat(pairs)
	if err != nil {
		return
	}

	return s.stat(ctx, path, opt)
}

// Write will write data into a file.
//
// This function will create a context by default.
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.WriteWithContext(ctx, path, r, pairs...)
}

// WriteWithContext will write data into a file.
func (s *Storage) WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWrite, err, path)
	}()
	var opt *pairStorageWrite
	opt, err = parsePairStorageWrite(pairs)
	if err != nil {
		return
	}

	return s.write(ctx, path, r, opt)
}

// WriteIndexSegment will write a part into an index based segment.
//
// This function will create a context by default.
func (s *Storage) WriteIndexSegment(seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.WriteIndexSegmentWithContext(ctx, seg, r, index, size, pairs...)
}

// WriteIndexSegmentWithContext will write a part into an index based segment.
func (s *Storage) WriteIndexSegmentWithContext(ctx context.Context, seg segment.Segment, r io.Reader, index int, size int64, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpWriteIndexSegment, err, seg.Path(), seg.ID())
	}()
	var opt *pairStorageWriteIndexSegment
	opt, err = parsePairStorageWriteIndexSegment(pairs)
	if err != nil {
		return
	}

	return s.writeIndexSegment(ctx, seg, r, index, size, opt)
}
