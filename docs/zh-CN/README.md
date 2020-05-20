# storage

[![构建状态](https://travis-ci.com/Xuanwo/storage.svg?branch=master)](https://travis-ci.com/Xuanwo/storage) [![GoDoc](https://godoc.org/github.com/Xuanwo/storage?status.svg)](https://godoc.org/github.com/Xuanwo/storage) [![Go Report Card](https://goreportcard.com/badge/github.com/Xuanwo/storage)](https://goreportcard.com/report/github.com/Xuanwo/storage) [![codecov](https://codecov.io/gh/Xuanwo/storage/branch/master/graph/badge.svg)](https://codecov.io/gh/Xuanwo/storage) [![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/Xuanwo/storage/blob/master/LICENSE) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/15867a455afc4f24a763a5ed1011e05a)](https://app.codacy.com/manual/Xuanwo/storage?utm_source=github.com&utm_medium=referral&utm_content=Xuanwo/storage&utm_campaign=Badge_Grade_Settings) [![加入聊天](https://img.shields.io/badge/chat-online-blue?style=flat&logo=telegram)](https://t.me/storage_dev)

面向应用程序的 Golang 统一存储层。

## 目标

- 可用于生产环境
- 高性能
- 供应商无关

## 功能

### 广泛的服务支持

- [zblob](./services/azblob/): [Azure Blob 存储](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](./services/cos/): [腾讯云对象存储](https://cloud.tencent.com/product/cos)
- [Dropbox](./services/dropbox/): [Dropbox](https://www.dropbox.com)
- [fs](./services/fs/)：本地文件系统
- [gcs](./services/gcs/): [Google 云存储](https://cloud.google.com/storage/)
- [kodo](./services/kodo/): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](./services/oss/): [Aliyun 对象存储](https://www.aliyun.com/product/oss)
- [qingstor](./services/qingstor/): [QingStor 对象存储](https://www.qingcloud.com/products/qingstor/)
- [s3](./services/s3/): [Amazon S3](https://aws.amazon.com/s3/)
- [uss](./services/uss/): [UPYUN 存储服务](https://www.upyun.com/products/file-storage)

### 服务级别操作支持

- List: 列出所有 Storager
- Get: 通过名称获取 Storager
- Create: 创建一个 Storager
- Delete: 删除一个 Storager

### 存储级别操作支持

基本操作

- Metadata: 获取存储服务的元数据
- Read: 读取文件内容
- Write: 将内容写入文件
- Stat: 获取文件的元数据
- Delete：删除文件或目录

扩展操作

- Copy: 在存储器中复制一个文件
- Move：在存储器中移动一个文件
- Reach: 生成一个可公开访问的 URL
- Statistical：获取存储服务的统计

多种列取风格支持

- ListDir：在目录下列出文件和目录
- ListPrefix：列出前缀下的文件

分段上传支持

- ListPrefixSegment: 列出前缀下面的段
- InitIndexSegment：创建一个基于索引的分段
- WriteIndexSegment: 将内容写入一个基于索引的分段
- CompleteSegment: 完成一个分段以创建一个文件
- AbortSegment: 中止一个分段

### 文件元数据支持

必须的元数据

- `id`：服务中全局唯一的键
- `name`：到当前工作路径的相对路径
- `size`：此对象的大小
- `updated_at`：此对象的最后更新时间

可选的元数据

- `content-md5`：[rfc2616](https://tools.ietf.org/html/rfc2616#section-14.15)中定义的 md5 摘要
- `content-type`：[rfc2616](https://tools.ietf.org/html/rfc2616#section-14.17)中定义的媒体类型
- `etag`：在[rfc2616](https://tools.ietf.org/html/rfc2616#section-14.19)中定义的实体标签
- `storage-class`：在[草案](./design/8-normalize-metadata-storage-class.md)中定义的对象存储级别

## 快速开始

```go
import (
    "log"

    "github.com/Xuanwo/storage"
    "github.com/Xuanwo/storage/coreutils"
    "github.com/Xuanwo/storage/types/pairs"
)

// 初始化一个 Storager
store, err := coreutils.OpenStorager("fs", pairs.WithWorkDir("/tmp"))
if err != nil {
    log.Fatalf("service init failed: %v", err)
}

// 使用 Storager API 来维护数据
r, err := store.Read("path/to/file")
if err != nil {
    log.Printf("storager read: %v", err)
}
```
