# storage

[![Build Status](https://travis-ci.com/Xuanwo/storage.svg?branch=master)](https://travis-ci.com/Xuanwo/storage)
[![GoDoc](https://godoc.org/github.com/Xuanwo/storage?status.svg)](https://godoc.org/github.com/Xuanwo/storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xuanwo/storage)](https://goreportcard.com/report/github.com/Xuanwo/storage)
[![codecov](https://codecov.io/gh/Xuanwo/storage/branch/master/graph/badge.svg)](https://codecov.io/gh/Xuanwo/storage)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/15867a455afc4f24a763a5ed1011e05a)](https://app.codacy.com/manual/Xuanwo/storage?utm_source=github.com&utm_medium=referral&utm_content=Xuanwo/storage&utm_campaign=Badge_Grade_Settings)
[![Join the chat](https://img.shields.io/badge/chat-online-blue?style=flat&logo=telegram)](https://t.me/storage_dev)

An application-oriented unified storage layer for Golang.

## Goal

- Production ready
- High performance
- Vendor agnostic

## Features

### Widely services support

- [azblob](./services/azblob/): [Azure Blob storage](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](./services/cos/): [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos)
- [dropbox](./services/dropbox/): [Dropbox](https://www.dropbox.com)
- [fs](./services/fs/): Local file system
- [gcs](./services/gcs/): [Google Cloud Storage](https://cloud.google.com/storage/)
- [kodo](./services/kodo/): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](./services/oss/): [Aliyun Object Storage](https://www.aliyun.com/product/oss)
- [qingstor](./services/qingstor/): [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/)
- [s3](./services/s3/): [Amazon S3](https://aws.amazon.com/s3/)
- [uss](./services/uss/): [UPYUN Storage Service](https://www.upyun.com/products/file-storage)

### Servicer operation support

- List: list all Storager in service
- Get: get a Storager via name
- Create: create a Storager
- Delete: delete a Storager

### Storager operation support

Basic operations

- Metadata: get storager's metadata
- Read: read file content
- Write: write content into file
- Stat: get file's metadata
- Delete: delete a file or directory

Extended operations

- Copy: copy a file inside storager
- Move: move a file inside storager
- Reach: generate a public accessible url
- Statistical: get storage service's statistics

Multiple list style support

- ListDir: list files and directories under a directory
- ListPrefix: list files under a prefix

Segment/Multipart support

- ListPrefixSegment: list segments under a prefix
- InitIndexSegment: initiate an index type segment
- WriteIndexSegment: write content into an index type segment
- CompleteSegment: complete a segment to create a file
- AbortSegment: abort a segment

### File metadata support

Required metadata

- `id`: unique key in service
- `name`: relative path towards service's work dir
- `size`: size of this object
- `updated_at`: last update time of this object

Optional metadata

- `content-md5`: md5 digest as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.15)
- `content-type`: media type as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.17)
- `etag`: entity tag as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.19)
- `storage-class`: object's storage class as defined in [storage proposal](./docs/design/8-normalize-metadata-storage-class.md)

## Quick Start

```go
import (
    "log"

    "github.com/Xuanwo/storage"
    "github.com/Xuanwo/storage/coreutils"
    "github.com/Xuanwo/storage/types/pairs"
)

// Init a service.
store, err := coreutils.OpenStorager("fs", pairs.WithWorkDir("/tmp"))
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// Use Storager API to maintain data.
r, err := store.Read("path/to/file")
if err != nil {
    log.Printf("storager read: %v", err)
}
```
