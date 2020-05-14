# storage

[![构建状态](https://travis-ci.com/Xuanwo/storage.svg?branch=master)](https://travis-ci.com/Xuanwo/storage) [![GoDoc](https://godoc.org/github.com/Xuanwo/storage?status.svg)](https://godoc.org/github.com/Xuanwo/storage) [![Go Report Card](https://goreportcard.com/badge/github.com/Xuanwo/storage)](https://goreportcard.com/report/github.com/Xuanwo/storage) [![codecov](https://codecov.io/gh/Xuanwo/storage/branch/master/graph/badge.svg)](https://codecov.io/gh/Xuanwo/storage) [![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/15867a455afc4f24a763a5ed1011e05a)](https://app.codacy.com/manual/Xuanwo/storage?utm_source=github.com&utm_medium=referral&utm_content=Xuanwo/storage&utm_campaign=Badge_Grade_Settings) [![加入聊天](https://img.shields.io/badge/chat-online-blue?style=flat&logo=telegram)](https://t.me/storage_dev)

面向应用程序的 Golang 统一存储层。

## 目标

- 可用于生产环境
- 高性能
- 无供应商锁定

## 功能

### 服务级别

- 使用相同的 API 对存储服务进行基本操作
  - List: 列出所有 Storager
  - Get: 通过名称获取 Storager
  - Create: 创建一个 Storager
  - Delete: 删除一个 Storager

### 存储级别

- 使用相同的 API 对存储服务进行基本操作
  - Read: 读取文件内容
  - Write: 将内容写入文件
  - List: 列取一个目录或者前缀下的文件
  - Stat: 获取文件的元数据
  - Delete: 删除一个文件
  - Metadata: 获取存储服务的元数据
- 使用相同的 API 对存储服务进行高级操作
  - Copy: 复制一个文件
  - Move: 移动一个文件
  - Reach: 生成一个可公开访问的 URL
  - Statistical: 获取存储服务的统计数据
  - Segment: 对分块/分段的完整支持

### 文件级别

- 元数据
  - Content Length / Size: 通过 [RFC 2616](https://tools.ietf.org/html/rfc2616) 实现完整支持
  - Content MD5 / ETag: 通过 [proposal](docs/design/14-normalize-content-hash-check.md) 实现完整支持
  - Content Type: 通过 [RFC 2616](https://tools.ietf.org/html/rfc2616) 实现完整支持
  - Storage Class: 通过 [proposal](docs/design/8-normalize-metadata-storage-class.md) 实现完整支持

## 安装

通过 `go get` 安装

```bash
go get github.com/Xuanwo/storage
```

Import

```go
import "github.com/Xuanwo/storage"
```

## 快速开始

```go
// 初始化服务。
store, err := coreutils.OpenStorager("fs", pairs.WithWorkDir("/tmp"))
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// 使用 Storager API 来维护数据。
r, err := store。Read("path/to/file")
if err != nil {
    log.Printf("storager read: %v", err)
}
```

## 服务

- [zblob](./services/azblob/): [Azure Blob 存储](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](./services/cos/): [腾讯云对象存储](https://cloud.tencent.com/product/cos)
- [Dropbox](./services/dropbox/): [Dropbox](https://www.dropbox.com)
- [fs](./services/fs/): 本地文件系统
- [gcs](./services/gcs/): [Google 云存储](https://cloud.google.com/storage/)
- [kodo](./services/kodo/): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](./services/oss/): [Aliyun 对象存储](https://www.aliyun.com/product/oss)
- [qingstor](./services/qingstor/): [QingStor 对象存储](https://www.qingcloud.com/products/qingstor/)
- [s3](./services/s3/): [Amazon S3](https://aws.amazon.com/s3/)
- [uss](./services/uss/): [UPYUN 存储服务](https://www.upyun.com/products/file-storage)
