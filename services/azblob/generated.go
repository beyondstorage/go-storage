// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package azblob

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
	ps "github.com/Xuanwo/storage/types/pairs"
)

var _ credential.Provider
var _ endpoint.Provider
var _ segment.Segment
var _ storage.Storager
var _ storageclass.Type
var _ services.ServiceError

// Type is the type for azblob
const Type = "azblob"

type pairServiceCreate struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
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
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairServiceDelete struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
}

func parseServicePairDelete(opts ...*types.Pair) (*pairServiceDelete, error) {
	result := &pairServiceDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairServiceGet struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
}

func parseServicePairGet(opts ...*types.Pair) (*pairServiceGet, error) {
	result := &pairServiceGet{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairServiceList struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	StoragerFunc storage.StoragerFunc
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
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[ps.StoragerFunc]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.StoragerFunc,
			Value: nil,
		}
	}
	if ok {
		result.StoragerFunc = v.(storage.StoragerFunc)
	}
	return result, nil
}

type pairServiceNew struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	Credential *credential.Provider
	Endpoint   endpoint.Provider
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
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[ps.Credential]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.Credential,
			Value: nil,
		}
	}
	if ok {
		result.Credential = v.(*credential.Provider)
	}
	v, ok = values[ps.Endpoint]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.Endpoint,
			Value: nil,
		}
	}
	if ok {
		result.Endpoint = v.(endpoint.Provider)
	}
	return result, nil
}

type pairStorageDelete struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
}

func parseStoragePairDelete(opts ...*types.Pair) (*pairStorageDelete, error) {
	result := &pairStorageDelete{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairStorageList struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	FileFunc types.ObjectFunc
}

func parseStoragePairList(opts ...*types.Pair) (*pairStorageList, error) {
	result := &pairStorageList{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[ps.FileFunc]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.FileFunc,
			Value: nil,
		}
	}
	if ok {
		result.FileFunc = v.(types.ObjectFunc)
	}
	return result, nil
}

type pairStorageMetadata struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
}

func parseStoragePairMetadata(opts ...*types.Pair) (*pairStorageMetadata, error) {
	result := &pairStorageMetadata{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairStorageNew struct {
	// Pre-defined pairs

	// Meta-defined pairs
	Name       string
	HasWorkDir bool
	WorkDir    string
}

func parseStoragePairNew(opts ...*types.Pair) (*pairStorageNew, error) {
	result := &pairStorageNew{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs

	// Parse meta-defined pairs
	v, ok = values[ps.Name]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.Name,
			Value: nil,
		}
	}
	if ok {
		result.Name = v.(string)
	}
	v, ok = values[ps.WorkDir]
	if ok {
		result.HasWorkDir = true
		result.WorkDir = v.(string)
	}
	return result, nil
}

type pairStorageRead struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
}

func parseStoragePairRead(opts ...*types.Pair) (*pairStorageRead, error) {
	result := &pairStorageRead{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}
	return result, nil
}

type pairStorageStat struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
}

func parseStoragePairStat(opts ...*types.Pair) (*pairStorageStat, error) {
	result := &pairStorageStat{}

	values := make(map[string]interface{})
	for _, v := range opts {
		values[v.Key] = v.Value
	}

	var v interface{}
	var ok bool

	// Parse pre-defined pairs
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	return result, nil
}

type pairStorageWrite struct {
	// Pre-defined pairs
	Context context.Context

	// Meta-defined pairs
	HasChecksum         bool
	Checksum            string
	HasReadCallbackFunc bool
	ReadCallbackFunc    func([]byte)
	Size                int64
	HasStorageClass     bool
	StorageClass        storageclass.Type
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
	v, ok = values[ps.Context]
	if ok {
		result.Context = v.(context.Context)
	} else {
		result.Context = context.Background()
	}

	// Parse meta-defined pairs
	v, ok = values[ps.Checksum]
	if ok {
		result.HasChecksum = true
		result.Checksum = v.(string)
	}
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}
	v, ok = values[ps.Size]
	if !ok {
		return nil, &services.PairError{
			Op:    "parse",
			Err:   services.ErrPairRequired,
			Key:   ps.Size,
			Value: nil,
		}
	}
	if ok {
		result.Size = v.(int64)
	}
	v, ok = values[ps.StorageClass]
	if ok {
		result.HasStorageClass = true
		result.StorageClass = v.(storageclass.Type)
	}
	return result, nil
}

// CreateWithContext adds context support for Create.
func (s *Service) CreateWithContext(ctx context.Context, name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.service.Create")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Create(name, pairs...)
}

// DeleteWithContext adds context support for Delete.
func (s *Service) DeleteWithContext(ctx context.Context, name string, pairs ...*types.Pair) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.service.Delete")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Delete(name, pairs...)
}

// GetWithContext adds context support for Get.
func (s *Service) GetWithContext(ctx context.Context, name string, pairs ...*types.Pair) (st storage.Storager, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.service.Get")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Get(name, pairs...)
}

// ListWithContext adds context support for List.
func (s *Service) ListWithContext(ctx context.Context, pairs ...*types.Pair) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.service.List")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.List(pairs...)
}

// DeleteWithContext adds context support for Delete.
func (s *Storage) DeleteWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.Delete")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Delete(path, pairs...)
}

// ListWithContext adds context support for List.
func (s *Storage) ListWithContext(ctx context.Context, path string, pairs ...*types.Pair) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.List")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.List(path, pairs...)
}

// MetadataWithContext adds context support for Metadata.
func (s *Storage) MetadataWithContext(ctx context.Context, pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.Metadata")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Metadata(pairs...)
}

// ReadWithContext adds context support for Read.
func (s *Storage) ReadWithContext(ctx context.Context, path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.Read")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Read(path, pairs...)
}

// StatWithContext adds context support for Stat.
func (s *Storage) StatWithContext(ctx context.Context, path string, pairs ...*types.Pair) (o *types.Object, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.Stat")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Stat(path, pairs...)
}

// WriteWithContext adds context support for Write.
func (s *Storage) WriteWithContext(ctx context.Context, path string, r io.Reader, pairs ...*types.Pair) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "github.com/Xuanwo/storage/services/azblob.storage.Write")
	defer span.Finish()

	pairs = append(pairs, ps.WithContext(ctx))
	return s.Write(path, r, pairs...)
}
