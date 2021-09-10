- Author:  Abyss-w <mad.hatter@foxmail.com>
- Start Date: 2021-9-10
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-0: Write Empty File Behavior

## Background

In our definition, the `Write` API will upload a file to the target object. We do not do anything about uploading empty files. However, some services do not support uploading empty files. As [proposed](https://forum.beyondstorage.io/t/topic/204) in our community, for s3, when we upload an empty file, the service will hang permanently. Here is an example of s3:

```go
_, err = fs.s.Write(path,nil, 0)
if err != nil {
    fs.logger.Error("write", zap.String("path", path), zap.Error(err))
    return nil, nil, err
}
```

If we want to upload an empty file, we need to do so:

```go
_, err = fs.s.Write(path, bytes.NewReader([]byte{}), 0)
if err != nil {
    fs.logger.Error("write", zap.String("path", path), zap.Error(err))
    return nil, nil, err
}
```

This is not convenient for the users.

## Proposal

I propose to allow the user to pass in an empty file when calling `Write`. 

- For services that do not support uploading empty files, we should check if the `reader` is `nil`. If it is `nil`, we need to create a `reader` before calling the API. Like s3, kodo, etc.
- For services that support uploading empty files, we don't need to check if the `reader` is `nil`, we can call the API directly. Like oss, etc.

## Rationale

N/A

## Compatibility

This behavior will affect all services that do not support uploading empty files.

## Implementation

- Update definitions to reflect changes
- Update integrations tests to make sure all service passed

