- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-10-11
- RFC PR: [beyondstorage/go-storage#837](https://github.com/beyondstorage/go-storage/issues/837)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-837: Support Feature Flag

- Updates:
  - [GSP-1: Unify storager behavior](./1-unify-storager-behavior.md): Abandon interface splitting
  - [GSP-44: Add CommitAppend in Appender](./44-commit-append.md): Move `CommitAppend` operations to `Storager`
  - [GSP-49: Add CreateDir Operation](./49-add-create-dir-operation.md): Move `CreateDir` operation to `Storager`
  - [GSP-86: Add Create Link Operation](./86-add-create-link-operation.md): Move `CreateLink` operation to `Storager`
  - [GSP-109: Redesign Features](./109-redesign-features.md): Reimplement `Features`
  - [GSP-729: Redesign HTTP Signer](./729-redesign-http-signer.md): Move HTTP signer related operations to `Storager`
  - [GSP-826: Add Multipart HTTP Signer Support](./826-add-multipart-http-signer-support.md): Move multipart HTTP signer related operations to `Storager`

## Background

In [GSP-87: Feature Gates](./87-feature-gates.md), we have the concept of `features` to resolve inconsistent behavior. And then in [GSP-109: Redesign Features](./109-redesign-features.md), we have a huge redesign about it. We have the following features for now:

- LoosePair
- VirtualDir
- VirtualLink
- VirtualObjectMetadata

It's obvious that they are only a subset of all features we support.

In fact, we can categorize all features into three kinds:

- Native support
- No native support but could be virtual
- Can't support

In the past, we only care about features that "No native support but could be virtual". Because other features can be handled easily: if we support, everything will work fine; If not, we just don't implement the interface.

However, [GSP-751: Write Empty File Behavior](./751-write-empty-file-behavior.md) rise a new question: How to handle the different behavior of the same interface?

Moreover, we do have different capabilities and limitations in storage services, for application or user, the capabilities to access the service is not accurately available by way of interface type assertion.

Maybe it's time for us to drop the idea that we introduced in [GSP-1: Unify storager behavior](./1-unify-storager-behavior.md).

## Proposal

I propose to support feature flag for operational rampups, and also provide interfaces for users or applications to acquire service capabilities.

### Reimplement Features

- Generate `ServiceFeatures` and `StorageFeatures` structs on go-storage side to represent the features supported by the storage container and storage service, respectively.
  - `CanXxx()` function is to check if the operation, or the feature that no native support but could be virtual is supported.
- Declare the supported features on service side.
  - Reserve the generated enable feature pairs `enable_xxx` and `EnableXxx()` functions for features that are not natively supported but can be supported virtually to enable the feature when init the storage.
  - `ServiceFeatures` and `StorageFeatures` pairs and structs will be deprecated on service side, as will the `WithStorageFeatures()`/`WithServiceFeatures()` functions.

```go
// StorageFeatures indicates features supported by storage.
type StorageFeatures struct {
	// native support or can't support 
	createAppend                   bool
	writeAppend                    bool
	commitAppend                   bool
	createBlock                    bool
	writeBlock                     bool
	combineBlock                   bool
	listBlock                      bool
	createDir                      bool
	fetch                          bool
	createLink                     bool
	move                           bool
	createMultipart                bool
	writeMultipart                 bool
	completeMultipart              bool
	listMultipart                  bool
	createPage                     bool
	writePage                      bool
	reach                          bool
	create                         bool
	delete                         bool
	metadata                       bool
	list                           bool
	read                           bool
	stat                           bool
	write                          bool
	querySignHTTPDelete            bool
	querySignHTTPRead              bool
	querySignHTTPWrite             bool
	querySignHTTPCreateMultipart   bool
	querySignHTTPCompleteMultipart bool
	querySignHTTPWriteMultipart    bool
	querySignHTTPlistMultipart     bool
	
	// no native support but could be virtual
	loosePair             bool
	virtualDir            bool
	virtualLink           bool
	virtualObjectMetadata bool
}

// CanCopy returns whether this storage support Copy operation or not.
func (f StorageFeatures) CanCopy() bool {
	return f.copy
}

...

// CanVirtualDir returns whether this storage support virtual_dir feature or not.
func (f *StorageFeatures) CanVirtualDir() {
	f.virtualDir = true
}

...
```

### Add operations back to Storager

- Add all operations in interfaces for storage service like `Copy`, `Move`, etc back to `Storager`.
- Support `Features()` which returns the supported features in `Storager` and `Servicer` interfaces, respectively.

```go
type Storager interface {
	Features() StorageFeatures
	
	String() string
	// Storage
	Create(path string, pairs ...Pair) (o *Object)
	Delete(path string, pairs ...Pair) (err error)
	DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error)
	...
	
	// Copy
	Copy(src string, dst string, pairs ...Pair) (err error)
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)
	
	...
	
	mustEmbedUnimplementedStorager()
}

// UnimplementedStorager must be embedded to have forward compatible implementations.
type UnimplementedStorager struct {
}

func (s UnimplementedStorager) mustEmbedUnimplementedStorager() {

}
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

...
```

- Stub implementations for all non-supported operations will be generated.
- The implementation of `Features()` will be generated according to the declared supported features.

Caller need to check the feature support flag before use them.

```go
if store.Features().CanCopy() {
	err := store.Copy(oldpath, newpath)
	if err != nil {
		return err
	}
}
```

## Rationale

### Feature Flag

Feature flags (also known as feature toggles or feature flipping) is a technique in development that allows you to enable or disable a feature without deploying any codes. The idea behind feature flags is to build conditional feature branches into code in order to make logic available only to certain groups of users at a time. If the flag is `on`, new code is executed, if the flag is `off`, the code is skipped.

Feature flag can be used in the following scenarios:

- Adding a new feature to an application.
- Enhancing an existing feature in an application.
- Hiding or disabling a feature.
- Extending an interface.

Feature flags solutions:

- [Unleash](https://www.getunleash.io/) is a feature management lets you turn new features on/off in production with no need for redeployment.
- [go-feature-flag](https://github.com/thomaspoignant/go-feature-flag) is a simple and complete feature flag solution, without any complex backend system to install.
- [Etsy's Feature flagging API](https://github.com/etsy/feature) used for operational rampups and A/B testing.

## Compatibility

- Interface type assertion for storage capability will be deprecated: the following interfaces will be deprecated, and the related operations will be added back to `Storager` interface. Caller need to check supported features before use them.
  - Appender
  - Blocker
  - Copier
  - Direr
  - Mover
  - Multiparter
  - Pager
  - Reacher
  - StorageHTTPSigner
  - MultipartHTTPSigner
- All API call that use `ServiceFeatures` and `StorageFeatures` could be affected. We could migrate as follows:
  - Mark `ServiceFeatures` and `StorageFeatures` related as deprecated.
  - Release a new version for go-storage and all services bump to this version with all references to `ServiceFeatures`, `StorageFeatures` etc updated.
  - Remove deprecated structs in the next major version.

## Implementation

- For go-storage
  - Deprecate the current generated `ServiceFeatures` and `StorageFeatures` pairs and structs for services, and re-generate them on go-storage side.
  - Remove interfaces other than `Servicer` and `Storager`, and add the related operations back to `Storager`.
- For storage services
  - Remove `implement` fields in `service.toml`.
  - Generate the implementation of `Features()` for `Storage` and `Service`.
- For go-integration-test
  - go-integration-test should not use type assertions, but run the test cases according to features.
  