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

Next are some situations that users will encounter when calling `Write`.

### nil io.Reader and zero size

- As described above, we will create a `io.Reader` of size 0 and upload it.
- We will upload data of length `size` as long as the length of `io.Reader` is the same as `size`.

### nil io.Reader and valid size

- We will return an error. In this case, the user's action is wrong, so we should return an error to alert the user.

### valid io.Reader and zero size

- We will upload an empty file with 0 size. We should follow the user's wishes and upload data of size 0.
- If the upload is successful, we will return the size of `0`.

### valid io.Reader and valid size

- If the size is smaller than the length of `io.Reader`, we will upload a file of size. If the upload is successful, size is returned.
- If the size is larger than the length of `io.Reader`, we will return an error.

## Rationale

N/A

## Compatibility

This change will not break services and users. We can do this as follows:

- Add this behavior to `Write` in `go-storage`.
- Add integration tests in `go-integration-test`.
- All services implement this behavior.

## Implementation

- `go-storage`
  - Update the description of `Write` in `definitions` to include this behavior.
- `go-integration-test`
  - Add the following cases to the integration tests in `storager`.
    - When write a file with a nil `io.Reader` and `0` size
    - When write a file with a nil `io.Reader` and `non-zero` size
    - When write a file with a non-nil `io.Reader` and `0` size
    - When write a file with a valid `io.Reader` and length greater than size
- `go-service-*`
  - Implement this behavior.
  - Ensure that integration tests pass.