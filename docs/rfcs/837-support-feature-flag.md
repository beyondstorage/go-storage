- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-10-11
- RFC PR: [beyondstorage/go-storage#837](https://github.com/beyondstorage/go-storage/issues/837)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-837: Support Feature Flag

- Updates:
  - [GSP-1: Unify storager behavior](./1-unify-storager-behavior.md): Abandon interface splitting

## Background

In [GSP-87: Feature Gates](./87-feature-gates.md), we have the concept of `features` first. And then in [GSP-109: Redesign Features](./109-redesign-features.md), we have a huge redesign about it. We have the following features for now:

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

Maybe it's time for us to drop the idea that we introduced in [GSP-1: Unify storager behavior](./1-unify-storager-behavior.md).

## Proposal

I propose to support feature flag for operational rampups, and also provide interfaces for users or applications to acquire service capabilities.

### Support Feature

Generate an uint64 `Feature` to represent the features supported by the storage service:

```go
// Feature is a uint64 which represents the features storage service have.
type Feature uint64

// All features that storage used.
const (
	FeatureCreate = 1 << iota
	FeatureDelete
	...
	
	FeatureCopy
	
	FeatureLoosePair
	FeatureVirtualDir
	FeatureVirtualLink
	FeatureVirtualObjectMetadata
)

// CanCopy returns whether this storage support Copy operation or not.
func (f Feature) CanCopy() bool {
	return c&FeatureCopy != 0
}

...
```

### Add operations back to Storager

Add all storage related operations like `Copy`, `Move`, etc back to `Storager`, and support returning an uint64 to represent `Feature` in `Storager`:

```go
type Storager interface {
	Feature() Feature
	
	// Storage
	Create(path string, pairs ...Pair) (o *Object)
	Delete(path string, pairs ...Pair) (err error)
	DeleteWithContext(ctx context.Context, path string, pairs ...Pair) (err error)
	...
	
	// Copy
	Copy(src string, dst string, pairs ...Pair) (err error)
	CopyWithContext(ctx context.Context, src string, dst string, pairs ...Pair) (err error)
	
	...
}
```

Service implementer needs to implement all the operations. `ErrCapabilityInsufficient` should be returned when an operation is not supported.

Caller need to check `Feature` before use them.

```go
if store.Feature().FeatureCopy() {
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

### Alternative Way

Instead of adding operations in all interfaces except the current `Servicer` and `Storager` back to `Storager`, maybe we can keep the current interfaces and generate the following interface for service according to `service.toml`:

```go
type Storage interface {
	Storager
	Appender
	Feature() Feature
}
```

Service implementer should implement all the supported interfaces.

## Compatibility

The following interfaces will be deprecated, and the related operations will be added back to `Storager` interface. Caller need to check `Feature` before use them.

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

## Implementation

- For go-storage
  - Generate `Feature`
  - Deprecate interfaces other than `Servicer` and `Storager`, and add the related operations back to `Storager`
- For storage services
  - Remove `features` and `implement` fields in `service.toml`
  - Implement `Feature() Feature` for `Storage` and `Service`
- For go-integration-test
  - go-integration-test should run the test cases according to `Feature`.