- Author: JinnyYi <github.com/JinnyYi>
- Start Date: 2021-08-16
- RFC PR: [beyondstorage/go-storage#706](https://github.com/beyondstorage/go-storage/issues/706)
- Tracking Issue: [beyondstorage/go-storage#707](https://github.com/beyondstorage/go-storage/issues/707)

# GSP-706: Support HTTP Signer

Previous discussion:

- [Allow generate signed url for upload objects](https://github.com/beyondstorage/go-storage/issues/646)

## Background

Authentication is the process of proving user's identity to the system. In addition to add signatures to the `Authorization` header of requests, users can also add signatures to the URL of the resource.

A signed URL is a URL that provides limited permission and time to make a request. Signed URLs contain authentication information in their query string. Using query parameters to authenticate requests is useful when users want to express a request entirely in a URL. A use case scenario for signed URL is that users can grant access to the resource.

## Proposal

I propose to add the following interface containing operations that support the generation of signed URL for RESTful services:

```go
type HTTPSigner interface {
    QuerySignHTTP(op, path string, expire time.Duration, ps ...types.Pair) (signedReq *http.Request, err error)
}
```

`HTTPSigner` is the interface for `Signer` related operations which support calculating request signature.

`QuerySignHTTP` returns `*http.Request` with query string parameters containing signature in `URL` to represent the client's request for the specified operation.

**Parameters**

- op: is a const string representing operation name defined in `types` package.
  - `op` SHOULD be the supported operation by service.
- path: is the path of object.
  - `path` COULD be relative or absolute path.
- expire: provides the time period, with type `time.Duration`, for which the generated `signedReq.URL` is valid.
  - Different services have different valid value ranges for `expire`.
- ps: is the arguments for this operation.

**Returns**

- signedReq: represents an HTTP request to be sent by service.
  - `URL` SHOULD NOT be nil and SHOULD be the request's signed URL.
- err: returning error if errors are encountered. It's nil if no error.

From service side:

- Services SHOULD maintain the supported authorized access operation list and check the validity of `op`.
- Services SHOULD return `http.Request` pointer with signature in the query string of `URL`, which is constructed by specific storage service.

From user side:

- A clock calibration is required for validation of expiration.

## Rationale

N/A

## Compatibility

This proposal will deprecate `Reacher` interface.

## Implementation

- Add new interface and operations in definitions.
- Update docs in site.
- Implement integration test.
- Implement `HTTPSigner` for services.
