- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-07-07
- RFC PR: [beyondstorage/go-storage#654](https://github.com/beyondstorage/go-storage/pull/654)
- Tracking Issue: [beyondstorage/go-storage#657](https://github.com/beyondstorage/go-storage/issues/657)

# GSP-654: Unify List Behavior

Previous Discussion:

- [Specify the behavior of List](https://github.com/beyondstorage/specs/issues/135)

## Background

[go-storage] exposes a list operation that lets you enumerate the keys contained in a bucket, or the files contained in a directory.

The list operation has been changed many times and tended to be stable since [Proposal: Unify List Operation](./23-unify-list-operation.md)

```go
func List(path string, pairs ...Pair) (oi *ObjectIterator, err error)
```

**Parameters**

path - The directory path for file system, or a file hosting service like dropbox. Also, it could be a prefix filter for object storage.

pairs - It contains all the filters for `List`. Usually we can specify `ListMode` for list operation.

`ListMode` is the type for `List`.

- `ListModePrefix` means this list will use prefix type. The returned file or object names must contain the prefix.
- `ListModeDir` means this list will use dir type. That means list files or objects hierarchically.
- `ListModePart` means this list will use part type. Generally, it's used to retrieve a list of in-progress multipart uploads. Only services that support multipart upload could list keys with `ListModePart`.
- `ListModeBlock` means this list will use block type. Generally, it's used to retrieve the list of blocks. It's used for block related operations.

**Returns**

oi - An object iterator. You can retrieve all the objects by `Next` and `IterateDone` will be returned while there are no items anymore.

err - It's nil if no error.


For the current implementation:

- Services have no default value for `ListMode` and will get `ListModeInvalidError` when no `ListMode` passed in which is not user-friendly.
- Whether we need to check `VirtualDir` while `ListMode` is `ListModeDir`. In [GSP-109: Redesign Features](./109-redesign-features.md), we introduced a new feature called `VirtualDir` to support simulated dir behavior and any operation related to dir should check `VirtualDir` in services that don't have native support for dir.
- Which fields should be assigned and how to handle object mode for the returning object iterator.

## Proposal

So I propose to specify the behavior of `List`.

### Support default `ListMode` for services

Service SHOULD support default `ListMode` if there's no `ListMode` passed in.

- File hosting services could be `ListModeDir`.
- Object storage services could be `ListModePrefix`.

### Implement `ListModeDir` without checking `VirtualDir`

Service SHOULD support `ListModeDir` for `List` and implement it without the check for `VirtualDir`.

Services like s3, oss or azblob, have a flat structure instead of a hierarchy like file system. However, for the sake of organizational simplicity, they support the `folder` concept as a means of grouping objects. The purpose of the prefix and delimiter is to help you organize and then browse the keys hierarchically. If prefix is specified and delimiter is set to a forward slash (/), only the objects in the "folder" are listed.

### Following [Object Lazy Stat]

There's no need `Stat` in `List`.

[go-storage] support [Object Lazy Stat], which allow the user to get the object information only when they really needed. Only `ID`, `Path`, `ObjectMode` are required during list operation, other fields could be fetched via lazy stat logic.

### How to Handle Object Mode?

[ObjectMode](https://beyondstorage.io/docs/go-storage/internal/core-concept#object) describes what users can operate on this object.

- For `ListModePrefix`:
    - Currently, it's only supported in object storage services, and the object mode usually have `ModeRead`.
- For `ListModeDir`:
    - For object storage, objects whose names contain the same string from the prefix and the next occurrence of the delimiter are grouped in `CommonPrefixes`, objects in the `CommonPrefixes` could be looked as folder, so the object mode should be `ModeDir`. Others should be `ModeRead`.
    - For file hosting services, we could set object mode according to the file attributes. If the object identifies a directory, it should have `ModeDir`. Or the object is a symbolic link, it should have `ModeLink`. Or the object identifies a regular file, it should have `ModeRead`, `ModeAppend` and `ModePage`.
- For `ListModePart`, the objects are in-progress multipart uploads, so object mode should be `ModePart`.
- For `ListModeBlock`, the objects are the list of blocks that have been uploaded, so object mode should be `ModeBlock`.

## Rationale

N/A

## Compatibility

Service should `List` with default `ListMode` instead of returning `ListModeInvalidError` if no `ListMode` passed in.

## Implementation

- All services should support default `ListMode`.
- All services should implement `ListModeDir`.

[go-storage]: https://github.com/beyondstorage/go-storage
[Object Lazy Stat]: https://beyondstorage.io/docs/go-storage/internal/object-lazy-stat
