## storage

[![Build Status](https://travis-ci.com/Xuanwo/storage.svg?branch=master)](https://travis-ci.com/Xuanwo/storage)
[![GoDoc](https://godoc.org/github.com/Xuanwo/storage?status.svg)](https://godoc.org/github.com/Xuanwo/storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xuanwo/storage)](https://goreportcard.com/report/github.com/Xuanwo/storage)
[![codecov](https://codecov.io/gh/Xuanwo/storage/branch/master/graph/badge.svg)](https://codecov.io/gh/Xuanwo/storage)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE)

A unified storage layer for Golang.

### Goal

- Production ready
- High performance
- Vendor lock free

### Current Status

This lib is in heavy development, break changes could be introduced at any time. All public interface or functions expected to be stable at `v1.0.0`.

### Installation

Install will `go get`

```bash
go get -u github.com/Xuanwo/storage
```

Import

```go
import "github.com/Xuanwo/storage"
```

### Quickstart


```go
// Init a service.
srv, store, err := coreutils.Open("qingstor://static:test_access_key:test_secret_key@https:qingstor.com:443/test_bucket_name")
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// Use Storager API to maintain data.
ch := make(chan *types.Object, 1)
defer close(ch)

err := store.ListDir("prefix", pairs.WithFileFunc(func(*types.Object){
    ch <- o
}))
if err != nil {
    log.Printf("storager listdir failed: %v", err)
}
```

### Services

| Service | Description | Status |
| ------- | ----------- | ------ |
| fs | Local file system | stable |
| qingstor | [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/) | stable |
| s3 | [Amazon S3](https://aws.amazon.com/s3/) | alpha (-segments, -unittests) |
| oss | [Aliyun Object Storage](https://www.aliyun.com/product/oss) | planned |
| gcs | [Google Cloud Storage](https://cloud.google.com/storage/) | planned |
| kodo | [qiniu kodo](https://www.qiniu.com/products/kodo) | planned |
| uss | [UPYUN Storage Service](https://www.upyun.com/products/file-storage) | planned |
| cos | [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos) | planned |