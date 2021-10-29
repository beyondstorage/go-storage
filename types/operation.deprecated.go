package types

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Appender is the interface for Append related operations.
//
// Deprecated: Moved to Storager.
type Appender interface {
	// CommitAppend will commit and finish an append process.
	//
	// Deprecated: Moved to Storager.
	CommitAppend(o *Object, pairs ...Pair) (err error)
	// CommitAppendWithContext will commit and finish an append process.
	//
	// Deprecated: Moved to Storager.
	CommitAppendWithContext(ctx context.Context, o *Object, pairs ...Pair) (err error)

	// CreateAppend will create an append object.
	//
	// ## Behavior
	//
	// - CreateAppend SHOULD create an appendable object with position 0 and size 0.
	// - CreateAppend SHOULD NOT return an error as the object exist.
	//   - Service SHOULD check and delete the object if exists.
	//
	// Deprecated: Moved to Storager.
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	// CreateAppendWithContext will create an append object.
	//
	// ## Behavior
	//
	// - CreateAppend SHOULD create an appendable object with position 0 and size 0.
	// - CreateAppend SHOULD NOT return an error as the object exist.
	//   - Service SHOULD check and delete the object if exists.
	//
	// Deprecated: Moved to Storager.
	CreateAppendWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WriteAppend will append content to an append object.
	//
	// Deprecated: Moved to Storager.
	WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteAppendWithContext will append content to an append object.
	//
	// Deprecated: Moved to Storager.
	WriteAppendWithContext(ctx context.Context, o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedAppender()
}

// UnimplementedAppender must be embedded to have forward compatible implementations.
//
// Appender is the interface for Append related operations.
//
// Deprecated: Moved to Storager.
type UnimplementedAppender struct {
}

func (s UnimplementedAppender) mustEmbedUnimplementedAppender() {

}
func (s UnimplementedAppender) String() string {
	return "UnimplementedAppender"
}
func (s UnimplementedAppender) CommitAppend(o *Object, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("commit_append")
	return
}
func (s UnimplementedAppender) CommitAppendWithContext(ctx context.Context, o *Object, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("commit_append")
	return
}
func (s UnimplementedAppender) CreateAppend(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_append")
	return
}
func (s UnimplementedAppender) CreateAppendWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_append")
	return
}
func (s UnimplementedAppender) WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_append")
	return
}
func (s UnimplementedAppender) WriteAppendWithContext(ctx context.Context, o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_append")
	return
}

// Blocker is the interface for Block related operations.
//
// Deprecated: Moved to Storager.
type Blocker interface {
	// CombineBlock will combine blocks into an object.
	//
	// Deprecated: Moved to Storager.
	CombineBlock(o *Object, bids []string, pairs ...Pair) (err error)
	// CombineBlockWithContext will combine blocks into an object.
	//
	// Deprecated: Moved to Storager.
	CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error)

	// CreateBlock will create a new block object.
	//
	// ## Behavior
	// - CreateBlock SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if
	// exists.
	//
	// Deprecated: Moved to Storager.
	CreateBlock(path string, pairs ...Pair) (o *Object, err error)
	// CreateBlockWithContext will create a new block object.
	//
	// ## Behavior
	// - CreateBlock SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if
	// exists.
	//
	// Deprecated: Moved to Storager.
	CreateBlockWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListBlock will list blocks belong to this object.
	//
	// Deprecated: Moved to Storager.
	ListBlock(o *Object, pairs ...Pair) (bi *BlockIterator, err error)
	// ListBlockWithContext will list blocks belong to this object.
	//
	// Deprecated: Moved to Storager.
	ListBlockWithContext(ctx context.Context, o *Object, pairs ...Pair) (bi *BlockIterator, err error)

	// WriteBlock will write content to a block.
	//
	// Deprecated: Moved to Storager.
	WriteBlock(o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)
	// WriteBlockWithContext will write content to a block.
	//
	// Deprecated: Moved to Storager.
	WriteBlockWithContext(ctx context.Context, o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedBlocker()
}

// UnimplementedBlocker must be embedded to have forward compatible implementations.
//
// Blocker is the interface for Block related operations.
//
// Deprecated: Moved to Storager.
type UnimplementedBlocker struct {
}

