# go-storage

## [网站](https://beyondstorage.io) | [文档](https://beyondstorage.io/docs/go-storage/index) | [社区](https://beyondstorage.io/community)

[![Build Status](https://github.com/beyondstorage/go-storage/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-storage/actions?query=workflow%3A%22Unit+Test%22)
[![Go dev](https://pkg.go.dev/badge/github.com/beyondstorage/go-storage/v4)](https://pkg.go.dev/github.com/beyondstorage/go-storage/v4)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/beyondstorage/go-storage/blob/master/LICENSE)
[![go storage dev](https://img.shields.io/matrix/beyondstorage@go-storage:matrix.org.svg?label=go-storage&logo=matrix)](https://matrix.to/#/#beyondstorage@go-storage:matrix.org)

一个**厂商中立**的存储库。

## 愿景

**只需一次编写，即可在任意存储服务中运行。**

## 目标

- 厂商中立
- 随时投产
- 极致性能

## 示例

```go
package main

import (
    "log"

    "github.com/beyondstorage/go-storage/v4/services"
    "github.com/beyondstorage/go-storage/v4/types"

    // 添加 fs 支持
    _ "github.com/beyondstorage/go-service-fs/v3"
    // 添加 s3 支持
    _ "github.com/beyondstorage/go-service-s3/v2"
    // 添加 gcs 支持
    _ "github.com/beyondstorage/go-service-gcs/v2"
    // 添加 azblob 支持
    _ "github.com/beyondstorage/go-service-azblob/v2"
    // 更多支持，可在 BeyondStorage 下获取
    _ "github.com/beyondstorage/go-service-xxx" 
)

func main() {
    // 连接字符串并初始化存储器
    store, err := services.NewStoragerFromString("s3://bucket_name/path/to/workdir")
    if err != nil {
        log.Fatalf("service init failed: %v", err)
    }

    // 将 io.Reader 中的数据写入 hello.txt
    n, err := store.Write("hello.txt", r, length)

    // 从 hello.txt 中读取数据到 io.Writer
    n, err := store.Read("hello.txt", w)

    // 检查 hello.txt 是否存在并获取其元数据
    o, err := store.Stat("hello.txt")

    // 使用对象的函数获取元数据
    length, ok := o.GetContentLength()
    
    // 列表将在路径下创建一个对象的迭代器
    it, err := store.List("path")
    
    for {
    	// 使用迭代器检索下一个对象，直到迭代完成
    	o, err := it.Next()
    	if errors.Is(err, types.IteraoorDone) {
    		break
        }
    }

    // 删除 hello.txt
    err = store.Delete("hello.txt")
}
```

更多示例可以在 [go-storage-example](https://github.com/beyondstorage/go-storage-example) 找到。

## 特点

### 支持多种本地服务

目前已经有 **14** 个稳定的服务通过了所有的 [集成测试](https://github.com/beyondstorage/go-integration-test)。

- [azblob](https://github.com/beyondstorage/go-service-azblob/): [Azure Blob storage](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](https://github.com/beyondstorage/go-service-cos/): [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos)
- [dropbox](https://github.com/beyondstorage/go-service-dropbox/): [Dropbox](https://www.dropbox.com)
- [fs](https://github.com/beyondstorage/go-service-fs/): Local file system
- [gcs](https://github.com/beyondstorage/go-service-gcs/): [Google Cloud Storage](https://cloud.google.com/storage/)
- [kodo](https://github.com/beyondstorage/go-service-kodo/): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](https://github.com/beyondstorage/go-service-oss/): [Aliyun Object Storage](https://www.aliyun.com/product/oss)
- [qingstor](https://github.com/beyondstorage/go-service-qingstor/): [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/)
- [s3](https://github.com/beyondstorage/go-service-s3/): [Amazon S3](https://aws.amazon.com/s3/)
- [ftp](https://github.com/beyondstorage/go-service-ftp/): FTP
- [gdrive](https://github.com/beyondstorage/go-service-gdrive): [Google Drive](https://www.google.com/drive/)
- [ipfs](https://github.com/beyondstorage/go-service-ipfs): [InterPlanetary File System](https://ipfs.io)
- [memory](https://github.com/beyondstorage/go-service-memory): data that only in memory
- [minio](https://github.com/beyondstorage/go-service-minio): [MinIO](https://min.io)

另有 **3** 个公测版本的服务已实现了所需功能，但还没有通过 [集成测试](https://github.com/beyondstorage/go-integration-test)。

- [uss](https://github.com/beyondstorage/go-service-uss/): [UPYUN Storage Service](https://www.upyun.com/products/file-storage)
- [hdfs](https://github.com/beyondstorage/go-service-hdfs): [Hadoop Distributed File System](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html#Introduction)
- [tar](https://github.com/beyondstorage/go-service-tar): tar files

最后还有 **4** 个处于内测阶段的服务仍在开发中。

- [onedrive](https://github.com/beyondstorage/go-service-onedrive): [Microsoft OneDrive](https://www.microsoft.com/en-ww/microsoft-365/onedrive/online-cloud-storage)
- [storj](https://github.com/beyondstorage/go-service-storj): [StorJ](https://www.storj.io/)
- [webdav](https://github.com/beyondstorage/go-service-webdav): [WebDAV](http://www.webdav.org/)
- [zip](https://github.com/beyondstorage/go-service-zip): zip files

更多关于服务的想法可以在 [Service Integration Tracking](https://github.com/beyondstorage/go-storage/issues/536) 找到。

### 完整且易扩展的接口

基本操作

- 元数据: 获取 `存储器` 元数据 
```go
meta := store.Metadata()
_ := meta.GetWorkDir() // 获取对象的工作目录
_, ok := meta.GetWriteSizeMaximum() // 获取写操作的最大尺寸
```
- 读取: 读取 `对象` 的内容
```go
// 在偏移量 1024 处读取 2048 字节到 io.Writer
n, err := store.Read("path", w, pairs.WithOffset(1024), pairs.WithSize(2048))
```
- 写入: 将内容写入 `对象` 中
```go
// 从 io.Reader 写入 2048 字节
n, err := store.Write("path", r, 2048)
```
- 统计: 获取 `对象` 元数据并检查是否存在
```go
o, err := store.Stat("path")
if errors.Is(err, services.ErrObjectNotExist) {
	// 对象不存在
}
length, ok := o.GetContentLength() // 获取对象的内容长度
```
- 删除: 删除一个 `对象`
```go
err := store.Delete("path") // 删除对象 "路径"
```
- 列表: 列出给定前缀或目录中的 `对象` 
```go
it, err := store.List("path")
for {
	o, err := it.Next()
	if err != nil && errors.Is(err, types.IteratorDone) {
        // 列表结束
    }
    length, ok := o.GetContentLength() // 获取对象的内容长度
}
```

扩展操作

- 拷贝: 复制一个 `对象` 到存储库内
```go
err := store.(Copier).Copy(src, dst) // 从 src 复制一个对象到 dst
```
- 移动: 移动一个 `对象` 到存储库内
```go
err := store.(Mover).Move(src, dst) // 从 src 移动一个对象到 dst
```
- 链接: 为 `对象` 生成一个可访问的公共 url
```go
url, err := store.(Reacher).Reach("path") // 生成一个对象的 url
```
- 目录: `对象` 目录
```go
o, err := store.(Direr).CreateDir("path") // 创建一个对象目录
```

大文件操作

- 分段: 允许进行分段上传
```go
ms := store.(Multiparter)

// 创建一个分段对象
o, err := ms.CreateMultipart("path")
// 将 io.reader 中的 1024 字节分段写入索引 1
n, part, err := ms.WriteMultipart(o, r, 1024, 1)
// 完成分段对象创建
err := ms.CompleteMultipart(o, []*Part{part})
```
- 追加: 允许追加到一个对象上
```go
as := store.(Appender)

// 创建一个可追加的对象
o, err := as.CreateAppend("path")
// 从 io.Reader 写入 1024 字节
n, err := as.WriteAppend(o, r, 1024)
// 提交一个待追加的对象
err = as.CommitAppend(o)
```
- 块: 允许将一个对象与块 id 进行组合
```go
bs := store.(Blocker)

// 创建一个块对象
o, err := bs.CreateBlock("path")
// 将 io.reader 中的 1024 字节写入 id 为 ”id-abc“ 的块
n, err := bs.WriteBlock(o, r, 1024, "id-abc")
// 通过块 id 组合区块
err := bs.CombineBlock(o, []string{"id-abc"})
```
- 页面：允许进行随机写入
```go
ps := store.(Pager)

// 创建一个页面对象
o, err := ps.CreatePage("path")
// Write 1024 bytes from io.Reader at offset 2048
n, err := ps.WritePage(o, r, 1024, 2048)
```

### 综合元数据

全局对象元数据

- `id`: 服务中的唯一键
- `name`: 服务工作目录的相对路径
- `mode`: 对象的模式可以由以下几种进行组合：`read`, `dir`, `part` and [more](https://github.com/beyondstorage/go-storage/blob/master/types/object.go#L11) 
- `etag`: 实体标签，定义于 [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.19) 
- `content-length`: 对象的内容大小
- `content-md5`: [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.15) 中定义的 Md5 简介
- `content-type`: [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.17) 中定义的媒体类型
- `last-modified`: 对象的最后更新时间

系统对象元数据

服务系统对象元数据，如 `storage-class` 等。

```go
o, err := store.Stat("path")

// 通过 go-service-s3 提供的 API 获取服务系统元数据
om := s3.GetObjectSystemMetadata(o)
_ = om.StorageClass // 此对象的存储类型
_ = om.ServerSideEncryptionCustomerAlgorithm // 此对象的 sse 算法
```

### 强类型的接口

自我维护的代码生成器 [定义](https://github.com/beyondstorage/go-storage/tree/master/cmd/definitions) 有助于生成我们所有的 API、配对和元数据。

生成的 pairs 可用作 API 的可选参数。

```go
func WithContentMd5(v string) Pair {
    return Pair{
        Key:   "content_md5",
        Value: v,
    }
}
```

生成的对象元数据可用于从对象中获取内容 md5。

```go
func (o *Object) GetContentMd5() (string, bool) {
    o.stat()
    
    if o.bit&objectIndexContentMd5 != 0 {
        return o.contentMd5, true
    }
    
    return "", false
}
```

### 服务器端加密

服务器端加密支持在 system pair 和 system metadata 中使用, 并且我们可以通过 [Default Pairs](https://beyondstorage.io/docs/go-storage/pairs/index#default-pairs) 来简化工作。

```go

func NewS3SseC(key []byte) (types.Storager, error) {
    defaultPairs := s3.DefaultStoragePairs{
        Write: []types.Pair{
            // 需要, 必须为 AES256
            s3.WithServerSideEncryptionCustomerAlgorithm(s3.ServerSideEncryptionAes256),
            // 你的 AES-256 密钥, 需要是一个 32 字节的二进制值
            s3.WithServerSideEncryptionCustomerKey(key),
        },
        // 现在你必须提供客户密钥才能读取加密数据
        Read: []types.Pair{
            // 需要, 必须为 AES256
            s3.WithServerSideEncryptionCustomerAlgorithm(s3.ServerSideEncryptionAes256),
            // Required, your AES-256 key, a 32-byte binary value
            s3.WithServerSideEncryptionCustomerKey(key),
        }}
    
    return s3.NewStorager(..., s3.WithDefaultStoragePairs(defaultPairs))
}
```

## 赞助商

<a href="https://vercel.com?utm_source=beyondstorage&utm_campaign=oss">
    <img src="./docs/images/vercel.svg">
</a>
