package types

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Operation names in Appender.
const (
	// OpAppenderCommitAppend is the operation name for commit_append in appender.
	OpAppenderCommitAppend = "commit_append"
	// OpAppenderCreateAppend is the operation name for create_append in appender.
	OpAppenderCreateAppend = "create_append"
	// OpAppenderWriteAppend is the operation name for write_append in appender.
	OpAppenderWriteAppend = "write_append"
)

var _ = OpAppenderCommitAppend
var _ = OpAppenderCreateAppend
var _ = OpAppenderWriteAppend

// Appender is the interface for Append related operations.
type Appender interface {

	// CommitAppend will commit and finish an append process.
	CommitAppend(o *Object, pairs ...Pair) (err error)
	// CommitAppendWithContext will commit and finish an append process.
	CommitAppendWithContext(ctx context.Context, o *Object, pairs ...Pair) (err error)

	// CreateAppend will create an append object.
	//
	// ## Behavior
	//
	// - CreateAppend SHOULD create an appendable object with position 0 and size 0.
	// - CreateAppend SHOULD NOT return an error as the object exist.
	//   - Service SHOULD check and delete the object if exists.
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	// CreateAppendWithContext will create an append object.
	//
	// ## Behavior
	//
	// - CreateAppend SHOULD create an appendable object with position 0 and size 0.
	// - CreateAppend SHOULD NOT return an error as the object exist.
	//   - Service SHOULD check and delete the object if exists.
	CreateAppendWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WriteAppend will append content to an append object.
	WriteAppend(o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteAppendWithContext will append content to an append object.
	WriteAppendWithContext(ctx context.Context, o *Object, r io.Reader, size int64, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedAppender()
}

// UnimplementedAppender must be embedded to have forward compatible implementations.
type UnimplementedAppender struct{}

func (s UnimplementedAppender) mustEmbedUnimplementedAppender() {}

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

// Operation names in Blocker.
const (
	// OpBlockerCombineBlock is the operation name for combine_block in blocker.
	OpBlockerCombineBlock = "combine_block"
	// OpBlockerCreateBlock is the operation name for create_block in blocker.
	OpBlockerCreateBlock = "create_block"
	// OpBlockerListBlock is the operation name for list_block in blocker.
	OpBlockerListBlock = "list_block"
	// OpBlockerWriteBlock is the operation name for write_block in blocker.
	OpBlockerWriteBlock = "write_block"
)

var _ = OpBlockerCombineBlock
var _ = OpBlockerCreateBlock
var _ = OpBlockerListBlock
var _ = OpBlockerWriteBlock

// Blocker is the interface for Block related operations.
type Blocker interface {

	// CombineBlock will combine blocks into an object.
	CombineBlock(o *Object, bids []string, pairs ...Pair) (err error)
	// CombineBlockWithContext will combine blocks into an object.
	CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error)

	// CreateBlock will create a new block object.
	//
	// ## Behavior
	// - CreateBlock SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if exists.
	CreateBlock(path string, pairs ...Pair) (o *Object, err error)
	// CreateBlockWithContext will create a new block object.
	//
	// ## Behavior
	// - CreateBlock SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if exists.
	CreateBlockWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListBlock will list blocks belong to this object.
	ListBlock(o *Object, pairs ...Pair) (bi *BlockIterator, err error)
	// ListBlockWithContext will list blocks belong to this object.
	ListBlockWithContext(ctx context.Context, o *Object, pairs ...Pair) (bi *BlockIterator, err error)

	// WriteBlock will write content to a block.
	WriteBlock(o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)
	// WriteBlockWithContext will write content to a block.
	WriteBlockWithContext(ctx context.Context, o *Object, r io.Reader, size int64, bid string, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedBlocker()
}

// UnimplementedBlocker must be embedded to have forward compatible implementations.
type UnimplementedBlocker struct{}

func (s UnimplementedBlocker) mustEmbedUnimplementedBlocker() {}

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

// Operation names in Copier.
const (
	// OpCopierCopy is the operation name for copy in copier.
	OpCopierCopy = "copy"
)

var _ = OpCopierCopy

// Copier is the interface for Copy.
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
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object if exists.
	// - A successful copy opration should be complete, which means the dst object's content and metadata should be the same as src object.
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
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object if exists.
	// - A successful copy opration should be complete, which means the dst object's content and metadata should be the same as src object.
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)

	mustEmbedUnimplementedCopier()
}

