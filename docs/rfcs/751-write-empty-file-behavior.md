- Author:  Abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-09-10
- RFC PR: [beyondstorage/go-storage#751](https://github.com/beyondstorage/go-storage/pull/751)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-751: Write Empty File Behavior

- Previous discussion:
    - [Specify the behavior for writing empty file](https://forum.beyondstorage.io/t/topic/204)

## Background

```go
func (s *Storage) Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
```

In our definition, the `Write` API will upload a file to the path. We do not do anything about uploading empty files(i.e. `io.Reader`is nil ). However, some services do not support uploading nil `io.Reader`. As [proposed](https://forum.beyondstorage.io/t/topic/204) in our community, for s3, when we upload an empty file, the service will hang permanently. Here is an example of s3:

```go
_, err = s3.s.Write(path,nil, 0)
if err != nil {
    s3.logger.Error("write", zap.String("path", path), zap.Error(err))
    return nil, nil, err
}
```

If we want to upload an empty file, we need to do so:

```go
_, err = s3.s.Write(path, bytes.NewReader([]byte{}), 0)
if err != nil {
    s3.logger.Error("write", zap.String("path", path), zap.Error(err))
    return nil, nil, err
}
```

This is not convenient for the users.

## Proposal

I propose to allow the user to pass in a nil `io.Reader` when calling `Write`.

- For services that do not support uploading nil `io.Reader`, but support upload object with size 0 , we should check if `io.Reader` is `nil`. If it is `nil`, we need to create a `io.Reader` before calling the API. For example, `s3`, `kodo`, etc.
- For services that support uploading nil `io.Reader`, we don't need to check if the `reader` is `nil`, we can call the API directly. Like `oss`, etc.
- For services that do not support uploading both nil `io.Reader` and objects of size 0, we should mark that they do not support uploading empty files, and check whether `io.Reader` is nil or size 0. Like `azfile`, etc.

## Rationale

N/A

## Compatibility

This behavior will affect all services that do not support uploading empty files.

## Implementation

- Update definitions to reflect changes
- Update integrations tests to make sure all service passed

