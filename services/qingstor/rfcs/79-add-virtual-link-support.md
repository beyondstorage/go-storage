- Author: abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-08-16
- RFC PR: [beyondstorage/go-service-qingstor#79](https://github.com/beyondstorage/go-service-qingstor/pull/79)
- Tracking Issue: [beyondstorage/go-service-qingstor#64](https://github.com/beyondstorage/go-service-qingstor/issues/64)

# RFC-79: Add Virtual Link Support

## Background

As you can see from the [official documentation for qingstor](https://docs.qingcloud.com/qingstor/), it does not support symbolic links itself. However, qingstor supports [user-defined object metadata](https://docs.qingcloud.com/qingstor/api/common/metadata), and we can use this to simulate the implementation of symlink.

## Proposal

I propose to use user-defined object metadata to support virtual link.

```go
input := &service.PutObjectInput{
    XQSMetadata: &map[string]string{
        "x-qs-meta-bs-link-target": rt,
    },
}
```

- `PutObjectInput` is used to store the fields we need when calling `PutObjectWithContext` API to upload an object.
- `XQSMetadata` is a map that stores user-defined metadata.
  - `"x-qs-meta-bs-link-target"` is the name of user-defined metadata, the middle `bs` is used to avoid conflicts.
  - `rt` is the symlink target, it is an absolute path.

## Rationale

### User-defined metadata

As the [Object Metadata](https://docs.qingcloud.com/qingstor/api/common/metadata) says,there are two types of object metadata that can be changed by the user: standard HTTP headers and user-defined metadata. We can define our own metadata to store the fields we need when uploading an object. Note that the name of the user-defined metadata must start with `x-qs-meta`.

### Drawbacks

As qingstor itself does not support symlink, we can only simulate it. And the object created is not really a symlink object. When we call `stat`, we can only tell if it is a symlink by using user-defined metadata,

```go
if v, ok := metadata["x-qs-meta-bs-link-target"]; ok {
	// The path is a symlink object 
}
```

Calling `HeadObject` in `list` will increase the execution cost of `list` which we cannot afford. So we will relax the qingstor condition. We will not support getting the exact symlink object type in `list` when user has virtual link enabled, if the user wants to get the exact object mode they need to call `stat`.

## Compatibility

N/A

## Implementation

- Implement `virtual_link` in go-service-qingstor
- Support `stat`
- Setup linker tests