// UnimplementedCopier must be embedded to have forward compatible implementations.
type UnimplementedCopier struct{}

func (s UnimplementedCopier) mustEmbedUnimplementedCopier() {}

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

// Operation names in Direr.
const (
	// OpDirerCreateDir is the operation name for create_dir in direr.
	OpDirerCreateDir = "create_dir"
)

var _ = OpDirerCreateDir

// Direr is the interface for Directory.
type Direr interface {

	// CreateDir will create a new dir object.
	CreateDir(path string, pairs ...Pair) (o *Object, err error)
	// CreateDirWithContext will create a new dir object.
	CreateDirWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	mustEmbedUnimplementedDirer()
}

// UnimplementedDirer must be embedded to have forward compatible implementations.
type UnimplementedDirer struct{}

func (s UnimplementedDirer) mustEmbedUnimplementedDirer() {}

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

// Operation names in Fetcher.
const (
	// OpFetcherFetch is the operation name for fetch in fetcher.
	OpFetcherFetch = "fetch"
)

var _ = OpFetcherFetch

// Fetcher is the interface for Fetch.
type Fetcher interface {

	// Fetch will fetch from a given url to path.
	//
	// ## Behavior
	//
	// - Fetch SHOULD NOT return an error as the object exists.
	// - A successful fetch operation should be complete, which means the object's content and metadata should be the same as requiring from the url.
	Fetch(path string, url string, pairs ...Pair) (err error)
	// FetchWithContext will fetch from a given url to path.
	//
	// ## Behavior
	//
	// - Fetch SHOULD NOT return an error as the object exists.
	// - A successful fetch operation should be complete, which means the object's content and metadata should be the same as requiring from the url.
	FetchWithContext(ctx context.Context, path string, url string, pairs ...Pair) (err error)

	mustEmbedUnimplementedFetcher()
}

// UnimplementedFetcher must be embedded to have forward compatible implementations.
type UnimplementedFetcher struct{}

func (s UnimplementedFetcher) mustEmbedUnimplementedFetcher() {}

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

// Operation names in HTTPSigner.
const (
	// OpHTTPSignerQuerySignHTTP is the operation name for query_sign_http in httpSigner.
	OpHTTPSignerQuerySignHTTP = "query_sign_http"
)

var _ = OpHTTPSignerQuerySignHTTP

// HTTPSigner is the interface for Signer.
type HTTPSigner interface {

	// QuerySignHTTP will return `*http.Request` with query string parameters containing signature in `URL` to represent the client's request.
	QuerySignHTTP(op string, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)
	// QuerySignHTTPWithContext will return `*http.Request` with query string parameters containing signature in `URL` to represent the client's request.
	QuerySignHTTPWithContext(ctx context.Context, op string, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error)

	mustEmbedUnimplementedHTTPSigner()
}

// UnimplementedHTTPSigner must be embedded to have forward compatible implementations.
type UnimplementedHTTPSigner struct{}

func (s UnimplementedHTTPSigner) mustEmbedUnimplementedHTTPSigner() {}

func (s UnimplementedHTTPSigner) String() string {
	return "UnimplementedHTTPSigner"
}

func (s UnimplementedHTTPSigner) QuerySignHTTP(op string, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http")
	return
}
func (s UnimplementedHTTPSigner) QuerySignHTTPWithContext(ctx context.Context, op string, path string, expire time.Duration, pairs ...Pair) (req *http.Request, err error) {
	err = NewOperationNotImplementedError("query_sign_http")
	return
}

// Operation names in Linker.
const (
	// OpLinkerCreateLink is the operation name for create_link in linker.
	OpLinkerCreateLink = "create_link"
)

var _ = OpLinkerCreateLink

