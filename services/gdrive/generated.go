// Code generated by go generate via cmd/definitions; DO NOT EDIT.
package gdrive

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

var (
	_ types.Storager
	_ services.ServiceError
	_ strings.Reader
	_ time.Duration
	_ http.Request
)

// Type is the type for gdrive
const Type = "gdrive"

// ObjectSystemMetadata stores system metadata for object.
type ObjectSystemMetadata struct {
}

// GetObjectSystemMetadata will get ObjectSystemMetadata from Object.
//
// - This function should not be called by service implementer.
// - The returning ObjectServiceMetadata is read only and should not be modified.
func GetObjectSystemMetadata(o *types.Object) ObjectSystemMetadata {
	sm, ok := o.GetSystemMetadata()
	if ok {
		return sm.(ObjectSystemMetadata)
	}
	return ObjectSystemMetadata{}
}

// setObjectSystemMetadata will set ObjectSystemMetadata into Object.
//
// - This function should only be called once, please make sure all data has been written before set.
func setObjectSystemMetadata(o *types.Object, sm ObjectSystemMetadata) {
	o.SetSystemMetadata(sm)
}

// StorageSystemMetadata stores system metadata for object.
type StorageSystemMetadata struct {
}

// GetStorageSystemMetadata will get StorageSystemMetadata from Storage.
//
// - This function should not be called by service implementer.
// - The returning StorageServiceMetadata is read only and should not be modified.
func GetStorageSystemMetadata(s *types.StorageMeta) StorageSystemMetadata {
	sm, ok := s.GetSystemMetadata()
	if ok {
		return sm.(StorageSystemMetadata)
	}
	return StorageSystemMetadata{}
}

// setStorageSystemMetadata will set StorageSystemMetadata into Storage.
//
// - This function should only be called once, please make sure all data has been written before set.
func setStorageSystemMetadata(s *types.StorageMeta, sm StorageSystemMetadata) {
	s.SetSystemMetadata(sm)
}

type Factory struct {
	Credential string
	Name       string
	WorkDir    string
}

func (f *Factory) FromString(conn string) (err error) {
	slash := strings.IndexByte(conn, '/')
	question := strings.IndexByte(conn, '?')

	var partService, partStorage, partParams string

	if question != -1 {
		if len(conn) > question {
			partParams = conn[question+1:]
		}
		conn = conn[:question]
	}

	if slash != -1 {
		partService = conn[:slash]
		partStorage = conn[slash:]
	} else {
		partService = conn
	}

	if partService != "" {
		f.Credential = partService
	}
	if partStorage != "" {
		slash := strings.IndexByte(partStorage[1:], '/')
		if slash == -1 {
			f.Name = partStorage[1:]
		} else {
			f.Name, f.WorkDir = partStorage[1:slash+1], partStorage[slash+1:]
		}

	}
	if partParams != "" {
		xs := strings.Split(partParams, "&")
		for _, v := range xs {
			var key, value string
			vs := strings.SplitN(v, "=", 2)
			key = vs[0]
			if len(vs) > 1 {
				value = vs[1]
			}
			switch key {
			case "credential":
				f.Credential = value
			case "name":
				f.Name = value
			case "work_dir":
				f.WorkDir = value
			}
		}
	}
	return nil
}
func (f *Factory) WithPairs(ps ...types.Pair) (err error) {
	for _, v := range ps {
		switch v.Key {
		case "credential":
			f.Credential = v.Value.(string)
		case "name":
			f.Name = v.Value.(string)
		case "work_dir":
			f.WorkDir = v.Value.(string)
		}
	}
	return nil
}
func (f *Factory) FromMap(m map[string]interface{}) (err error) {
	return errors.New("FromMap not implemented")
}
func (f *Factory) NewServicer() (srv types.Servicer, err error) {
	return f.newService()
}
func (f *Factory) NewStorager() (sto types.Storager, err error) {
	return f.newStorage()
}
func (f *Factory) serviceFeatures() (s types.ServiceFeatures) {
	return
}
func (f *Factory) storageFeatures() (s types.StorageFeatures) {
	s.Copy = true
	s.Create = true
	s.CreateDir = true
	s.Delete = true
	s.List = true
	s.Metadata = true
	s.Read = true
	s.Stat = true
	s.Write = true
	s.WriteEmptyObject = true
	return
}

var _ types.Servicer = &Service{}

