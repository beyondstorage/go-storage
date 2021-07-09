- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-04-23
- RFC PR: [beyondstorage/specs#40](https://github.com/beyondstorage/specs/issues/40)
- Tracking Issue: N/A

# AOS-40: Unify Object Metadata

- Updated By:
  - [GSP-117](./117-rename-service-to-system-as-the-opposite-to-global.md): Rename `service metadata` to `system metadata`

## Background

Object may have different kind of metadata:

- last updated time
- content length
- content type
- storage class
- ...

Some of them are defined in standards, such as: `content-type`, `content-length` which defined in [RFC 2616](https://tools.ietf.org/html/rfc2616). Some of them are only used in special services, such as: `x-amz-storage-class` which only used in the AWS S3 service.

In order to unify object metadata behavior so that we can handle them in the same way, we need to clearly define and classify them.

[Proposal: Normalize Metadata](./6-normalize-metadata.md) is an attempt, it normalized metadata by [message-header](https://www.iana.org/assignments/message-headers/message-headers.xhtml), and handle all private metadata seperately.

- `Content-MD5` -> `content-md5`
- `x-aws-storage-class` / `x-qs-storage-class` -> `storage-class`

The problem is: How to implement cross-storage operations？

Giving a set object, how to migrate them to another storage services without changing its behavior? For example, permission on local file system, storage class on object storage class.

## Proposal

So, I propose to split all metadata into four groups which based on the definer:

- global metadata: defined via framework
    - id
    - link target
    - mode
    - multipart id
    - path
    - ...
- standard metadata: defined via existing RFCs
    - content md5
    - content type
    - etag
    - last-modified
    - ...
- service metadata: defined via service implementations
    - x-amz-storage-class
    - x-qs-storage-class
    - ...
- user metadata: defined via user input.

Both `global` and `standard` metadata will be included in object type.
`service` metadata will be stored in a struct that defined by service.
`user-defined` metadata will be stored in a hash map inside object type.

Use `golang` as example:

- `id` will be `id string`
- `content md5` will be `contentMd5 string`
- service metadata will be `serviceMetadata interface{}`
- user metadata will be `userMetadata map[string]string`

For service metadata,  we will introduce `Strong Typed Service Metadata`. We can generate following struct for every service:

```go
// For Service A
type ObjectMetadata struct {
    ServerSideEncryption                 string
    ServerSideEncryptionAwsKmsKeyID      string
    ServerSideEncryptionBucketKeyEnabled bool
    StorageClass                         string
    ...
}

// For Service B
type ObjectMetadata struct {
	ContentSha256 []byte
	...
}
```

And add following generated functions in service packages:

```go
// Only be used in service to set object metadata.
func setObjectMetadata(o *Object, om *ObjectMetadata) {}
// Only be used outside service package to get service metadata.
func GetObjectMetadata(o *Object) *ObjectMetadata {}
```

## Rationale

N/A

## Compatibility

All API call that used service metadata could be affected.

## Implementation

This proposal been implemented partially in go-storage with service metadata implemented as `map[string]string`, which made it hard to use.