// Linker is the interface for link
type Linker interface {

	// CreateLink Will create a link object.
	//
	// # Behavior
	//
	// - `path` and `target` COULD be relative or absolute path.
	// - If `target` not exists, CreateLink will still create a link object to path.
	// - If `path` exists:
	//   - If `path` is a symlink object, CreateLink will remove the symlink object and create a new link object to path.
	//   - If `path` is not a symlink object, CreateLink will return an ErrObjectModeInvalid error when the service does not support overwrite.
	// - A link object COULD be returned in `Stat` or `List`.
	// - CreateLink COULD implement virtual_link feature when service without native support.
	//   - Users SHOULD enable this feature by themselves.
	CreateLink(path string, target string, pairs ...Pair) (o *Object, err error)
	// CreateLinkWithContext Will create a link object.
	//
	// # Behavior
	//
	// - `path` and `target` COULD be relative or absolute path.
	// - If `target` not exists, CreateLink will still create a link object to path.
	// - If `path` exists:
	//   - If `path` is a symlink object, CreateLink will remove the symlink object and create a new link object to path.
	//   - If `path` is not a symlink object, CreateLink will return an ErrObjectModeInvalid error when the service does not support overwrite.
	// - A link object COULD be returned in `Stat` or `List`.
	// - CreateLink COULD implement virtual_link feature when service without native support.
	//   - Users SHOULD enable this feature by themselves.
	CreateLinkWithContext(ctx context.Context, path string, target string, pairs ...Pair) (o *Object, err error)

	mustEmbedUnimplementedLinker()
}

// UnimplementedLinker must be embedded to have forward compatible implementations.
type UnimplementedLinker struct{}

func (s UnimplementedLinker) mustEmbedUnimplementedLinker() {}

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

// Operation names in Mover.
const (
	// OpMoverMove is the operation name for move in mover.
	OpMoverMove = "move"
)

var _ = OpMoverMove

// Mover is the interface for Move.
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
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object if exists.
	// - A successful move operation SHOULD be complete, which means the dst object's content and metadata should be the same as src object.
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
	//   - Service that has native support for `overwrite` doesn't NEED to check the dst object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the dst object if exists.
	// - A successful move operation SHOULD be complete, which means the dst object's content and metadata should be the same as src object.
	MoveWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)

	mustEmbedUnimplementedMover()
}

// UnimplementedMover must be embedded to have forward compatible implementations.
type UnimplementedMover struct{}

func (s UnimplementedMover) mustEmbedUnimplementedMover() {}

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

// Operation names in Multiparter.
const (
	// OpMultiparterCompleteMultipart is the operation name for complete_multipart in multiparter.
	OpMultiparterCompleteMultipart = "complete_multipart"
	// OpMultiparterCreateMultipart is the operation name for create_multipart in multiparter.
	OpMultiparterCreateMultipart = "create_multipart"
	// OpMultiparterListMultipart is the operation name for list_multipart in multiparter.
	OpMultiparterListMultipart = "list_multipart"
	// OpMultiparterWriteMultipart is the operation name for write_multipart in multiparter.
	OpMultiparterWriteMultipart = "write_multipart"
)

var _ = OpMultiparterCompleteMultipart
var _ = OpMultiparterCreateMultipart
var _ = OpMultiparterListMultipart
var _ = OpMultiparterWriteMultipart

// Multiparter is the interface for Multipart related operations.
type Multiparter interface {

	// CompleteMultipart will complete a multipart upload and construct an Object.
	CompleteMultipart(o *Object, parts []*Part, pairs ...Pair) (err error)
	// CompleteMultipartWithContext will complete a multipart upload and construct an Object.
	CompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error)

	// CreateMultipart will create a new multipart.
	//
	// ## Behavior
	//
	// - CreateMultipart SHOULD NOT return an error as the object exists.
	CreateMultipart(path string, pairs ...Pair) (o *Object, err error)
	// CreateMultipartWithContext will create a new multipart.
	//
	// ## Behavior
	//
	// - CreateMultipart SHOULD NOT return an error as the object exists.
	CreateMultipartWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// ListMultipart will list parts belong to this multipart.
	ListMultipart(o *Object, pairs ...Pair) (pi *PartIterator, err error)
	// ListMultipartWithContext will list parts belong to this multipart.
	ListMultipartWithContext(ctx context.Context, o *Object, pairs ...Pair) (pi *PartIterator, err error)

	// WriteMultipart will write content to a multipart.
	WriteMultipart(o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error)
	// WriteMultipartWithContext will write content to a multipart.
	WriteMultipartWithContext(ctx context.Context, o *Object, r io.Reader, size int64, index int, pairs ...Pair) (n int64, part *Part, err error)

	mustEmbedUnimplementedMultiparter()
}

