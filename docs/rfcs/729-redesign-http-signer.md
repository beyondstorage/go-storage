- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-27
- RFC PR: [beyondstorage/go-storage#729](https://github.com/beyondstorage/go-storage/issues/729)
- Tracking Issue: [beyondstorage/go-storage#731](https://github.com/beyondstorage/go-storage/issues/731)

# GSP-729: Redesign HTTP Signer

- Updates:
  - [GSP-706]: Deprecate `HTTPSigner` and `QuerySignHTTP()`.

- Updated By:
  - [GSP-837: Support Feature Flag](./837-support-feature-flag.md): Move HTTP signer related operations to `Storager`

Previous discussion:
- [How to pass `partIndex` into `QuerySignHTTP` for `WriteMultipart`?](https://forum.beyondstorage.io/t/how-to-pass-partindex-into-querysignhttp-for-writemultipart/192)
- [#beyondstorage@gsp-706:matrix.org](https://matrix.to/#/#beyondstorage@gsp-706:matrix.org)

## Background

In [GSP-706], we introduced `HTTPSigner` interface to support HTTP authenticating requests, which contains `QuerySignHTTP()`. `QuerySignHTTP()` is used to authenticate requests by using query parameters.

During the implementation for services, we found the following problems:
- There's no appropriate way to pass in some parameters for some operations like [How to pass partIndex into QuerySignHTTP for WriteMultipart](https://forum.beyondstorage.io/t/how-to-pass-partindex-into-querysignhttp-for-writemultipart/192).
- Supporting all the authenticating request operations in one function makes it slightly complicated and hard to maintain, especially for the services that support query string authentication like s3, gcs, etc.

## Proposal

I propose to split `HTTPSigner` into multiple interfaces according to the current interface classification, and define HTTP signer interfaces one by one. For now, we will introduce the following interface:

```go
type StorageHTTPSigner interface {
    QuerySignHTTPRead(path string, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
    QuerySignHTTPWrite(path string, size int64, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
}
```

- `StorageHTTPSigner` is the interface associated with `Storager`, which support using query parameters to authenticate requests.
  - `QuerySignHTTPRead()` and `QuerySignHTTPWrite()` are the supported signature operations for read and write in `StorageHTTPSigner` for now. Other operations SHOULD be introduced by new GSP.
- Other interfaces like `MultipartHTTPSigner`, `AppendHTTPSigner` etc associated with `Multiparter`, `Appender` etc SHOULD be introduced by new GSP if needed.
- Compared to the corresponding basic operation(`Read` and `Write` in `Storager`), the parameters of the `StorageHTTPSigner` operations have the following differences:
  - `expire` is required.
  - `io.Reader` typed parameter for write operations or `io.Writer` typed parameter for read operations SHOULD be removed.
  - Other parameters SHOULD be consistent.
  
From service side:

- If part of the operations are unsupported, `services.ErrCapabilityInsufficient` error can be returned directly.

## Rationale

### Query string authentication in s3

As described in the [authentication overview](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html#auth-methods-intro), you can provide authentication information using query string parameters.

Take `aws-sdk-go` as an example, it provides `xxxRequest()` to generate a "aws/request.Request" representing the client's request for almost all the operations, like:

```go
// PutObjectRequest generates a "aws/request.Request" representing the client's request for the PutObject operation.
func (c *S3) PutObjectRequest(input *PutObjectInput) (req *request.Request, output *PutObjectOutput) {}
```

For query string authentication, we can use `Presign` to get the request's signed URL or `PresignReqest` to get the signed url and a set of header that were signed.

```go
// Presign returns the request's signed URL.
func (r *Request) Presign(expire time.Duration) (string, error) {}

// PresignRequest behaves just like presign, with the addition of returning a set of headers that were signed.
func (r *Request) PresignRequest(expire time.Duration) (string, http.Header, error) {}
```

As for go-storage, if only s3 service is considered, we can add signature operations for all the public APIs, like: 

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

Signed URL is to authenticate an HTTP request to cloud storage.

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

### Add signatures to a URL in oss

You can generate a signed URL and provide the URL to a visitor for temporary access:

```go
// SignURL signs the URL.
func (bucket Bucket) SignURL(objectKey string, method HTTPMethod, expiredInSec int64, options ...Option) (string, error) {}
```

"oss/Option" is HTTP options and used to set URL parameter, HTTP header and function argument, some of them will be involved in authentication.

You can add signatures to URLs that are contained in PUT and GET requests. `PutObjectWithURL` and `GetObjectWithURL` are the public APIs for uploading or downloading an object with signed URL.

```go
// PutObjectWithURL uploads an object with the signed URL.
func (bucket Bucket) PutObjectWithURL(signedURL string, reader io.Reader, options ...Option) error

// GetObjectWithURL downloads the object and returns the reader instance,  with the signed URL.
func (bucket Bucket) GetObjectWithURL(signedURL string, options ...Option) (io.ReadCloser, error)
```

## Compatibility

`HTTPSigner` and `QuerySignHTTP()` will be deprecated.

## Implementation

- Mark `HTTPSigner` and `QuerySignHTTP()` as deprecated.
- Add new interface and operations in definitions.
- Update all services that support HTTP signer.

[GSP-706]: ./706-support-http-signer.md
