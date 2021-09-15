- Author:  Abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-09-10
- RFC PR: [beyondstorage/go-storage#751](https://github.com/beyondstorage/go-storage/pull/751)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-751: Write Empty File Behavior

- Previous discussion:
  - [Specify the behavior for writing empty file](https://forum.beyondstorage.io/t/topic/204)

## Background

```go
func (s *Storage) Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {}
```

In our definition, the `Write` function will upload a file to the path. We do not do anything about uploading empty files(i.e. `io.Reader`is nil ). However, some services do not support uploading nil `io.Reader`. As [proposed](https://forum.beyondstorage.io/t/topic/204) in our community, for s3, when we upload an empty file, the service will hang permanently. Here is an example of s3:

```go
_, err = store.Write(path,nil, 0)
if err != nil {
return err
}
```

If we want to upload an empty file, we need to do so:

```go
_, err = store.Write(path, bytes.NewReader([]byte{}), 0)
if err != nil {
    return err
}
```

This is not convenient for the users.

## Proposal

I propose to allow the user to pass in a nil `io.Reader` and `0` size to create an empty file when calling `Write`. For services that do not support nil Reader uploads we can do something like this.

```go
func (s *Storage) Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error) {
// We can create a Reader with size 0 like this.
r = bytes.NewReader([]byte{})
// Then we can call the API with r and 0 size.
}
```

- What will happen if we got a nil `io.Reader` but `size != 0`?
  - We will return an error. We can't make decisions for users.
- What will happen if we got a valid `io.Reader` but `size = 0`?
  - We will upload files of `io.Reader` length.
  - If the upload is successful, we will return the size as the length of the `io.Reader`.

## Rationale

N/A

## Compatibility

This change will not break services and users. We can do this as follows:

- Add this behavior to `Write` in `go-storage`.
- Services bump to the latest version of `go-storage' (if the latest version does not contain this behavior, bump to the latest master branch), and all descriptions of this behavior are updated. Then changes are made to the services that need to be modified.

## Implementation

- `go-storage`
  - Update the description of `Write` in `definitions` to include this behavior.
- `go-integration-test`
  - Add the following cases to the integration tests in `storager`.
    - When write a file with a nil `io.Reader` and `0` size
    - When write a file with a nil `io.Reader` and `non-zero` size
    - When write a file with a non-nil `io.Reader` and `0` size
- `go-service-*`
  - Implement this behavior.
  - Ensure that integration tests pass.