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
srv := qingstor.New()
err := srv.Init(
    pairs.WithAccessKey("test_access_key"),
    pairs.WithSecretKey("test_secret_key"),
)
if err != nil {
    log.Printf("service init failed: %v", err)
}

// Get a storage instance from an initiated service.
store, err := srv.Get("test_bucket_name")
if err != nil {
    log.Printf("service get bucket failed: %v", err)
}

// Init a storage.
err := store.Init(pairs.WithWorkDir("/prefix"))
if err != nil {
    log.Printf("storager init failed: %v", err)
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