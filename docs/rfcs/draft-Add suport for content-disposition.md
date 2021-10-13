- Author: Joey-1445601153,  <github@Joey-1445611153>
- Start Date: 2021-10-12
- RFC PR: [beyondstorage/go-storage#xx](https://github.com/beyondstorage/go-storage/issues/xx)
- Tracking Issue: [beyondstorage/go-storage#836](https://github.com/beyondstorage/go-storage/issues/836)

# GSP-xx: Add support for Content-Disposition

## Background

The Content-Disposition header provides a mechanism, allowing each component of a message to be tagged with an indication of its desired presentation. It is wildly used by storage products, such as [Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-properties), [AWS](https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPOST.html),  [TentcentCloud](https://cloud.tencent.com/developer/section/1189916), etc.

Now We don't support Content-Disposition. Add support for Content-Disposition will allow user to decide how to show message. 

## Proposal

So I propose following changes:

- Add content-disposition pair to global pairs
- Add process of content-disposition field
  - For write operation: User can use `content-disposition` to set the object metadata
  - For read operation: User can set `content-disposition` for this request

Add content-disposition pair to global pairs

Add property for pair:

```
[content_disposition]
type = "string"
```

Add process of content-disposition field

Add content-disposition in service.tomal in go-service-* 

Add content-disposition process at read&write operation like:

```
if opt.HasContentDisposition {
	input.ContentDisposition = &opt.ContentDisposition
}
```

## Rational

- None

## Compatibility

No breaking changes.

## Implementation

- Add content-disposition pair to pairs.toml in go-storage
- Add content-disposition pair to info_object_meta.toml in go-storage
- Add process of content-disposition field in read&write relevant operation for each service

## Previous Discussion

This proposal came out of [Beyondstorage forum Topic#227](https://forum.beyondstorage.io/t/topic/227)

