// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package gcs

import (
	"context"
	"io"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	"github.com/aos-dev/go-storage/v2/pkg/httpclient"
	"github.com/aos-dev/go-storage/v2/pkg/segment"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
	ps "github.com/aos-dev/go-storage/v2/types/pairs"
)

var _ credential.Provider
var _ endpoint.Provider
var _ segment.Segment
var _ storage.Storager
var _ services.ServiceError
var _ httpclient.Options

// Type is the type for gcs
const Type = "gcs"

// Service available pairs.
const (
	// StorageClass will // StorageClass
	PairStorageClass = "gcs_storage_class"
)

// Service available infos.
const (
	InfoObjectMetaStorageClass = "gcs-storage-class"
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
	ps.Project:    struct{}{},
	// Optional pairs
	// Generated pairs
	ps.HTTPClientOptions: struct{}{},
}

// pairServiceNew is the parsed struct
type pairServiceNew struct {
	pairs []*types.Pair

	// Required pairs
	Credential *credential.Provider
	Project    string
	// Optional pairs
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
	v, ok = values[ps.Project]
	if !ok {
		return nil, services.NewPairRequiredError(ps.Project)
	}
	if ok {
		result.Project = v.(string)
	}
	// Handle optional pairs
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
	// Optional pairs
	// Generated pairs
}

// pairServiceCreate is the parsed struct
type pairServiceCreate struct {
	pairs []*types.Pair

	// Required pairs
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

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairServiceDeleteMap holds all available pairs
var pairServiceDeleteMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairServiceDelete is the parsed struct
type pairServiceDelete struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
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

	// Handle required pairs
	// Handle optional pairs
	// Handle generated pairs

	return result, nil
}

// pairServiceGetMap holds all available pairs
var pairServiceGetMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairServiceGet is the parsed struct
type pairServiceGet struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
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

	// Handle required pairs
	// Handle optional pairs
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
	ps.Name: struct{}{},
	// Optional pairs
	ps.WorkDir: struct{}{},
	// Generated pairs
	ps.HTTPClientOptions: struct{}{},
}

// pairStorageNew is the parsed struct
type pairStorageNew struct {
	pairs []*types.Pair

	// Required pairs
	Name string
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
