# storage

[![Build Status](https://travis-ci.com/Xuanwo/storage.svg?branch=master)](https://travis-ci.com/Xuanwo/storage)
[![GoDoc](https://godoc.org/github.com/Xuanwo/storage?status.svg)](https://godoc.org/github.com/Xuanwo/storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xuanwo/storage)](https://goreportcard.com/report/github.com/Xuanwo/storage)
[![codecov](https://codecov.io/gh/Xuanwo/storage/branch/master/graph/badge.svg)](https://codecov.io/gh/Xuanwo/storage)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)

An application-oriented unified storage layer for Golang.

## Goal

- Production ready
- High performance
- Vendor lock free

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

## Quickstart


```go
// Init a service.
srv, store, err := coreutils.Open("qingstor://hmac:test_access_key:test_secret_key@https:qingstor.com:443/test_bucket_name")
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
    log.Printf("storager listdir failed: %v", err)
}
```

## Services

| Service | Description | Status |
| ------- | ----------- | ------ |
| [azblob](#azblob) | [Azure Blob storage](https://docs.microsoft.com/en-us/azure/storage/blobs/) | alpha (-segments, -unittests) |
| [cos](#cos) | [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos) | alpha (-segments, -unittests) |
| [dropbox](#dropbox) | [Dropbox](https://www.dropbox.com) | alpha (-unittests) |
| [fs](#fs) | Local file system | stable (-segments)|
| [gcs](#gcs) | [Google Cloud Storage](https://cloud.google.com/storage/) | alpha (-segments, -unittests) |
| [kodo](#kodo) | [qiniu kodo](https://www.qiniu.com/products/kodo) | alpha (-segments, -unittests) |
| [oss](#oss) | [Aliyun Object Storage](https://www.aliyun.com/product/oss) | alpha (-segments, -unittests) |
| [qingstor](#qingstor) | [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/) | stable |
| [s3](#s3) | [Amazon S3](https://aws.amazon.com/s3/) | alpha (-segments, -unittests) |
| [uss](#uss) | [UPYUN Storage Service](https://www.upyun.com/products/file-storage) | planned |

### azblob

`azblob://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

### cos

`cos://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

### dropbox

`dropbox://apikey:<api_key>/path/to/dir`

### fs

`fs:///path/to/dir`

### gcs

`gcs://apikey:<api_key>/<bucket_name>/<prefix>?project=<project_id>`

### kodo

`kodo://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`

### oss

`oss://hmac:<access_key>:<secret_key>@<protocol>:<host>:<port>/<bucket_name>/<prefix>`

### qingstor

`qingstor://hmac:<access_key>:<secret_key>@<protocol>:<host>:<port>/<bucket_name>/<prefix>`

### s3

`s3://hmac:<access_key>:<secret_key>/<bucket_name>/<prefix>`