// Deprecated: Use types.ServiceFeatures instead.
type ServiceFeatures = types.ServiceFeatures

// Deprecated: Use types.DefaultServicePairs instead.
type DefaultServicePairs = types.DefaultServicePairs

func (s *Service) Features() types.ServiceFeatures {
	return s.features
}

type pairServiceCreate struct {
	pairs []types.Pair
}

func (s *Service) parsePairServiceCreate(opts []types.Pair) (pairServiceCreate, error) {
	result :=
		pairServiceCreate{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairServiceCreate{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Service) Create(name string, pairs ...types.Pair) (store types.Storager, err error) {
	err = types.NewOperationNotImplementedError("create")
	return
}
func (s *Service) CreateWithContext(ctx context.Context, name string, pairs ...types.Pair) (store types.Storager, err error) {
	err = types.NewOperationNotImplementedError("create")
	return
}

type pairServiceDelete struct {
	pairs []types.Pair
}

func (s *Service) parsePairServiceDelete(opts []types.Pair) (pairServiceDelete, error) {
	result :=
		pairServiceDelete{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairServiceDelete{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Service) Delete(name string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("delete")
	return
}
func (s *Service) DeleteWithContext(ctx context.Context, name string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("delete")
	return
}

type pairServiceGet struct {
	pairs []types.Pair
}

func (s *Service) parsePairServiceGet(opts []types.Pair) (pairServiceGet, error) {
	result :=
		pairServiceGet{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairServiceGet{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Service) Get(name string, pairs ...types.Pair) (store types.Storager, err error) {
	err = types.NewOperationNotImplementedError("get")
	return
}
func (s *Service) GetWithContext(ctx context.Context, name string, pairs ...types.Pair) (store types.Storager, err error) {
	err = types.NewOperationNotImplementedError("get")
	return
}

type pairServiceList struct {
	pairs []types.Pair
}

func (s *Service) parsePairServiceList(opts []types.Pair) (pairServiceList, error) {
	result :=
		pairServiceList{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairServiceList{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Service) List(pairs ...types.Pair) (sti *types.StoragerIterator, err error) {
	err = types.NewOperationNotImplementedError("list")
	return
}
func (s *Service) ListWithContext(ctx context.Context, pairs ...types.Pair) (sti *types.StoragerIterator, err error) {
	err = types.NewOperationNotImplementedError("list")
	return
}

var _ types.Storager = &Storage{}

// Deprecated: Use types.StorageFeatures instead.
type StorageFeatures = types.StorageFeatures

// Deprecated: Use types.DefaultStoragePairs instead.
type DefaultStoragePairs = types.DefaultStoragePairs

func (s *Storage) Features() types.StorageFeatures {
	return s.features
}

type pairStorageCombineBlock struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCombineBlock(opts []types.Pair) (pairStorageCombineBlock, error) {
	result :=
		pairStorageCombineBlock{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCombineBlock{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CombineBlock(o *types.Object, bids []string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("combine_block")
	return
}
func (s *Storage) CombineBlockWithContext(ctx context.Context, o *types.Object, bids []string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("combine_block")
	return
}

type pairStorageCommitAppend struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCommitAppend(opts []types.Pair) (pairStorageCommitAppend, error) {
	result :=
		pairStorageCommitAppend{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCommitAppend{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CommitAppend(o *types.Object, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("commit_append")
	return
}
func (s *Storage) CommitAppendWithContext(ctx context.Context, o *types.Object, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("commit_append")
	return
}

type pairStorageCompleteMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCompleteMultipart(opts []types.Pair) (pairStorageCompleteMultipart, error) {
	result :=
		pairStorageCompleteMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCompleteMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CompleteMultipart(o *types.Object, parts []*types.Part, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("complete_multipart")
	return
}
func (s *Storage) CompleteMultipartWithContext(ctx context.Context, o *types.Object, parts []*types.Part, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("complete_multipart")
	return
}

type pairStorageCopy struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCopy(opts []types.Pair) (pairStorageCopy, error) {
	result :=
		pairStorageCopy{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCopy{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Copy(src string, dst string, pairs ...types.Pair) (err error) {
	ctx := context.Background()
	return s.CopyWithContext(ctx, src, dst, pairs...)
}
func (s *Storage) CopyWithContext(ctx context.Context, src string, dst string, pairs ...types.Pair) (err error) {
	defer func() {
		err =
			s.formatError("copy", err, src, dst)
	}()
	pairs = append(pairs, s.defaultPairs.Copy...)
	var opt pairStorageCopy

	opt, err = s.parsePairStorageCopy(pairs)
	if err != nil {
		return
	}
	return s.copy(ctx, strings.ReplaceAll(src, "\\", "/"), strings.ReplaceAll(dst, "\\", "/"), opt)
}

type pairStorageCreate struct {
	pairs         []types.Pair
	HasObjectMode bool
	ObjectMode    types.ObjectMode
}

func (s *Storage) parsePairStorageCreate(opts []types.Pair) (pairStorageCreate, error) {
	result :=
		pairStorageCreate{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "object_mode":
			if result.HasObjectMode {
				continue
			}
			result.HasObjectMode = true
			result.ObjectMode = v.Value.(types.ObjectMode)
		default:
			return pairStorageCreate{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Create(path string, pairs ...types.Pair) (o *types.Object) {
	pairs = append(pairs, s.defaultPairs.Create...)
	var opt pairStorageCreate

	// Ignore error while handling local functions.
	opt, _ = s.parsePairStorageCreate(pairs)
	return s.create(path, opt)
}

type pairStorageCreateAppend struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreateAppend(opts []types.Pair) (pairStorageCreateAppend, error) {
	result :=
		pairStorageCreateAppend{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreateAppend{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreateAppend(path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_append")
	return
}
func (s *Storage) CreateAppendWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_append")
	return
}

type pairStorageCreateBlock struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreateBlock(opts []types.Pair) (pairStorageCreateBlock, error) {
	result :=
		pairStorageCreateBlock{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreateBlock{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreateBlock(path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_block")
	return
}
func (s *Storage) CreateBlockWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_block")
	return
}

type pairStorageCreateDir struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreateDir(opts []types.Pair) (pairStorageCreateDir, error) {
	result :=
		pairStorageCreateDir{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreateDir{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreateDir(path string, pairs ...types.Pair) (o *types.Object, err error) {
	ctx := context.Background()
	return s.CreateDirWithContext(ctx, path, pairs...)
}
func (s *Storage) CreateDirWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	defer func() {
		err =
			s.formatError("create_dir", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.CreateDir...)
	var opt pairStorageCreateDir

	opt, err = s.parsePairStorageCreateDir(pairs)
	if err != nil {
		return
	}
	return s.createDir(ctx, strings.ReplaceAll(path, "\\", "/"), opt)
}

type pairStorageCreateLink struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreateLink(opts []types.Pair) (pairStorageCreateLink, error) {
	result :=
		pairStorageCreateLink{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreateLink{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreateLink(path string, target string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_link")
	return
}
func (s *Storage) CreateLinkWithContext(ctx context.Context, path string, target string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_link")
	return
}

type pairStorageCreateMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreateMultipart(opts []types.Pair) (pairStorageCreateMultipart, error) {
	result :=
		pairStorageCreateMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreateMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreateMultipart(path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_multipart")
	return
}
func (s *Storage) CreateMultipartWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_multipart")
	return
}

type pairStorageCreatePage struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageCreatePage(opts []types.Pair) (pairStorageCreatePage, error) {
	result :=
		pairStorageCreatePage{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageCreatePage{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) CreatePage(path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_page")
	return
}
func (s *Storage) CreatePageWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	err = types.NewOperationNotImplementedError("create_page")
	return
}

type pairStorageDelete struct {
	pairs         []types.Pair
	HasObjectMode bool
	ObjectMode    types.ObjectMode
}

func (s *Storage) parsePairStorageDelete(opts []types.Pair) (pairStorageDelete, error) {
	result :=
		pairStorageDelete{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "object_mode":
			if result.HasObjectMode {
				continue
			}
			result.HasObjectMode = true
			result.ObjectMode = v.Value.(types.ObjectMode)
		default:
			return pairStorageDelete{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Delete(path string, pairs ...types.Pair) (err error) {
	ctx := context.Background()
	return s.DeleteWithContext(ctx, path, pairs...)
}
func (s *Storage) DeleteWithContext(ctx context.Context, path string, pairs ...types.Pair) (err error) {
	defer func() {
		err =
			s.formatError("delete", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.Delete...)
	var opt pairStorageDelete

	opt, err = s.parsePairStorageDelete(pairs)
	if err != nil {
		return
	}
	return s.delete(ctx, strings.ReplaceAll(path, "\\", "/"), opt)
}

type pairStorageFetch struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageFetch(opts []types.Pair) (pairStorageFetch, error) {
	result :=
		pairStorageFetch{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageFetch{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Fetch(path string, url string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("fetch")
	return
}
func (s *Storage) FetchWithContext(ctx context.Context, path string, url string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("fetch")
	return
}

type pairStorageList struct {
	pairs       []types.Pair
	HasListMode bool
	ListMode    types.ListMode
}

func (s *Storage) parsePairStorageList(opts []types.Pair) (pairStorageList, error) {
	result :=
		pairStorageList{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "list_mode":
			if result.HasListMode {
				continue
			}
			result.HasListMode = true
			result.ListMode = v.Value.(types.ListMode)
		default:
			return pairStorageList{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) List(path string, pairs ...types.Pair) (oi *types.ObjectIterator, err error) {
	ctx := context.Background()
	return s.ListWithContext(ctx, path, pairs...)
}
func (s *Storage) ListWithContext(ctx context.Context, path string, pairs ...types.Pair) (oi *types.ObjectIterator, err error) {
	defer func() {
		err =
			s.formatError("list", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.List...)
	var opt pairStorageList

	opt, err = s.parsePairStorageList(pairs)
	if err != nil {
		return
	}
	return s.list(ctx, strings.ReplaceAll(path, "\\", "/"), opt)
}

type pairStorageListBlock struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageListBlock(opts []types.Pair) (pairStorageListBlock, error) {
	result :=
		pairStorageListBlock{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageListBlock{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) ListBlock(o *types.Object, pairs ...types.Pair) (bi *types.BlockIterator, err error) {
	err = types.NewOperationNotImplementedError("list_block")
	return
}
func (s *Storage) ListBlockWithContext(ctx context.Context, o *types.Object, pairs ...types.Pair) (bi *types.BlockIterator, err error) {
	err = types.NewOperationNotImplementedError("list_block")
	return
}

type pairStorageListMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageListMultipart(opts []types.Pair) (pairStorageListMultipart, error) {
	result :=
		pairStorageListMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageListMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) ListMultipart(o *types.Object, pairs ...types.Pair) (pi *types.PartIterator, err error) {
	err = types.NewOperationNotImplementedError("list_multipart")
	return
}
func (s *Storage) ListMultipartWithContext(ctx context.Context, o *types.Object, pairs ...types.Pair) (pi *types.PartIterator, err error) {
	err = types.NewOperationNotImplementedError("list_multipart")
	return
}

type pairStorageMetadata struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageMetadata(opts []types.Pair) (pairStorageMetadata, error) {
	result :=
		pairStorageMetadata{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageMetadata{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Metadata(pairs ...types.Pair) (meta *types.StorageMeta) {
	pairs = append(pairs, s.defaultPairs.Metadata...)
	var opt pairStorageMetadata

	// Ignore error while handling local functions.
	opt, _ = s.parsePairStorageMetadata(pairs)
	return s.metadata(opt)
}

type pairStorageMove struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageMove(opts []types.Pair) (pairStorageMove, error) {
	result :=
		pairStorageMove{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageMove{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Move(src string, dst string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("move")
	return
}
func (s *Storage) MoveWithContext(ctx context.Context, src string, dst string, pairs ...types.Pair) (err error) {
	err = types.NewOperationNotImplementedError("move")
	return
}

type pairStorageQuerySignHTTPCompleteMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPCompleteMultipart(opts []types.Pair) (pairStorageQuerySignHTTPCompleteMultipart, error) {
	result :=
		pairStorageQuerySignHTTPCompleteMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPCompleteMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPCompleteMultipart(o *types.Object, parts []*types.Part, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_complete_multipart")
	return
}
func (s *Storage) QuerySignHTTPCompleteMultipartWithContext(ctx context.Context, o *types.Object, parts []*types.Part, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_complete_multipart")
	return
}

type pairStorageQuerySignHTTPCreateMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPCreateMultipart(opts []types.Pair) (pairStorageQuerySignHTTPCreateMultipart, error) {
	result :=
		pairStorageQuerySignHTTPCreateMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPCreateMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPCreateMultipart(path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_create_multipart")
	return
}
func (s *Storage) QuerySignHTTPCreateMultipartWithContext(ctx context.Context, path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_create_multipart")
	return
}

type pairStorageQuerySignHTTPDelete struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPDelete(opts []types.Pair) (pairStorageQuerySignHTTPDelete, error) {
	result :=
		pairStorageQuerySignHTTPDelete{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPDelete{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPDelete(path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_delete")
	return
}
func (s *Storage) QuerySignHTTPDeleteWithContext(ctx context.Context, path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_delete")
	return
}

type pairStorageQuerySignHTTPListMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPListMultipart(opts []types.Pair) (pairStorageQuerySignHTTPListMultipart, error) {
	result :=
		pairStorageQuerySignHTTPListMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPListMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPListMultipart(o *types.Object, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_list_multipart")
	return
}
func (s *Storage) QuerySignHTTPListMultipartWithContext(ctx context.Context, o *types.Object, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_list_multipart")
	return
}

type pairStorageQuerySignHTTPRead struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPRead(opts []types.Pair) (pairStorageQuerySignHTTPRead, error) {
	result :=
		pairStorageQuerySignHTTPRead{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPRead{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPRead(path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_read")
	return
}
func (s *Storage) QuerySignHTTPReadWithContext(ctx context.Context, path string, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_read")
	return
}

type pairStorageQuerySignHTTPWrite struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPWrite(opts []types.Pair) (pairStorageQuerySignHTTPWrite, error) {
	result :=
		pairStorageQuerySignHTTPWrite{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPWrite{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPWrite(path string, size int64, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_write")
	return
}
func (s *Storage) QuerySignHTTPWriteWithContext(ctx context.Context, path string, size int64, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_write")
	return
}

type pairStorageQuerySignHTTPWriteMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageQuerySignHTTPWriteMultipart(opts []types.Pair) (pairStorageQuerySignHTTPWriteMultipart, error) {
	result :=
		pairStorageQuerySignHTTPWriteMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageQuerySignHTTPWriteMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) QuerySignHTTPWriteMultipart(o *types.Object, size int64, index int, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_write_multipart")
	return
}
func (s *Storage) QuerySignHTTPWriteMultipartWithContext(ctx context.Context, o *types.Object, size int64, index int, expire time.Duration, pairs ...types.Pair) (req *http.Request, err error) {
	err = types.NewOperationNotImplementedError("query_sign_http_write_multipart")
	return
}

type pairStorageRead struct {
	pairs         []types.Pair
	HasIoCallback bool
	IoCallback    func([]byte)
	HasOffset     bool
	Offset        int64
	HasSize       bool
	Size          int64
}

func (s *Storage) parsePairStorageRead(opts []types.Pair) (pairStorageRead, error) {
	result :=
		pairStorageRead{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "io_callback":
			if result.HasIoCallback {
				continue
			}
			result.HasIoCallback = true
			result.IoCallback = v.Value.(func([]byte))
		case "offset":
			if result.HasOffset {
				continue
			}
			result.HasOffset = true
			result.Offset = v.Value.(int64)
		case "size":
			if result.HasSize {
				continue
			}
			result.HasSize = true
			result.Size = v.Value.(int64)
		default:
			return pairStorageRead{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Read(path string, w io.Writer, pairs ...types.Pair) (n int64, err error) {
	ctx := context.Background()
	return s.ReadWithContext(ctx, path, w, pairs...)
}
func (s *Storage) ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...types.Pair) (n int64, err error) {
	defer func() {
		err =
			s.formatError("read", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.Read...)
	var opt pairStorageRead

	opt, err = s.parsePairStorageRead(pairs)
	if err != nil {
		return
	}
	return s.read(ctx, strings.ReplaceAll(path, "\\", "/"), w, opt)
}

type pairStorageStat struct {
	pairs         []types.Pair
	HasObjectMode bool
	ObjectMode    types.ObjectMode
}

func (s *Storage) parsePairStorageStat(opts []types.Pair) (pairStorageStat, error) {
	result :=
		pairStorageStat{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "object_mode":
			if result.HasObjectMode {
				continue
			}
			result.HasObjectMode = true
			result.ObjectMode = v.Value.(types.ObjectMode)
		default:
			return pairStorageStat{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Stat(path string, pairs ...types.Pair) (o *types.Object, err error) {
	ctx := context.Background()
	return s.StatWithContext(ctx, path, pairs...)
}
func (s *Storage) StatWithContext(ctx context.Context, path string, pairs ...types.Pair) (o *types.Object, err error) {
	defer func() {
		err =
			s.formatError("stat", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.Stat...)
	var opt pairStorageStat

	opt, err = s.parsePairStorageStat(pairs)
	if err != nil {
		return
	}
	return s.stat(ctx, strings.ReplaceAll(path, "\\", "/"), opt)
}

type pairStorageWrite struct {
	pairs          []types.Pair
	HasContentMd5  bool
	ContentMd5     string
	HasContentType bool
	ContentType    string
	HasIoCallback  bool
	IoCallback     func([]byte)
}

func (s *Storage) parsePairStorageWrite(opts []types.Pair) (pairStorageWrite, error) {
	result :=
		pairStorageWrite{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		case "content_md5":
			if result.HasContentMd5 {
				continue
			}
			result.HasContentMd5 = true
			result.ContentMd5 = v.Value.(string)
		case "content_type":
			if result.HasContentType {
				continue
			}
			result.HasContentType = true
			result.ContentType = v.Value.(string)
		case "io_callback":
			if result.HasIoCallback {
				continue
			}
			result.HasIoCallback = true
			result.IoCallback = v.Value.(func([]byte))
		default:
			return pairStorageWrite{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) Write(path string, r io.Reader, size int64, pairs ...types.Pair) (n int64, err error) {
	ctx := context.Background()
	return s.WriteWithContext(ctx, path, r, size, pairs...)
}
func (s *Storage) WriteWithContext(ctx context.Context, path string, r io.Reader, size int64, pairs ...types.Pair) (n int64, err error) {
	defer func() {
		err =
			s.formatError("write", err, path)
	}()
	pairs = append(pairs, s.defaultPairs.Write...)
	var opt pairStorageWrite

	opt, err = s.parsePairStorageWrite(pairs)
	if err != nil {
		return
	}
	return s.write(ctx, strings.ReplaceAll(path, "\\", "/"), r, size, opt)
}

type pairStorageWriteAppend struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageWriteAppend(opts []types.Pair) (pairStorageWriteAppend, error) {
	result :=
		pairStorageWriteAppend{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageWriteAppend{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) WriteAppend(o *types.Object, r io.Reader, size int64, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_append")
	return
}
func (s *Storage) WriteAppendWithContext(ctx context.Context, o *types.Object, r io.Reader, size int64, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_append")
	return
}

type pairStorageWriteBlock struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageWriteBlock(opts []types.Pair) (pairStorageWriteBlock, error) {
	result :=
		pairStorageWriteBlock{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageWriteBlock{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) WriteBlock(o *types.Object, r io.Reader, size int64, bid string, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_block")
	return
}
func (s *Storage) WriteBlockWithContext(ctx context.Context, o *types.Object, r io.Reader, size int64, bid string, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_block")
	return
}

type pairStorageWriteMultipart struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageWriteMultipart(opts []types.Pair) (pairStorageWriteMultipart, error) {
	result :=
		pairStorageWriteMultipart{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageWriteMultipart{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) WriteMultipart(o *types.Object, r io.Reader, size int64, index int, pairs ...types.Pair) (n int64, part *types.Part, err error) {
	err = types.NewOperationNotImplementedError("write_multipart")
	return
}
func (s *Storage) WriteMultipartWithContext(ctx context.Context, o *types.Object, r io.Reader, size int64, index int, pairs ...types.Pair) (n int64, part *types.Part, err error) {
	err = types.NewOperationNotImplementedError("write_multipart")
	return
}

type pairStorageWritePage struct {
	pairs []types.Pair
}

func (s *Storage) parsePairStorageWritePage(opts []types.Pair) (pairStorageWritePage, error) {
	result :=
		pairStorageWritePage{pairs: opts}

	for _, v := range opts {
		switch v.Key {
		default:
			return pairStorageWritePage{}, services.PairUnsupportedError{Pair: v}
		}
	}
	return result, nil
}
func (s *Storage) WritePage(o *types.Object, r io.Reader, size int64, offset int64, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_page")
	return
}
func (s *Storage) WritePageWithContext(ctx context.Context, o *types.Object, r io.Reader, size int64, offset int64, pairs ...types.Pair) (n int64, err error) {
	err = types.NewOperationNotImplementedError("write_page")
	return
}
func init() {
	services.RegisterFactory(Type, &Factory{})
}
