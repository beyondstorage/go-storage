# storage

[![Build Status](https://github.com/aos-dev/go-storage/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/aos-dev/go-storage/actions?query=workflow%3A%22Unit+Test%22)
[![Go dev](https://pkg.go.dev/badge/github.com/aos-dev/go-storage/v3)](https://pkg.go.dev/github.com/aos-dev/go-storage/v3)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/aos-dev/go-storage/v3/blob/master/LICENSE)
[![Join the chat](https://img.shields.io/badge/chat-online-blue?style=flat&logo=zulip)](https://aos-dev.zulipchat.com/join/c3sqj64sp53tlau7oojg3yll/)

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

- Metadata: get storager metadata
- Read: read file content
- Write: write content into file
- Stat: get file's metadata
- Delete: delete a file or directory
- List: list file in prefix or dir styles

Extended operations

- Copy: copy a file inside storager
- Move: move a file inside storager
- Reach: generate a public accessible url

Multi object modes support

- Multipart: allow doing multipart uploads
- Append: allow appending to an object
- Block: allow combining an object with block ids.
- Page: allow doing random writes

### Object metadata support

Common metadata

- `id`: unique key in service
- `name`: relative path towards service's work dir
- `type`: object type cloud be `file`, `dir`, `link` or `unknown`

Optional metadata

- `size`: object's content size.
- `updated-at`: object's last updated time.
- `content-md5`: md5 digest as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.15)
- `content-type`: media type as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.17)
- `etag`: entity tag as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.19)
- `storage-class`: object's storage class as defined
  in [storage proposal](./docs/design/8-normalize-metadata-storage-class.md)

## Quick Start

```go
import (
    "log"

    "github.com/aos-dev/go-storage/v3"
    "github.com/aos-dev/go-storage/v3/pairs"
    "github.com/aos-dev/go-services-fs"
)

// Init a service.
store, err := fs.NewStorager(pairs.WithWorkDir("/tmp"))
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// Use Storager API to maintain data.
var buf bytes.Buffer

n, err := store.Read("path/to/file", &buf)
if err != nil {
    log.Printf("storager read: %v", err)
}
```

## Sponsor

<a href="https://vercel.com?utm_source=aos-dev&utm_campaign=oss">
    <img src="./docs/images/vercel.svg">
</a>