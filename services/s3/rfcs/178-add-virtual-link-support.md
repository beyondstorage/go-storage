- Author: abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-08-11
- RFC PR: [beyondstorage/go-service-s3#178](https://github.com/beyondstorage/go-service-s3/pull/178)
- Tracking Issue: [beyondstorage/go-service-s3#144](https://github.com/beyondstorage/go-service-s3/issues/144)

# RFC-178: Add Virtual Link Support

## Background

Like the one presented in [GSP-86 Add Create Link Operation](https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/86-add-create-link-operation.md), s3 has no native support for symlink. But we can use [user-defined object metadata](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/UsingMetadata.html#UserMetadata) to simulate it.

## Proposal

I propose to use user-defined object metadata to implement virtual_link feature to support symlink in s3.

```go
input := &s3.PutObjectInput{
    Metadata: map[string]*string{
			"x-amz-meta-bs-link-target": &rt,
		},
}
```

- `PutObjectInput` in s3 is used to store the fields we need when calling `PutObjectWithContext` API to upload an object.
- `Metadata` is a map that stores user-defined metadata.
  - `"x-amz-meta-bs-link-target"` is the name of user-defined metadata, the middle `bs` is used to avoid conflicts.
  - `rt` is the value of `"x-amz-meta-bs-link-target"`, which is the target of the symlink, it is an absolute path.

## Rationale

### User-defined metadata

As stated in the document [user-defined object metadata](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/userguide/UsingMetadata.html#UserMetadata), we can define our own metadata to store the information we need when uploading an object. One thing to note is that user-defined metadata names must start with `"x-amz-meta-"`.

### Drawbacks

As s3 itself does not support symlink, we can only simulate it. And the object created is not really a symlink object. When we call `stat`, we can only tell if it is a symlink by using user-defined object metadata.

```go
if v, ok := metadata["x-amz-meta-bs-link-target"]; ok {
	// The path is a symlink object. 
}
```

Calling `HeadObject` in `list` will increase the execution cost of `list` which we cannot afford. So we will relax the s3 condition. We will not support getting the exact symlink object type in `list` when the user has virtual linking enabled, if the user wants to get the exact object mode they need to call `stat`.

## Compatibility

N/A

## Implementation

- Implement `virtual_link` in go-service-s3
- Support `stat`
- Setup linker tests