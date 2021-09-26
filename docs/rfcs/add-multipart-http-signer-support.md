- Author: abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-09-26
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-0: Add Multipart HTTP Signer Support

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
    signHTTPCreateMultipart(path string, expire time.Duration, ps ...types.Paires) (req *http.Request, err error)
    signHTTPWriteMultipart(o *Object, size int64, index int, expire time.Duration, ps ...types.Pairs) (req *http.Request, err error)
}
```

- `MultipartHTTPSigner` is the interface associated with `Multipart`, which support using query parameters to authenticate requests.
  - `SignHTTPCreateMultipart` and `SignHTTPWriteMultipart` are the supported signature operations for `createMultipart` and `writeMultipart` in `MultipartHTTPSigner` for now.
- Compared to the corresponding basic operation (`createMultipart` and `writeMultipart` in `Multiparter`), the parameters of the `MultipartHTTPSigner` operations have the following differences:
  - `expire` is required.
  - `io.Reader` typed parameter for writeMultipart operations SHOULD be removed.
  - Other parameters SHOULD be consistent.

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
  - Add integration tests for `MultipartHTTPSigner`. Split into the following two functions:
    - `TestMultipartHTTPSignerCreateMultipart(t *testing.T, store types.Storager) {}`
    - `TestMultipartHTTPSignerWriteMultipart(t *testing.T, store types.Storager) {}`
- `go-service-*`
  - Update all services that support HTTP signer.
  - Ensure that integration tests pass.
