- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-27
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP: Redesign HTTP Signer

- Updates:
  - [GSP-706]: Deprecate `QuerySignHTTP()`.

Previous discussion:
- https://github.com/beyondstorage/go-storage/issues/726
- https://matrix.to/#/#beyondstorage@gsp-706:matrix.org

## Background

## Proposal

Split `QuerySignHTTP()` into multiple operations according to `op`:

```go
type HTTPSigner interface {
    QuerySignHTTPRead(path string, expire time.Duration, ps ...types.Pair) (signedReq *http.Request, err error)
    QuerySignHTTPWrite(path string, size int64, expire time.Duration, ps ...types.Pair) (signedReq *http.Request, err error)
    // ...
}
```

- `Read` and `Write` (add `Multiparter` related operations if needed) are the supported operations for now. Other operations SHOULD be introduced by new GSP.
- Compared to the corresponding basic operation, the parameters of the `HTTPSinger` operations have the following differences:
  - `expire` is required.
  - `io.Reader` typed parameter for write operations or `io.Writer` typed parameter for read operations SHOULD be removed.
  - Other parameters SHOULD be consistent.
  
From service side:

- If part of the operations are unsupported, `services.ErrCapabilityInsufficient` error can be returned directly.

## Rationale

### Query string authorization in s3

As described in the [authentication overview](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html#auth-methods-intro), you can provide authentication information using query string parameters.

Take `aws-sdk-go` as an example, it provides `xxxRequest()` to generate a "aws/request.Request" representing the client's request for the corresponding operation. Then for query string authorization, we can use `Presign` to get the request's signed URL or `PresignReqest` to get the signed url and a set of header that were signed.

As for go-storage, if only s3 service is considered, we can add signature operations for all the interfaces, like: 

```go
type Storager interface {
    // ...
    // Delete will delete an object from service.
    Delete(path string, pairs ...Pair) (err error)
    // ...
    // HTTP signer for Delete operation.
    QuerySignHTTPDelete(path string, expire time.Duration, pairs ...Pair) (signedReq *http.Request, err error)

    // Read will read the file's data.
    Read(path string, w io.Writer, pairs ...Pair) (n int64, err error)
    // ...
    // HTTP signer for Read operation.
    QuerySignHTTPRead(path string, expire time.Duration, pairs ...Pair) (signedReq *http.Request, err error)

    // Write will write data into a file.
    Write(path string, r io.Reader, size int64, pairs ...Pair) (n int64, err error)
    // ...
    // HTTP signer for Write operation.
    QuerySignHTTPWrite(path string, size int64, expire time.Duration, pairs ...Pair) (signedReq *http.Request, err error)

    // ...
}
```

But it will generate a lot of unsupported operations for many services.

### SignedURL in gcs

```go
// Methods which can be used in signed URLs.
var signedURLMethods = map[string]bool{"DELETE": true, "GET": true, "HEAD": true, "POST": true, "PUT": true}

// SignedURLOptions allows you to restrict the access to the signed URL.
type SignedURLOptions struct {
    // ...
    // SignBytes is a function for implementing custom signing.
    SignBytes func([]byte) ([]byte, error)
    // Method is the HTTP method to be used with the signed URL.
    Method string
    // Expires is the expiration time on the signed URL.
    Expires time.Time
    // ContentType is the content type header the client must provide to use the generated signed URL.
    // Optional.
    ContentType string
    // Headers is a list of extension headers the client must provide in order to use the generated signed URL.
    Headers []string
    // QueryParameters is a map of additional query parameters.
    QueryParameters url.Values
    // MD5 is the base64 encoded MD5 checksum of the file.
    // Optional.
    MD5 string
    // ...
}

// SignedURL returns a URL for the specified object.
func SignedURL(bucket, name string, opts *SignedURLOptions) (string, error) {}
```

Combining the design and application scenarios of go-storage, although unsupported operations could be avoided, the drawbacks are obvious:

- A suitable struct compatible with most services is needed, so we have to think more for it.
- It's complex for the implementation in services, and also difficult to maintain.
- It's hard for users as they need to know how to pass in the parameters correctly for different services.

### Authorized access in oss

You can generate a signed URL and provide the URL to a visitor for temporary access:

```go
// SignURL signs the URL.
func (bucket Bucket) SignURL(objectKey string, method HTTPMethod, expiredInSec int64, options ...Option) (string, error) {}
```

"oss/Option" is HTTP options and used to set URL parameter, HTTP header and function argument. `PutObjectWithURL` or `GetObjectWithURL` is used for uploading or downloading an object with the URL.

## Compatibility

`QuerySignHTTP()` will be deprecated.

## Implementation

N/A

[GSP-706]: ./706-support-http-signer.md
