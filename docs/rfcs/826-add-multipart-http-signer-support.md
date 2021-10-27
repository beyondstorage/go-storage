- Author: abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-09-26
- RFC PR: [beyondstorage/go-storage#826](https://github.com/beyondstorage/go-storage/pull/826)
- Tracking Issue: [beyondstorage/go-storage#827](https://github.com/beyondstorage/go-storage/issues/827)

# GSP-826: Add Multipart HTTP Signer Support

- Updated By:
  - [GSP-837: Support Feature Flag](./837-support-feature-flag.md): Move multipart HTTP signer related operations to `Storager`

Previous discussion:

- [Multipart related operations ](https://forum.beyondstorage.io/t/topic/226)

# Background

In [GSP-729](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/729-redesign-http-signer.md), we split out the `HTTPSigner` proposed in [GSP-706](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/706-support-http-signer.md). We had the following problems in the implementation of GSP-706:

- There's no appropriate way to pass in some parameters for some operations like [How to pass partIndex into QuerySignHTTP for WriteMultipart](https://forum.beyondstorage.io/t/how-to-pass-partindex-into-querysignhttp-for-writemultipart/192).
- `CreateMultipart` needs to support Expires parameter.

## Proposal

I propose to add `MultipartHTTPSigner` interface.

```go
type MultipartHTTPSigner interface {
    QuerySignHTTPCreateMultipart(path string, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
    QuerySignHTTPWriteMultipart(o *Object, size int64, index int, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
    QuerySignHTTPListMultipart(o *Object, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
    QuerySignHTTPCompleteMultipart(o *Object, parts []*Part, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
}
```

- `MultipartHTTPSigner` is the interface associated with `Multipart`, which support using query parameters to authenticate requests.
  - `QuerySignHTTPCreateMultipart`, `QuerySignHTTPWriteMultipart`, `QuerySignHTTPListMultipart` and `QuerySignHTTPCompleteMultipart` are the supported signature operations for `createMultipart`, `writeMultipart`, `listMultipart` and `completeMultipart` in `MultipartHTTPSigner` for now.
- Compared to the corresponding basic operation (`createMultipart`, `writeMultipart`, `listMultipart` and `completeMultipart` in `Multiparter`), the parameters of the `MultipartHTTPSigner` operations have the following differences:
  - `expire` is required.
  - `io.Reader` typed parameter for writeMultipart operations SHOULD be removed.
  - Other parameters SHOULD be consistent.

Also I propose to support `delete` in `StorageHTTPSigner`(Proposed in [GSP-729](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/729-redesign-http-signer.md#proposal)).

```go
type StorageHTTPSigner interface {
    QuerySignHTTPDelete(path string, expire time.Duration, ps ...types.Pair) (req *http.Request, err error)
}
```

- `QuerySignHTTPDelete` is the supported signature operation for `delete` in `StorageHTTPSigner` for now.
- Compared to the corresponding basic operation (`delete` in `Storager`),  the parameters of the `QuerySignHTTPDelete` operation has the following differences:
  - `expire` is required.
  - Original parameters SHOULD be consistent.

From service side:

- If part of the operations are unsupported, `services.ErrCapabilityInsufficient` error can be returned directly.

## Rationale

Rationale is mentioned in the [GSP-729](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/729-redesign-http-signer.md#rationale).

## Compatibility

N/A

## Implementation

- `go-storage`
  - Add new interface and operations in definitions.
- `go-integration-test`
  - Add integration tests for `MultipartHTTPSigner`. Create a new file `multipart_http_signer.go`. Split into the following four functions:
    - `TestMultipartHTTPSignerCreateMultipart(t *testing.T, store types.Storager) {}`
    - `TestMultipartHTTPSignerWriteMultipart(t *testing.T, store types.Storager) {}`
    - `TestMultipartHTTPSignerListMultipart(t *testing.T, store types.Storager) {}`
    - `TestMultipartHTTPSignerCompleteMultipart(t *testing.T, store types.Storager) {}`
  - Add integration tests for `QuerySignHTTPDelete` in `storage_http_signer.go`.
    - `TestStorageHTTPSignerDelete(t *testing.T, store types.Storager) {}`
- `go-service-*`
  - Update all services that support HTTP signer.
  - Ensure that integration tests pass.