// UnimplementedMultiparter must be embedded to have forward compatible implementations.
type UnimplementedMultiparter struct{}

func (s UnimplementedMultiparter) mustEmbedUnimplementedMultiparter() {}

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

// Operation names in Pager.
const (
	// OpPagerCreatePage is the operation name for create_page in pager.
	OpPagerCreatePage = "create_page"
	// OpPagerWritePage is the operation name for write_page in pager.
	OpPagerWritePage = "write_page"
)

var _ = OpPagerCreatePage
var _ = OpPagerWritePage

// Pager is the interface for Page related operations which support random write.
type Pager interface {

	// CreatePage will create a new page object.
	//
	// ## Behavior
	//
	// - CreatePage SHOULD NOT return an error as the object exists.
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	// CreatePageWithContext will create a new page object.
	//
	// ## Behavior
	//
	// - CreatePage SHOULD NOT return an error as the object exists.
	CreatePageWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// WritePage will write content to specific offset.
	WritePage(o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)
	// WritePageWithContext will write content to specific offset.
	WritePageWithContext(ctx context.Context, o *Object, r io.Reader, size int64, offset int64, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedPager()
}

// UnimplementedPager must be embedded to have forward compatible implementations.
type UnimplementedPager struct{}

func (s UnimplementedPager) mustEmbedUnimplementedPager() {}

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

// Operation names in Reacher.
const (
	// OpReacherReach is the operation name for reach in reacher.
	OpReacherReach = "reach"
)

var _ = OpReacherReach

// Reacher is the interface for Reach.
type Reacher interface {

	// Reach will provide a way, which can reach the object.
	Reach(path string, pairs ...Pair) (url string, err error)
	// ReachWithContext will provide a way, which can reach the object.
	ReachWithContext(ctx context.Context, path string, pairs ...Pair) (url string, err error)

	mustEmbedUnimplementedReacher()
}

// UnimplementedReacher must be embedded to have forward compatible implementations.
type UnimplementedReacher struct{}

func (s UnimplementedReacher) mustEmbedUnimplementedReacher() {}

func (s UnimplementedReacher) String() string {
	return "UnimplementedReacher"
}

func (s UnimplementedReacher) Reach(path string, pairs ...Pair) (url string, err error) {
	err = NewOperationNotImplementedError("reach")
	return
}
func (s UnimplementedReacher) ReachWithContext(ctx context.Context, path string, pairs ...Pair) (url string, err error) {
	err = NewOperationNotImplementedError("reach")
	return
}

// Operation names in Servicer.
const (
	// OpServicerCreate is the operation name for create in servicer.
	OpServicerCreate = "create"
	// OpServicerDelete is the operation name for delete in servicer.
	OpServicerDelete = "delete"
	// OpServicerGet is the operation name for get in servicer.
	OpServicerGet = "get"
	// OpServicerList is the operation name for list in servicer.
	OpServicerList = "list"
)

var _ = OpServicerCreate
var _ = OpServicerDelete
var _ = OpServicerGet
var _ = OpServicerList

// Servicer can maintain multipart storage services.
type Servicer interface {
	String() string

	// Create will create a new storager instance.
	Create(name string, pairs ...Pair) (store Storager, err error)
	// CreateWithContext will create a new storager instance.
	CreateWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error)

	// Delete will delete a storager instance.
	Delete(name string, pairs ...Pair) (err error)
	// DeleteWithContext will delete a storager instance.
	DeleteWithContext(ctx context.Context, name string, pairs ...Pair) (err error)

	// Get will get a valid storager instance for service.
	Get(name string, pairs ...Pair) (store Storager, err error)
	// GetWithContext will get a valid storager instance for service.
	GetWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error)

	// List will list all storager instances under this service.
	List(pairs ...Pair) (sti *StoragerIterator, err error)
	// ListWithContext will list all storager instances under this service.
	ListWithContext(ctx context.Context, pairs ...Pair) (sti *StoragerIterator, err error)

	mustEmbedUnimplementedServicer()
}

