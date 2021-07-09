- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-05-08
- RFC PR: [beyondstorage/specs#49](https://github.com/beyondstorage/specs/issues/49)
- Tracking Issue: N/A

# AOS-49: Add CreateDir Operation

## Background

Applications need the ability to create a directory. For now, our support is a bit wired.

In [fs](https://github.com/beyondstorage/go-service-fs), we support `CreateDir` by dirty hack:

```go
if s.isDirPath(rp) {
    // FIXME: Do we need to check r == nil && size == 0 ?
    return 0, s.createDir(rp)
}
```

In other storage service, user needs to create dir by special `content-type`:

```go
store.Write("abc/", nil, 0, ps.WithContentType("application/x-directory"))
```

We need to allow user create a directory in the same way.

## Proposal

So I propose to add a new operation `CreateDir` like we do on `append` / `multipart` Object.

```go
type Direr interface {
	CreateDir(path string, pairs ...Pair) (o *Object, err error)
}
```

`CreateDir` will return an Object with `dir` mode, and different service could have different implementations.

## Rationale

### Directory in Object Storage Services

Object Storage is a K-V Storage, and don't have the concept of directory natively. But most object storages support ListObjects via delimiter `/` to demonstrate a file system tree. With delimiter `/`, object storage services will organize objects end with `/` as common prefix.

## Compatibility

This proposal COULD break users who use `store.Write("abc/")` to create directory.

## Implementation

- Update specs to add `CreateDir` operations.
- Update go-storage to implement the changes.
- Update all services that support create dir.
