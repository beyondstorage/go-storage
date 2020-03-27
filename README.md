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
- Vendor lock free

## Features

### Servicer Level

- Basic operations across implemented storage services with the same API
  - List: list all Storager in service
  - Get: get a Storager via name
  - Create: create a Storager
  - Delete: delete a Storager

### Storager Level

- Basic operations across all storage services with the same API
  - Read: read file content
  - Write: write content into file
  - List: list files under a dir or prefix
  - Stat: get file's metadata
  - Delete: delete a file
  - Metadata: get storage service's metadata
- Advanced operations across implemented storage services with the same API
  - Copy: copy a file
  - Move: move a file
  - Reach: generate a public accessible url
  - Statistical: get storage service's statistics
  - Segment: Full support for Segment, aka, Multipart

### File Level

- Metadata
  - Content Length / Size: Full support via [RFC 2616](https://tools.ietf.org/html/rfc2616)
  - Content MD5 / ETag: Full support via [proposal](docs/design/14-normalize-content-hash-check.md)
  - Content Type: Full support via [RFC 2616](https://tools.ietf.org/html/rfc2616)
  - Storage Class: Full support via [proposal](docs/design/8-normalize-metadata-storage-class.md)  

## Current Status

This lib is in heavy development, break changes could be introduced at any time. All public interface or functions expected to be stable at `v1.0.0`.

## Installation

Install will `go get`

```bash
go get -u github.com/Xuanwo/storage
```

Import

```go
import "github.com/Xuanwo/storage"
```

## Quick Start

```go
// Init a service.
srv, store, err := coreutils.Open("fs", pairs.WithWorkDir("/tmp"))
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// Use Storager API to maintain data.
ch := make(chan *types.Object, 1)
defer close(ch)

err := store.List("prefix", pairs.WithFileFunc(func(*types.Object){
    ch <- o
}))
if err != nil {
    log.Printf("storager list failed: %v", err)
}
```

## Services

- [azblob](docs/services/azblob.md): [Azure Blob storage](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](docs/services/cos.md): [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos)
- [dropbox](docs/services/dropbox.md): [Dropbox](https://www.dropbox.com)
- [fs](docs/services/fs.md): Local file system
- [gcs](docs/services/gcs.md): [Google Cloud Storage](https://cloud.google.com/storage/)
- [kodo](docs/services/kodo.md): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](docs/services/oss.md): [Aliyun Object Storage](https://www.aliyun.com/product/oss)
- [qingstor](docs/services/qingstor.md): [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/)
- [s3](docs/services/s3.md): [Amazon S3](https://aws.amazon.com/s3/)
- [uss](docs/services/uss.md): [UPYUN Storage Service](https://www.upyun.com/products/file-storage)