// UnimplementedServicer must be embedded to have forward compatible implementations.
type UnimplementedServicer struct{}

func (s UnimplementedServicer) mustEmbedUnimplementedServicer() {}

func (s UnimplementedServicer) String() string {
	return "UnimplementedServicer"
}

func (s UnimplementedServicer) Create(name string, pairs ...Pair) (store Storager, err error) {
	err = NewOperationNotImplementedError("create")
	return
}
func (s UnimplementedServicer) CreateWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error) {
	err = NewOperationNotImplementedError("create")
	return
}

func (s UnimplementedServicer) Delete(name string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("delete")
	return
}
func (s UnimplementedServicer) DeleteWithContext(ctx context.Context, name string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("delete")
	return
}

func (s UnimplementedServicer) Get(name string, pairs ...Pair) (store Storager, err error) {
	err = NewOperationNotImplementedError("get")
	return
}
func (s UnimplementedServicer) GetWithContext(ctx context.Context, name string, pairs ...Pair) (store Storager, err error) {
	err = NewOperationNotImplementedError("get")
	return
}

func (s UnimplementedServicer) List(pairs ...Pair) (sti *StoragerIterator, err error) {
	err = NewOperationNotImplementedError("list")
	return
}
func (s UnimplementedServicer) ListWithContext(ctx context.Context, pairs ...Pair) (sti *StoragerIterator, err error) {
	err = NewOperationNotImplementedError("list")
	return
}

// Operation names in Storager.
const (
	// OpStoragerCreate is the operation name for create in storager.
	OpStoragerCreate = "create"
	// OpStoragerDelete is the operation name for delete in storager.
	OpStoragerDelete = "delete"
	// OpStoragerList is the operation name for list in storager.
	OpStoragerList = "list"
	// OpStoragerMetadata is the operation name for metadata in storager.
	OpStoragerMetadata = "metadata"
	// OpStoragerRead is the operation name for read in storager.
	OpStoragerRead = "read"
	// OpStoragerStat is the operation name for stat in storager.
	OpStoragerStat = "stat"
	// OpStoragerWrite is the operation name for write in storager.
	OpStoragerWrite = "write"
)

var _ = OpStoragerCreate
var _ = OpStoragerDelete
var _ = OpStoragerList
var _ = OpStoragerMetadata
var _ = OpStoragerRead
var _ = OpStoragerStat
var _ = OpStoragerWrite