func (s UnimplementedBlocker) mustEmbedUnimplementedBlocker() {

}
func (s UnimplementedBlocker) String() string {
	return "UnimplementedBlocker"
}
func (s UnimplementedBlocker) CombineBlock(o *Object, bids []string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("combine_block")
	return
}
func (s UnimplementedBlocker) CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("combine_block")
	return
}
func (s UnimplementedBlocker) CreateBlock(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_block")
	return
}
func (s UnimplementedBlocker) CreateBlockWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_block")
	return
}
func (s UnimplementedBlocker) ListBlock(o *Object, pairs ...Pair) (bi *BlockIterator, err error) {
	err = NewOperationNotImplementedError("list_block")
	return
}
func (s UnimplementedBlocker) ListBlockWithContext(ctx context.Context, o *Object, pairs ...Pair) (bi *BlockIterator, err error) {
	err = NewOperationNotImplementedError("list_block")
	return
}
func (s UnimplementedBlocker) WriteBlock(o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_block")
	return
}
func (s UnimplementedBlocker) WriteBlockWithContext(ctx context.Context, o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_block")
	return
}

// Copier is the interface for Copy.
//
// Deprecated: Moved to Storager.
type Copier interface {
	// Copy will copy an Object or multiple object in the service.
	//
	// ## Behavior
	//
	// - Copy only copy one and only one object.
	//   - Service DON'T NEED to support copy a non-empty directory or copy files recursively.
	//   - User NEED to implement copy a non-empty directory and copy recursively by themself.
	//   - Copy a file to a directory SHOULD return `ErrObjectModeInvalid`.
	// - Copy SHOULD NOT return an error as dst object exists.
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or
	// not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object
	// if exists.
	// - A successful copy opration should be complete, which means the dst object's content and metadata
	// should be the same as src object.
	//
	// Deprecated: Moved to Storager.
	Copy(src string, dst string, pairs ...Pair) (err error)
	// CopyWithContext will copy an Object or multiple object in the service.
	//
	// ## Behavior
	//
	// - Copy only copy one and only one object.
	//   - Service DON'T NEED to support copy a non-empty directory or copy files recursively.
	//   - User NEED to implement copy a non-empty directory and copy recursively by themself.
	//   - Copy a file to a directory SHOULD return `ErrObjectModeInvalid`.
	// - Copy SHOULD NOT return an error as dst object exists.
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or
	// not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object
	// if exists.
	// - A successful copy opration should be complete, which means the dst object's content and metadata
	// should be the same as src object.
	//
	// Deprecated: Moved to Storager.
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)

	mustEmbedUnimplementedCopier()
}

// UnimplementedCopier must be embedded to have forward compatible implementations.
//
// Copier is the interface for Copy.
//
// Deprecated: Moved to Storager.
type UnimplementedCopier struct {
}

func (s UnimplementedCopier) mustEmbedUnimplementedCopier() {

}
func (s UnimplementedCopier) String() string {
	return "UnimplementedCopier"
}
func (s UnimplementedCopier) Copy(src string, dst string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("copy")
	return
}
func (s UnimplementedCopier) CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("copy")
	return
}

// Direr is the interface for Directory.
//
// Deprecated: Moved to Storager.
type Direr interface {
	// CreateDir will create a new dir object.
	//
	// Deprecated: Moved to Storager.
	CreateDir(path string, pairs ...Pair) (o *Object, err error)
	// CreateDirWithContext will create a new dir object.
	//
	// Deprecated: Moved to Storager.
	CreateDirWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	mustEmbedUnimplementedDirer()
}

// UnimplementedDirer must be embedded to have forward compatible implementations.
//
// Direr is the interface for Directory.
//
// Deprecated: Moved to Storager.
type UnimplementedDirer struct {
}

func (s UnimplementedDirer) mustEmbedUnimplementedDirer() {

}
func (s UnimplementedDirer) String() string {
	return "UnimplementedDirer"
}
func (s UnimplementedDirer) CreateDir(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_dir")
	return
}
func (s UnimplementedDirer) CreateDirWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_dir")
	return
}

// Fetcher is the interface for Fetch.
//
// Deprecated: Moved to Storager.
type Fetcher interface {
	// Fetch will fetch from a given url to path.
	//
	// ## Behavior
	//
	// - Fetch SHOULD NOT return an error as the object exists.
	// - A successful fetch operation should be complete, which means the object's content and metadata
	// should be the same as requiring from the url.
	//
	// Deprecated: Moved to Storager.
	Fetch(path string, url string, pairs ...Pair) (err error)
	// FetchWithContext will fetch from a given url to path.
	//
	// ## Behavior
	//
	// - Fetch SHOULD NOT return an error as the object exists.
	// - A successful fetch operation should be complete, which means the object's content and metadata
	// should be the same as requiring from the url.
	//
	// Deprecated: Moved to Storager.
	FetchWithContext(ctx context.Context, path string, url string, pairs ...Pair) (err error)

	mustEmbedUnimplementedFetcher()
}

