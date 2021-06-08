package types

import (
	"context"
	"io"
)

// Appender is the interface for Append related operations.
type Appender interface {

	// CommitAppend will commit and finish an append process.
	CommitAppend(o *Object, pairs ...Pair) (err error)
	// CommitAppendWithContext will commit and finish an append process.
	CommitAppendWithContext(ctx context.Context, o *Object, pairs ...Pair) (err error)

	// CreateAppend will create an append object.
	CreateAppend(path string, pairs ...Pair) (o *Object, err error)
	// CreateAppendWithContext will create an append object.
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

// Blocker is the interface for Block related operations.
type Blocker interface {

	// CombineBlock will combine blocks into an object.
	CombineBlock(o *Object, bids []string, pairs ...Pair) (err error)
	// CombineBlockWithContext will combine blocks into an object.
	CombineBlockWithContext(ctx context.Context, o *Object, bids []string, pairs ...Pair) (err error)

	// CreateBlock will create a new block object.
	CreateBlock(path string, pairs ...Pair) (o *Object, err error)
	// CreateBlockWithContext will create a new block object.
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

// Copier is the interface for Copy.
type Copier interface {

	// Copy will copy an Object or multiple object in the service.
	Copy(src string, dst string, pairs ...Pair) (err error)
	// CopyWithContext will copy an Object or multiple object in the service.
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

// Fetcher is the interface for Fetch.
type Fetcher interface {

	// Fetch will fetch from a given url to path.
	Fetch(path string, url string, pairs ...Pair) (err error)
	// FetchWithContext will fetch from a given url to path.
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

// Mover is the interface for Move.
type Mover interface {

	// Move will move an object in the service.
	Move(src string, dst string, pairs ...Pair) (err error)
	// MoveWithContext will move an object in the service.
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

// Multiparter is the interface for Multipart related operations.
type Multiparter interface {

	// CompleteMultipart will complete a multipart upload and construct an Object.
	CompleteMultipart(o *Object, parts []*Part, pairs ...Pair) (err error)
	// CompleteMultipartWithContext will complete a multipart upload and construct an Object.
	CompleteMultipartWithContext(ctx context.Context, o *Object, parts []*Part, pairs ...Pair) (err error)

	// CreateMultipart will create a new multipart.
	CreateMultipart(path string, pairs ...Pair) (o *Object, err error)
	// CreateMultipartWithContext will create a new multipart.
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

// Pager is the interface for Page related operations which support random write.
type Pager interface {

	// CreatePage will create a new page object.
	CreatePage(path string, pairs ...Pair) (o *Object, err error)
	// CreatePageWithContext will create a new page object.
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
	List(path string, pairs ...Pair) (oi *ObjectIterator, err error)
	// ListWithContext will return list a specific path.
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
	Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
	// WriteWithContext will write data into a file.
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
