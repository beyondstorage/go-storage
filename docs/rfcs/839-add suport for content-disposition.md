- Author: Joey-1445601153,  <github@Joey-1445611153>
- Start Date: 2021-10-12
- RFC PR: [beyondstorage/go-storage#839](https://github.com/beyondstorage/go-storage/issues/839)
- Tracking Issue: [beyondstorage/go-storage#836](https://github.com/beyondstorage/go-storage/issues/836)

# GSP-839: Add support for Content-Disposition
Previous Disscussion:
- [Add support about Content-Disposition](https://forum.beyondstorage.io/t/topic/227)

## Background

The Content-Disposition header provides a mechanism, allowing each component of a message to be tagged with an indication of its desired presentation. It is wildly used by storage products, such as [Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-properties), [AWS](https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPOST.html),  [TencentCloud](https://cloud.tencent.com/developer/section/1189916), etc.

Now We don't support Content-Disposition. Add support for Content-Disposition will allow user to decide how to show message. 

## Proposal

So I propose following changes:

- Add content-disposition pair to global pairs
- Add content-disposition to object metadata
- Add process of content-disposition field
  - For write operation: User can use `content-disposition` to set the object metadata
  - For read operation: User can set `content-disposition` for this request

Add content-disposition pair to global pairs

Add content-disposition to object metadata

Add content-disposition process in read&write relevant operations if we need.

## Rational

- None

## Compatibility

No breaking changes.

## Implementation

- Add content-disposition pair to pairs.toml in go-storage
- Add content-disposition pair to info_object_meta.toml in go-storage
- Add process of content-disposition field in read&write relevant operations for necessery service
