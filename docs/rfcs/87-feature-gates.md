- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-05-30
- RFC PR: [beyondstorage/specs#87](https://github.com/beyondstorage/specs/issues/87)
- Tracking Issue: N/A

# GSP-87: Feature Gates

## Background

Behavior consistency is the most important thing for go-storage. However, we do have different capabilities and limitations in storage services. So the problem comes into how to handle them.

Our goals are:

- Behavior consistent by default (invalid operations should be return error)
- Give user abilities to loose the restriction
- Allow user to enable the features they want

In [GSP-16], we introduced loose mode which is a global flag that controls the behavior when services meet the unsupported pairs.

- If `loose` is on: This error will be ignored.
- If `loose` is off: Storager returns a compatibility-related error.

However, we removed `loose` in [GSP-20] because we can't figure out [error could be returned too early while in loose mode](https://github.com/beyondstorage/go-storage/issues/233). And loose is so general that it affects nearly all behavior of Storager.

In [types: Implement pair policy](https://github.com/beyondstorage/go-storage/pull/453), we try to figure out this problem by introducing `PairPolicy`.

`PairPolicy` controls the behavior of pairs:

```go
type PairPolicy struct {
   All bool
   
   ...

   // pairs for interface Storager
   Create           bool
   Delete           bool
   List             bool
   ListListMode     bool
   Metadata         bool
   Read             bool
   ReadSize         bool
   ReadOffset       bool
   ReadIoCallback   bool
   Stat             bool
   Write            bool
   WriteContentType bool
   WriteContentMd5  bool
   WriteIoCallback  bool
}
```

If `PairPolicy.Write` is on, Storager will return `services.PairUnsupportedError` while meeting not supported pairs.

But it doesn't solve the problem:

- `PairPolicy` is generated in go-storage and can't reflect the capabilities in service.
- `PairPolicy` only used for pair capabilities check, and can't fix the problem introduced in [GSP-86]

## Proposal

So I propose to treat loose behavior consistency as a feature and introduce feature gates in [go-storage].

### Features

`Feature` means userland optional abilities provided by [go-storage].

`Copier` is not a `Feature`.

It's decided by service providers, the user can't enable `Copier` for a service.

`LooseOperation` is a `Feature`.

It's decided by the user, and service providers can't affect its behavior.

Every `Feature` SHOULD be introduced via [go-storage] RFC process.

[go-storage] will generate feature gates struct for service:

```go
type StorageFeatures struct {
   LooseOpeartionAll bool
   LooseOperationWrite bool
   
   VirtualOperationAll bool
   VirtualOperationCreateDir bool
}

func WithStorageFeatures(v StorageFeatures) Pair {
   return Pair{
      Key:   pairStorageFeatures,
      Value: v,
   }
}
```

User can use `WithStorageFeatures` or `WithStorageFeatures` to enable features they want while `NewServicer` or `NewStorager`. Those features CANNOT be changed during runtime.

### New Feature: Loose Operation

We will format `PairPolicy` as a new feature called: `Loose Operation`.

By default, [go-storage] will return errors for not supported pairs. If loose operation feature has been enabled, [go-storage] will ignore those errors.

To enable loose operation, users need to add pairs like:

```go
s3.WithStorageFeatures(s3.StorageFeatures{
   LooseOpeartionAll: true,
   LooseOperationWrite: true,
})
```

> `New` function is special, and we will always enable `loose` for it.

### New Feature: Virtual Operation and Virtual Pair

We will introduce a new idea `Virtual` to represents the following state: A service doesn't support some feature, but we can simulate it by some methods.


For example:

- We can have `Virtual Operation`
    - S3 doesn't have native support for `CreateDir`, but we can put an Object end with `/` to simulate it.
    - S3 doesn't have native support for `Link`, but we can store the link target in object metadata to simulate it.
- We can have `Virtual Pair`
    - fs doesn't support `content_md5` for `write`, but we can store it in the file's xattr to simulate it.

Those features will affect users data in some way:

- User couldn't read those data without go-storage
- Those data could be changed by other API and go-storage can't detect them.

So the end-user has to decide whether enable those virtual features or not.

- If they are enabled, our services will run in the `virtual` mode, likes your virtual machine. And they have to afford the side effects.
- If they don't, our services will behave like they don't implement this feature.

The virtual feature will give us more power so that we can implement `Link` / `Rename` even when our service doesn't have native support.

To enable virtual operation, users need to add pairs like:

```go
s3.WithStorageFeatures(s3.StorageFeatures{
   VirtualOperationAll: true,
})
```

## Rationale

N/A

## Compatibility

This proposal will deprecate `PairPolicy`.

## Implementation

To implement virtual operation and virtual pair support, we will add new fields into `service.toml`.

```toml
[namespace.storage.op.write]
simulated = true
optional = ["content_type", "io_callback", "storage_class"]
virtual = ["content_md5"]
```

- `simulated`: Mark this operation as a `virtual operation`.
- `virtual`: The list of `virtual pairs`.

[GSP-16]: ./16-loose-mode.md
[GSP-20]: ./20-remove-loose-mode.md
[GSP-86]: https://github.com/beyondstorage/specs/pull/86
[go-storage]: https://github.com/beyondstorage/go-storage