// UnimplementedFetcher must be embedded to have forward compatible implementations.
//
// Fetcher is the interface for Fetch.
//
// Deprecated: Moved to Storager.
type UnimplementedFetcher struct {
}

func (s UnimplementedFetcher) mustEmbedUnimplementedFetcher() {

}
func (s UnimplementedFetcher) String() string {
	return "UnimplementedFetcher"
}
func (s UnimplementedFetcher) Fetch(path string, url string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("fetch")
	return
}
func (s UnimplementedFetcher) FetchWithContext(ctx context.Context, path string, url string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("fetch")
	return
}

// Linker is the interface for link.
//
// Deprecated: Moved to Storager.
type Linker interface {
	// CreateLink Will create a link object.
	//
	// # Behavior
	//
	// - `path` and `target` COULD be relative or absolute path.
	// - If `target` not exists, CreateLink will still create a link object to path.
	// - If `path` exists:
	//   - If `path` is a symlink object, CreateLink will remove the symlink object and create a new link object
	// to path.
	//   - If `path` is not a symlink object, CreateLink will return an ErrObjectModeInvalid error when
	// the service does not support overwrite.
	// - A link object COULD be returned in `Stat` or `List`.
	// - CreateLink COULD implement virtual_link feature when service without native support.
	//   - Users SHOULD enable this feature by themselves.
	//
	// Deprecated: Moved to Storager.
	CreateLink(path string, target string, pairs ...Pair) (o *Object, err error)
	// CreateLinkWithContext Will create a link object.
	//
	// # Behavior
	//
	// - `path` and `target` COULD be relative or absolute path.
	// - If `target` not exists, CreateLink will still create a link object to path.
	// - If `path` exists:
	//   - If `path` is a symlink object, CreateLink will remove the symlink object and create a new link object
	// to path.
	//   - If `path` is not a symlink object, CreateLink will return an ErrObjectModeInvalid error when
	// the service does not support overwrite.
	// - A link object COULD be returned in `Stat` or `List`.
	// - CreateLink COULD implement virtual_link feature when service without native support.
	//   - Users SHOULD enable this feature by themselves.
	//
	// Deprecated: Moved to Storager.
	CreateLinkWithContext(ctx context.Context, path string, target string, pairs ...Pair) (o *Object, err error)

	mustEmbedUnimplementedLinker()
}

// UnimplementedLinker must be embedded to have forward compatible implementations.
//
// Linker is the interface for link.
//
// Deprecated: Moved to Storager.
type UnimplementedLinker struct {
}

func (s UnimplementedLinker) mustEmbedUnimplementedLinker() {

}
func (s UnimplementedLinker) String() string {
	return "UnimplementedLinker"
}
func (s UnimplementedLinker) CreateLink(path string, target string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_link")
	return
}
func (s UnimplementedLinker) CreateLinkWithContext(ctx context.Context, path string, target string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_link")
	return
}

// Mover is the interface for Move.
//
// Deprecated: Moved to Storager.
type Mover interface {
	// Move will move an object in the service.
	//
	// ## Behavior
	//
	// - Move only move one and only one object.
	//   - Service DON'T NEED to support move a non-empty directory.
	//   - User NEED to implement move a non-empty directory by themself.
	//   - Move a file to a directory SHOULD return `ErrObjectModeInvalid`.
	// - Move SHOULD NOT return an error as dst object exists.
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or
	// not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object
	// if exists.
	// - A successful move operation SHOULD be complete, which means the dst object's content and metadata
	// should be the same as src object.
	//
	// Deprecated: Moved to Storager.
	Move(src string, dst string, pairs ...Pair) (err error)
	// MoveWithContext will move an object in the service.
	//
	// ## Behavior
	//
	// - Move only move one and only one object.
	//   - Service DON'T NEED to support move a non-empty directory.
	//   - User NEED to implement move a non-empty directory by themself.
	//   - Move a file to a directory SHOULD return `ErrObjectModeInvalid`.
	// - Move SHOULD NOT return an error as dst object exists.
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or
	// not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object
	// if exists.
	// - A successful move operation SHOULD be complete, which means the dst object's content and metadata
	// should be the same as src object.
	//
	// Deprecated: Moved to Storager.
	MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)

	mustEmbedUnimplementedMover()
}

// UnimplementedMover must be embedded to have forward compatible implementations.
//
// Mover is the interface for Move.
//
// Deprecated: Moved to Storager.
type UnimplementedMover struct {
}

