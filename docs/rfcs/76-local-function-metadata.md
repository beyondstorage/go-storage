- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2021-05-21
- RFC PR: [beyondstorage/specs#76](https://github.com/beyondstorage/specs/issues/76)
- Tracking Issue: N/A

# GSP-76: Local Function Metadata

## Background

[Metadata] is a function to retrieve Storage's Metadata:

```go
Metadata(pairs ...Pair) (meta *StorageMeta, err error)
MetadataWithContext(ctx context.Context, pairs ...Pair) (meta *StorageMeta)
```

In our current implementations, Metadata looks like following:

```go
func (s *Storage) metadata(ctx context.Context, opt pairStorageMetadata) (meta *StorageMeta, err error) {
	meta = NewStorageMeta()
	meta.Name = *s.properties.BucketName
	meta.WorkDir = s.workDir
	meta.SetLocation(*s.properties.Zone)
	return meta, nil
}
```

We will not send API/RPC call in this function, a.k.a., this function never returns errors.

But our user still need to check them:

```go
meta, err := s.Metadata()
if err != nil {
	return err
}
```

## Proposal

So I propose to make `Metadata` a local function:

```go
Metadata(pairs ...Pair) (meta *StorageMeta)
```

This function will not return error, and no need for `Context`, just like our `Create` API.

## Rationale

N/A

## Compatibility

This change will break all services, so we expect to be released in next major version.

## Implementation

- Update specs
- Update go-storage
- Update go-serivce-*
