---
author: Xuanwo <github@xuanwo.io>
status: candidate
updated_at: 2020-01-03
---

# Proposal: Normalize Metadata

## Background

Metadata is being chaos.

Currently, we have following metadata:

```go
const (
	Checksum = "checksum"
	Class    = "class"
	Count    = "count"
	Expire   = "expire"
	Host     = "host"
	Location = "location"
	Size     = "size"
	Type     = "type"
	WorkDir  = "work_dir"
)
```

In order to unify different value name across services, we have a virtual name map:

- checksum -> Content-MD5
- type -> Content-Type
- class -> x-qs-storage-class, x-amz-storage-class...
- size -> bucket's size
- count -> bucket's count
- location -> bucket's location or zone

After more values prompt to struct value, these name map could confuse end users:

- `object.Size` or `object.GetSize()` ?
- `object.Type` or `object.GetType()` ?
- Different services could use different hash algorithm, how to distinguish them?

## Proposal

So I propose following changes:

### Split storage meta and object meta

- Use a `map[string]interface{}` to replace `Metadata`, and generate functions on `Storage` directly.
- Rename `metadata.Storage` to `metadata.StorageMeta`
- Add `metadata.ObjectMeta`
- Add `metadata.StorageStatistic`
- Split `metadata.json` into `object_meta.json`, `storage_meta.json` and `storage_statistic.json`
- Refactor `metadata` code generator
- Make `SetXXX` return itself, so we can call them like a chain:

    ```go
    m := metadata.NewObjectMeta().
        SetContentType("application/json").
        SetStorageClass("cool").
        SetETag("xxxxx")
    ```

### Normalize metadata's name

We should normalize metadata name in following ways:

**If this metadata has been output as `Object` struct value, ignore them.**

**If this metadata defined in [Message Headers](https://www.iana.org/assignments/message-headers/message-headers.xhtml), use the normative style.**

For example, `content-md5`. 

*As HTTP/2 has been approved in [rfc7540](https://tools.ietf.org/html/rfc7540), we will use the lower case style of header.*

**If this metadata is a private meta, use the most common usage case.**

For example, `x-qs-storage-class` and `x-amz-storage-class` should be meta `storage-class`.

### Normalize metadata's value

Those meta value should also be normalized: every service should convert their own meta to storage shared meta.

For example, Amazon S3 have following storage class: `STANDARD`, `REDUCED_REDUNDANCY`, `STANDARD_IA`, `ONEZONE_IA`, `GLACIER`, `DEEP_ARCHIVE` and so on.

In Azure Storage, they have different `AccessTier`: `AccessTierHot`, `AccessTierCool`, `AccessTierArchive`. 

We need to handle all of them: 

- Store them as meta: `storage-class`
- Convert to storage shared meta while fetch from service
- Convert to service private meta while updating

## Rationale

None.

## Compatibility

Metadata's name and value could be changed.

## Implementation

Most of the work would be done by the author of this proposal.