func (s UnimplementedMover) mustEmbedUnimplementedMover() {

}
func (s UnimplementedMover) String() string {
	return "UnimplementedMover"
}
func (s UnimplementedMover) Move(src string, dst string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("move")
	return
}
func (s UnimplementedMover) MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("move")
	return
}

// MultipartHTTPSigner is the interface for Multiparter related operations which support authentication.
//
// Deprecated: Moved to Storager.
type MultipartHTTPSigner interface {
	// QuerySignHTTPCompleteMultipart will complete a multipart upload and construct an Object by
	// using query parameters to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPCompleteMultipart(o *Object, parts []*Part, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPCompleteMultipartWithContext will complete a multipart upload and construct
	// an Object by using query parameters to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPCompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	// QuerySignHTTPCreateMultipart will create a new multipart by using query parameters to authenticate
	// requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPCreateMultipart(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPCreateMultipartWithContext will create a new multipart by using query parameters
	// to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPCreateMultipartWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	// QuerySignHTTPListMultipart will list parts belong to this multipart by using query parameters
	// to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPListMultipart(o *Object, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPListMultipartWithContext will list parts belong to this multipart by using query
	// parameters to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPListMultipartWithContext(ctx context.Context, o *Object, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	// QuerySignHTTPWriteMultipart will write content to a multipart by using query parameters to authenticate
	// requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPWriteMultipart(o *Object, size int64, index int, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPWriteMultipartWithContext will write content to a multipart by using query parameters
	// to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPWriteMultipartWithContext(ctx context.Context, o *Object, size int64, index int, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	mustEmbedUnimplementedMultipartHTTPSigner()
}

// UnimplementedMultipartHTTPSigner must be embedded to have forward compatible implementations.
//
// MultipartHTTPSigner is the interface for Multiparter related operations which support authentication.
//
// Deprecated: Moved to Storager.
type UnimplementedMultipartHTTPSigner struct {
}

func (s UnimplementedMultipartHTTPSigner) mustEmbedUnimplementedMultipartHTTPSigner() {

}
func (s UnimplementedMultipartHTTPSigner) String() string {
	return "UnimplementedMultipartHTTPSigner"
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPCompleteMultipart(o *Object, parts []*Part, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_complete_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPCompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_complete_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPCreateMultipart(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_create_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPCreateMultipartWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_create_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPListMultipart(o *Object, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_list_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPListMultipartWithContext(ctx context.Context, o *Object, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_list_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPWriteMultipart(o *Object, size int64, index int, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_write_multipart")
	return
}
func (s UnimplementedMultipartHTTPSigner) QuerySignHTTPWriteMultipartWithContext(ctx context.Context, o *Object, size int64, index int, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_write_multipart")
	return
}

// Multiparter is the interface for Multipart related operations.
//
// Deprecated: Moved to Storager.
type Multiparter interface {
	// CompleteMultipart will complete a multipart upload and construct an Object.
	//
	// Deprecated: Moved to Storager.
	CompleteMultipart(o *Object, parts []*Part, pairs ...Pair) (err error)
	// CompleteMultipartWithContext will complete a multipart upload and construct an Object.
	//
	// Deprecated: Moved to Storager.
	CompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error)

	// CreateMultipart will create a new multipart.
	//
	// ## Behavior
	//
	// - CreateMultipart SHOULD NOT return an error as the object exists.
	//
	// Deprecated: Moved to Storager.
	CreateMultipart(path string, pairs ...Pair) (o *Object, err error)
	// CreateMultipartWithContext will create a new multipart.
	//
	// ## Behavior
	//
	// - CreateMultipart SHOULD NOT return an error as the object exists.
	//
	// Deprecated: Moved to Storager.
	CreateMultipartWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListMultipart will list parts belong to this multipart.
	//
	// Deprecated: Moved to Storager.
	ListMultipart(o *Object, pairs ...Pair) (pi *PartIterator, err error)
	// ListMultipartWithContext will list parts belong to this multipart.
	//
	// Deprecated: Moved to Storager.
	ListMultipartWithContext(ctx context.Context, o *Object, pairs ...Pair) (pi *PartIterator, err error)

	// WriteMultipart will write content to a multipart.
	//
	// Deprecated: Moved to Storager.
	WriteMultipart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error)
	// WriteMultipartWithContext will write content to a multipart.
	//
	// Deprecated: Moved to Storager.
	WriteMultipartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error)

	mustEmbedUnimplementedMultiparter()
}

// UnimplementedMultiparter must be embedded to have forward compatible implementations.
//
// Multiparter is the interface for Multipart related operations.
//
// Deprecated: Moved to Storager.
type UnimplementedMultiparter struct {
}

func (s UnimplementedMultiparter) mustEmbedUnimplementedMultiparter() {

}
func (s UnimplementedMultiparter) String() string {
	return "UnimplementedMultiparter"
}
func (s UnimplementedMultiparter) CompleteMultipart(o *Object, parts []*Part, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("complete_multipart")
	return
}
func (s UnimplementedMultiparter) CompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("complete_multipart")
	return
}
func (s UnimplementedMultiparter) CreateMultipart(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_multipart")
	return
}
func (s UnimplementedMultiparter) CreateMultipartWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_multipart")
	return
}
func (s UnimplementedMultiparter) ListMultipart(o *Object, pairs ...Pair) (pi *PartIterator, err error) {
	err = NewOperationNotImplementedError("list_multipart")
	return
}
func (s UnimplementedMultiparter) ListMultipartWithContext(ctx context.Context, o *Object, pairs ...Pair) (pi *PartIterator, err error) {
	err = NewOperationNotImplementedError("list_multipart")
	return
}
func (s UnimplementedMultiparter) WriteMultipart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error) {
	err = NewOperationNotImplementedError("write_multipart")
	return
}
func (s UnimplementedMultiparter) WriteMultipartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error) {
	err = NewOperationNotImplementedError("write_multipart")
	return
}