// Storager is the interface for storage service.
type Storager interface {
	String() string

	// Create will create a new object without any api call.
	//
	// ## Behavior
	//
	// - Create SHOULD NOT send any API call.
	// - Create SHOULD accept ObjectMode pair as object mode.
	Create(path string, pairs ...Pair) (o *Object)

	// Delete will delete an object from service.
	//
	// ## Behavior
	//
	// - Delete only delete one and only one object.
	//   - Service DON'T NEED to support remove all.
	//   - User NEED to implement remove_all by themself.
	// - Delete is idempotent.
	//   - Successful delete always return nil error.
	//   - Delete SHOULD never return `ObjectNotExist`
	//   - Delete DON'T NEED to check the object exist or not.
	Delete(path string, pairs ...Pair) (err error)
	// DeleteWithContext will delete an object from service.
	//
	// ## Behavior
	//
	// - Delete only delete one and only one object.
	//   - Service DON'T NEED to support remove all.
	//   - User NEED to implement remove_all by themself.
	// - Delete is idempotent.
	//   - Successful delete always return nil error.
	//   - Delete SHOULD never return `ObjectNotExist`
	//   - Delete DON'T NEED to check the object exist or not.
	DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error)

	// List will return list a specific path.
	//
	// ## Behavior
	//
	// - Service SHOULD support default `ListMode`.
	// - Service SHOULD implement `ListModeDir` without the check for `VirtualDir`.
	// - Service DON'T NEED to `Stat` while in `List`.
	List(path string, pairs ...Pair) (oi *ObjectIterator, err error)
	// ListWithContext will return list a specific path.
	//
	// ## Behavior
	//
	// - Service SHOULD support default `ListMode`.
	// - Service SHOULD implement `ListModeDir` without the check for `VirtualDir`.
	// - Service DON'T NEED to `Stat` while in `List`.
	ListWithContext(ctx context.Context, path string, pairs ...Pair) (oi *ObjectIterator, err error)

	// Metadata will return current storager metadata.
	Metadata(pairs ...Pair) (meta *StorageMeta)

	// Read will read the file's data.
	Read(path string, w io.Writer, pairs ...Pair) (n int64, err error)
	// ReadWithContext will read the file's data.
	ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...Pair) (n int64, err error)

	// Stat will stat a path to get info of an object.
	//
	// ## Behavior
	//
	// - Stat SHOULD accept ObjectMode pair as hints.
	//   - Service COULD have different implementations for different object mode.
	//   - Service SHOULD check if returning ObjectMode is match
	Stat(path string, pairs ...Pair) (o *Object, err error)
	// StatWithContext will stat a path to get info of an object.
	//
	// ## Behavior
	//
	// - Stat SHOULD accept ObjectMode pair as hints.
	//   - Service COULD have different implementations for different object mode.
	//   - Service SHOULD check if returning ObjectMode is match
	StatWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error)

	// Write will write data into a file.
	//
	// ## Behavior
	//
	// - Write SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if exists.
	// - A successful write operation SHOULD be complete, which means the object's content and metadata should be the same as specified in write request.
	Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteWithContext will write data into a file.
	//
	// ## Behavior
	//
	// - Write SHOULD NOT return an error as the object exist.
	//   - Service that has native support for `overwrite` doesn't NEED to check the object exists or not.
	//   - Service that doesn't have native support for `overwrite` SHOULD check and delete the object if exists.
	// - A successful write operation SHOULD be complete, which means the object's content and metadata should be the same as specified in write request.
	WriteWithContext(ctx context.Context, path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)

	mustEmbedUnimplementedStorager()
}

// UnimplementedStorager must be embedded to have forward compatible implementations.
type UnimplementedStorager struct{}

func (s UnimplementedStorager) mustEmbedUnimplementedStorager() {}

func (s UnimplementedStorager) String() string {
	return "UnimplementedStorager"
}

func (s UnimplementedStorager) Create(path string, pairs ...Pair) (o *Object) {
	return
}

func (s UnimplementedStorager) Delete(path string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("delete")
	return
}
func (s UnimplementedStorager) DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error) {
	err = NewOperationNotImplementedError("delete")
	return
}

func (s UnimplementedStorager) List(path string, pairs ...Pair) (oi *ObjectIterator, err error) {
	err = NewOperationNotImplementedError("list")
	return
}
func (s UnimplementedStorager) ListWithContext(ctx context.Context, path string, pairs ...Pair) (oi *ObjectIterator, err error) {
	err = NewOperationNotImplementedError("list")
	return
}

func (s UnimplementedStorager) Metadata(pairs ...Pair) (meta *StorageMeta) {
	return
}

func (s UnimplementedStorager) Read(path string, w io.Writer, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("read")
	return
}
func (s UnimplementedStorager) ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("read")
	return
}

func (s UnimplementedStorager) Stat(path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("stat")
	return
}
func (s UnimplementedStorager) StatWithContext(ctx context.Context, path string, pairs ...Pair) (o *Object, err error) {
	err = NewOperationNotImplementedError("stat")
	return
}

func (s UnimplementedStorager) Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write")
	return
}
func (s UnimplementedStorager) WriteWithContext(ctx context.Context, path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {
	err = NewOperationNotImplementedError("write")
	return
}
