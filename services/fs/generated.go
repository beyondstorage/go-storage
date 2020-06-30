// Code generated by go generate via internal/cmd/service; DO NOT EDIT.
package fs

import (
	"context"
	"io"

	"github.com/cns-io/go-storage/v2"
	"github.com/cns-io/go-storage/v2/pkg/credential"
	"github.com/cns-io/go-storage/v2/pkg/endpoint"
	"github.com/cns-io/go-storage/v2/pkg/httpclient"
	"github.com/cns-io/go-storage/v2/pkg/segment"
	"github.com/cns-io/go-storage/v2/services"
	"github.com/cns-io/go-storage/v2/types"
	"github.com/cns-io/go-storage/v2/types/info"
	ps "github.com/cns-io/go-storage/v2/types/pairs"
)

var _ credential.Provider
var _ endpoint.Provider
var _ segment.Segment
var _ storage.Storager
var _ services.ServiceError
var _ httpclient.Options

// Type is the type for fs
const Type = "fs"

// Service available pairs.
const ()

// Service available infos.
const ()

// pairStorageNewMap holds all available pairs
var pairStorageNewMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	ps.WorkDir: struct{}{},
	// Generated pairs
	ps.HTTPClientOptions: struct{}{},
}

// pairStorageNew is the parsed struct
type pairStorageNew struct {
	pairs []*types.Pair

	// Required pairs
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

// pairStorageCopyMap holds all available pairs
var pairStorageCopyMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageCopy is the parsed struct
type pairStorageCopy struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageCopy will parse *types.Pair slice into *pairStorageCopy
func parsePairStorageCopy(opts []*types.Pair) (*pairStorageCopy, error) {
	result := &pairStorageCopy{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageCopyMap[v.Key]; !ok {
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

// pairStorageMoveMap holds all available pairs
var pairStorageMoveMap = map[string]struct{}{
	// Required pairs
	// Optional pairs
	// Generated pairs
}

// pairStorageMove is the parsed struct
type pairStorageMove struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	// Generated pairs
}

// parsePairStorageMove will parse *types.Pair slice into *pairStorageMove
func parsePairStorageMove(opts []*types.Pair) (*pairStorageMove, error) {
	result := &pairStorageMove{
		pairs: opts,
	}

	values := make(map[string]interface{})
	for _, v := range opts {
		if _, ok := pairStorageMoveMap[v.Key]; !ok {
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
	ps.Offset: struct{}{},
	ps.Size:   struct{}{},
	// Generated pairs
	ps.ReadCallbackFunc: struct{}{},
}

// pairStorageRead is the parsed struct
type pairStorageRead struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	HasOffset bool
	Offset    int64
	HasSize   bool
	Size      int64
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
	v, ok = values[ps.Offset]
	if ok {
		result.HasOffset = true
		result.Offset = v.(int64)
	}
	v, ok = values[ps.Size]
	if ok {
		result.HasSize = true
		result.Size = v.(int64)
	}
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
	// Optional pairs
	ps.Size: struct{}{},
	// Generated pairs
	ps.ReadCallbackFunc: struct{}{},
}

// pairStorageWrite is the parsed struct
type pairStorageWrite struct {
	pairs []*types.Pair

	// Required pairs
	// Optional pairs
	HasSize bool
	Size    int64
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
	// Handle optional pairs
	v, ok = values[ps.Size]
	if ok {
		result.HasSize = true
		result.Size = v.(int64)
	}
	// Handle generated pairs
	v, ok = values[ps.ReadCallbackFunc]
	if ok {
		result.HasReadCallbackFunc = true
		result.ReadCallbackFunc = v.(func([]byte))
	}

	return result, nil
}

// Copy will copy an Object or multiple object in the service.
//
// This function will create a context by default.
func (s *Storage) Copy(src string, dst string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.CopyWithContext(ctx, src, dst, pairs...)
}

// CopyWithContext will copy an Object or multiple object in the service.
func (s *Storage) CopyWithContext(ctx context.Context, src string, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpCopy, err, src, dst)
	}()
	var opt *pairStorageCopy
	opt, err = parsePairStorageCopy(pairs)
	if err != nil {
		return
	}

	return s.copy(ctx, src, dst, opt)
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

// Move will move an object in the service.
//
// This function will create a context by default.
func (s *Storage) Move(src string, dst string, pairs ...*types.Pair) (err error) {
	ctx := context.Background()
	return s.MoveWithContext(ctx, src, dst, pairs...)
}

// MoveWithContext will move an object in the service.
func (s *Storage) MoveWithContext(ctx context.Context, src string, dst string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError(services.OpMove, err, src, dst)
	}()
	var opt *pairStorageMove
	opt, err = parsePairStorageMove(pairs)
	if err != nil {
		return
	}

	return s.move(ctx, src, dst, opt)
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