// Pager is the interface for Page related operations which support random write.
//
// Deprecated: Moved to Storager.
type Pager interface {
	// CreatePage will create a new page object.
	//
	// ## Behavior
	//
	// - CreatePage SHOULD NOT return an error as the object exists.
	//
	// Deprecated: Moved to Storager.
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	// CreatePageWithContext will create a new page object.
	//
	// ## Behavior
	//
	// - CreatePage SHOULD NOT return an error as the object exists.
	//
	// Deprecated: Moved to Storager.
	CreatePageWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WritePage will write content to specific offset.
	//
	// Deprecated: Moved to Storager.
	WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
	// WritePageWithContext will write content to specific offset.
	//
	// Deprecated: Moved to Storager.
	WritePageWithContext(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedPager()
}

// UnimplementedPager must be embedded to have forward compatible implementations.
//
// Pager is the interface for Page related operations which support random write.
//
// Deprecated: Moved to Storager.
type UnimplementedPager struct {
}

func (s UnimplementedPager) mustEmbedUnimplementedPager() {

}
func (s UnimplementedPager) String() string {
	return "UnimplementedPager"
}
func (s UnimplementedPager) CreatePage(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_page")
	return
}
func (s UnimplementedPager) CreatePageWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("create_page")
	return
}
func (s UnimplementedPager) WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_page")
	return
}
func (s UnimplementedPager) WritePageWithContext(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write_page")
	return
}

// StorageHTTPSigner is the interface for Storager related operations which support authentication.
//
// Deprecated: Moved to Storager.
type StorageHTTPSigner interface {
	// QuerySignHTTPDelete will delete an object from service by using query parameters to authenticate
	// requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPDelete(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPDeleteWithContext will delete an object from service by using query parameters
	// to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPDeleteWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	// QuerySignHTTPRead will read data from the file by using query parameters to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPRead(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPReadWithContext will read data from the file by using query parameters to authenticate
	// requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPReadWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	// QuerySignHTTPWrite will write data into a file by using query parameters to authenticate requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPWrite(path string, size int64, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPWriteWithContext will write data into a file by using query parameters to authenticate
	// requests.
	//
	// Deprecated: Moved to Storager.
	QuerySignHTTPWriteWithContext(ctx context.Context, path string, size int64, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	mustEmbedUnimplementedStorageHTTPSigner()
}

// UnimplementedStorageHTTPSigner must be embedded to have forward compatible implementations.
//
// StorageHTTPSigner is the interface for Storager related operations which support authentication.
//
// Deprecated: Moved to Storager.
type UnimplementedStorageHTTPSigner struct {
}

func (s UnimplementedStorageHTTPSigner) mustEmbedUnimplementedStorageHTTPSigner() {

}
func (s UnimplementedStorageHTTPSigner) String() string {
	return "UnimplementedStorageHTTPSigner"
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPDelete(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_delete")
	return
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPDeleteWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_delete")
	return
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPRead(path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_read")
	return
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPReadWithContext(ctx context.Context, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_read")
	return
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPWrite(path string, size int64, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_write")
	return
}
func (s UnimplementedStorageHTTPSigner) QuerySignHTTPWriteWithContext(ctx context.Context, path string, size int64, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http_write")
	return
}
