- Author: Joey-1445601153 <github.com/Joey-1445601153>
- Start Date: 2021-10-12
- RFC PR: [beyondstorage/go-storage#839](https://github.com/beyondstorage/go-storage/issues/839)
- Tracking Issue: [beyondstorage/go-storage#852](https://github.com/beyondstorage/go-storage/issues/852)

# GSP-839: Add Support for Content-Disposition

Previous Disscussion:
- [Add support about Content-Disposition](https://forum.beyondstorage.io/t/topic/227)

## Background

The `Content-Disposition` header provides a mechanism, allowing each component of a message to be tagged with an indication of its desired presentation. It is wildly used by storage products, such as [Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-properties), [AWS](https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPOST.html),  [TencentCloud](https://cloud.tencent.com/developer/section/1189916), etc.

Now we don't support `Content-Disposition`. Add support for `Content-Disposition` will allow user to specify the field. 

## Proposal

So I propose following changes:

- Add `content-disposition` pair to global pairs
- Add `content-disposition` to object metadata
- Add process of `content-disposition` field
  - For write operation: User can use `content-disposition` to set the object metadata
  - For read operation: User can set `content-disposition` for this request

### Write with content-disposition

User can take write operation with `content-disposition` as example:

```go
n, err := store.Write(path, r, length, pairs.WithContentDisposition("<content-disposition>"))
```

After write operation with `content-disposition`, presentational information of the object will be specified.

### Read with content-disposition

User can take read operation with `content-disposition` as example:

```go
n, err := store.Read(path, w, pairs.WithContentDisposition("<content-disposition>"))
```

After read operation with `content-disposition`, `content-disposition` filed in response header will be the value that is used in read operation. 

## Rational

N/A

## Compatibility

No breaking changes.

## Implementation

### go-storage implementation

- Add `content-disposition` pair to pairs.toml in go-storage
- Add `content-disposition` to info_object_meta.toml in go-storage

### service implementation

- Add `content-disposition` to optional pairs
- Add process of `content-disposition` field in read&write relevant operations
