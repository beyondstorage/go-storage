- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-07-04
- RFC PR: [beyondstorage/specs#134](https://github.com/beyondstorage/specs/pull/134)
- Tracking Issue: [beyondstorage/go-storage#624](https://github.com/beyondstorage/go-storage/issues/624)

# GSP-134: Write Behavior Consistency

Previous Discussions:

- [Add detailed behavior about move/copy](https://github.com/beyondstorage/specs/issues/130)
- [Make write Idempotent](https://github.com/beyondstorage/specs/issues/96)
- [copy/move doesn't need to support dst path that has an existing path](https://github.com/beyondstorage/go-integration-test/issues/29)

## Background

Behavior consistency is extremely critical to us, and write related operations is the most important part.

Write related operations including following operations:

- `Write`
- `CompleteMultipart`
- `CombineBlock`
- `CommitAppend`
- `Copy`
- `Move`

All of them will create a new object in our storage service. But our services could have different behavior while the destination path already exists. I will list a few representative service behaviors for analysis.

### fs

For `Write`:

`fs` could return `EEXIST` to indicate `File exists`. Users can set open mode `Create` to create a new file here without meeting `EEXIST` error.

For `Copy`:

`fs` doesn't have native copy support, we could implement it via reading from src and write into dst. For the writing part, most userland tools will use `Create` open mode to ignore `EEXIST` error.

So `Copy` will not return the `EEXIST` error.

For `Move`:

`fs` will ignore the `EEXIST` error and always replace the file. On the same filesystem, it's an atomic operation, otherwise, it should be implemented as `copy-delete`.

So `Move` will not return the `EEXIST` error.

### s3

Most s3 alike object storage services will not return object exist errors.

- For object storage services that don't have an object version, the last write will always succeed.
- For object storage services that have object versions, every write will create a new object version.

### ipfs

For `Write`:

`ipfs` supports `Create` like `fs`. With `Create`, `ipfs` will create the file if it does not exist, this is the only behavior, in other words, `ipfs` will continue to write to the file even if create is not specified.

If writing to an existing file, IPFS will overwrite it from the beginning, and if there are unwritten parts at the end of the original file, they will remain as well. `ipfs` provides an option `Truncate` to make sure that the write is clean.

```go
func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
   err = s.ipfs.FilesWrite(
      ctx, s.getAbsPath(path), r,
      ipfs.FilesWrite.Create(true),
      ipfs.FilesWrite.Parents(true),
      ipfs.FilesWrite.Truncate(true),
   )
   if err != nil {
      return 0, err
   }
   return size, nil
}
```

For `Copy` and `Move`:

`ipfs` have native support for `Copy` and `Move`. If the dst exists, the `Copy` operation will return `FileExists`, while the `Move` operation will overwrite.

### dropbox

For `Write`:

`dropbox` supports setting `WriteMode`. With `files.WriteModeAdd`, `Write` will not return `FileExists` too:

```go
input := &files.CommitInfo{
    Path: rp,
    Mode: &files.WriteMode{
        Tagged: dropbox.Tagged{
            Tag: files.WriteModeAdd,
        },
    },
}
```

For `Append`:

`dropbox` supports `WriteMode` like `Write`.

For `Copy` and `Move`:

`dropbox` supports `autorename`.

- If `autorename` is false, `dropbox` will return `to/conflict/file/` error.
- If `autorename` is true, `dropbox` will rename the final path with a suffix like `(1)`.


In [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md), we make `Delete` idempotent which means without any outsides changes, any `Delete` call on the same path will get the same result. But `Write` is more complex than `Delete` and we cannot reuse the experience from GSP-46 directly:

- `Write` is unlikely to be retried: the underlying `io.Reader` could only be consumed once.
- Make `Delete` idempotent is simple: we just need to ignore the `ObjectNotExist` error, but it's difficult for `Write`.

## Proposal

I propose to improve write behavior consistency in the following method:

- All write operations SHOULD NOT return an error as the object exists.
- All successful write operations SHOULD be complete

Compared to GSP-46, we will not enforce the requirement for idempotent, because most write operations are not idempotent. For example, the second call for `Move(a, b)` should report `ObjectNotExist`.

`Complete` means after a successful write, the object's content and metadata should be the same as specified in write request. For example, write `1024` bytes on an existing `2048` bytes file should result in a `1024` bytes file.

### Write Operations

Write operations are those that will create objects in the service, including:

- `Write`
- `Copy`
- `Move`
- `Fetch`
- `CreateDir`
- `CreateLink`
- `CreateMultipart`
- `CompleteMultipart`
- `CreateBlock`
- `CombineBlock`
- `CreateAppend`
- `CommitAppend`
- `CreatePage`

Other APIs are out of scope.

### Services Implementations

- Service that has native support for `overwrite` or `create` doesn't NEED to check the object exists or not.
- Service that doesn't have native support for `overwrite` or `create` SHOULD check and delete the object if exists.
  - `autoreanme` alike features are forbidden: service provider SHOULD NOT change the final destination path.
- Service SHOULD make sure write operations complete.

### Error Handling

If user call `Write(src, dst)` but dst is a dir, service SHOULD return `ErrObjectModeInvalid` error instead of

- removing the dst
- returning `ErrObjectExist`
- write into `dst/src`

Other write operations are the same.

## Rationale

### Check and return ObjectExist

For services

Adopting this method requires all services to check before write, even those services that never return `ObjectExist`.

For users

The API is unnatural to use.

If user check the existence by them self:

```go
_, err := store.Stat(dst)
if err == nil {
   err = store.Delete(dst)
   if err != nil {
      return err
    }
}
if err != nil && !errors.Is(err, ObjectNotExist) {
   ...
}
_, err := store.Copy(src, dst) // And we will check dst double times here.
```

If user doesn't care about whether dst exists:

```go
_, err := store.Copy(src, dst)
if err == nil {
   ...
}
if err != nil && !errors.Is(err, ObjectExist) {
   ...
}
_, err = store.Delete(dst)
if err != nil {
   ...
}
_, err := store.Copy(src, dst) // And we will check dst double times here.
if err != nil {
   ...
}
```

But with our current design:

```go
_, err := store.Copy(src, dst)
if err != nil {
   ...
}
```

If users do want to make sure the file does not exist before copy, they can `Stat` by themselves, and there is no extra call inside `Copy`.

### Add new pair to control the logic

Adding a new pair like `force` makes it more complex to maintain. To implement `force` correctly, we need to implement both the current design and `Check and return ObjectExist`.

## Compatibility

This proposal doesn't introduce any API changes.

## Implementation

- Specify the behavior in specs, docs and comments
- Add more tests in integration tests
- Make sure all services passed
