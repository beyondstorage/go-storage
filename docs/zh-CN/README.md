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
  - 列表：列出所有正在使用的存储器
  - 获取：通过名称获取存储器
  - 创建：创建一个存储器
  - 删除：删除存储器

### 存储器级别

- 使用相同的 API 的所有存储服务的基本操作
  - 阅读：读取文件内容
  - 写入：将内容写入文件
  - 列表：目录或前缀下的文件列表
  - 状态：获取文件的元数据
  - 删除：删除一个文件
  - 元数据：获取存储服务的元数据
- 具有相同API的已实现存储服务的高级操作
  - 复制：复制文件
  - 移动：移动一个文件
  - 达到：生成一个公共可访问的 url
  - 统计：获取存储服务的统计
  - 部分：对部分、卡片、多部分的充分支持

### 文件级别

- 元数据
  - 内容长度 / 大小：通过 [RFC 2616](https://tools.ietf.org/html/rfc2616) 完全支持
  - 内容MD5 / ETag：通过 [提议完全支持](docs/design/14-normalize-content-hash-check.md)
  - 内容类型：通过 [RFC 2616 获得完全支持](https://tools.ietf.org/html/rfc2616)
  - 存储类别：通过 [提议完全支持](docs/design/8-normalize-metadata-storage-class.md)

## 安装

安装 `将获得`

```bash
去获取 github.com/Xuanwo/存储
```

导入

```go
导入 "github.com/Xuanwo/storage"
```

## 快速开始

```go
// 初始化服务。
store, err := coreutils.OpenStorager("fs", 配对。WewWorkDir("/tmp"))
如果是err != nil 。
    log。Fatalf("服务失败: %v", 错误)
}

// 使用 Storager API 来维护数据。
r, err := store。读取("path/to/file")
如果是err != nil 。
    log。Printf("存储器读取： %v", 错误)
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
