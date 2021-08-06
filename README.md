# go-storage

## [Website](https://beyondstorage.io) | [Documentation](https://beyondstorage.io/docs/go-storage/index) | [Community](https://beyondstorage.io/community)

[![Build Status](https://github.com/beyondstorage/go-storage/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-storage/actions?query=workflow%3A%22Unit+Test%22)
[![Go dev](https://pkg.go.dev/badge/github.com/beyondstorage/go-storage/v4)](https://pkg.go.dev/github.com/beyondstorage/go-storage/v4)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/beyondstorage/go-storage/blob/master/LICENSE)
[![go storage dev](https://img.shields.io/matrix/beyondstorage@go-storage:matrix.org.svg?label=go-storage&logo=matrix)](https://matrix.to/#/#beyondstorage@go-storage:matrix.org)

A **vendor-neutral** storage library for Golang.

## Vision

**Write once, run on every storage service.**

## Goal

- Vendor agnostic
- Production ready
- High performance

## Examples

```go
package main

import (
    "log"

    "github.com/beyondstorage/go-storage/v4/services"
    "github.com/beyondstorage/go-storage/v4/types"

    // Add fs support
    _ "github.com/beyondstorage/go-service-fs/v3"
    // Add s3 support
    _ "github.com/beyondstorage/go-service-s3/v2"
    // Add gcs support
    _ "github.com/beyondstorage/go-service-gcs/v2"
    // Add azblob support
    _ "github.com/beyondstorage/go-service-azblob/v2"
    // More support could be found under BeyondStorage.
    _ "github.com/beyondstorage/go-service-xxx" 
)

func main() {
    // Init a Storager from connection string. 
    store, err := services.NewStoragerFromString("s3://bucket_name/path/to/workdir")
    if err != nil {
        log.Fatalf("service init failed: %v", err)
    }

    // Write data from io.Reader into hello.txt
    n, err := store.Write("hello.txt", r, length)

    // Read data from hello.txt to io.Writer
    n, err := store.Read("hello.txt", w)

    // Stat hello.txt to check existence or get its metadata
    o, err := store.Stat("hello.txt")

    // Use object's functions to get metadata
    length, ok := o.GetContentLength()
    
    // List will create an iterator of object under path.
    it, err := store.List("path")
    
    for {
    	// Use iterator.Next to retrieve next object until we meet IteratorDone.
    	o, err := it.Next()
    	if errors.Is(err, types.IteraoorDone) {
    		break
        }
    }

    // Delete hello.txt
    err = store.Delete("hello.txt")
}
```

More examples could be found at [go-storage-example](https://github.com/beyondstorage/go-storage-example).

## Features

### Widely native services support

**9** stable services that have passed all [integration tests](https://github.com/beyondstorage/go-integration-test).

- [azblob](https://github.com/beyondstorage/go-service-azblob/): [Azure Blob storage](https://docs.microsoft.com/en-us/azure/storage/blobs/)
- [cos](https://github.com/beyondstorage/go-service-cos/): [Tencent Cloud Object Storage](https://cloud.tencent.com/product/cos)
- [dropbox](https://github.com/beyondstorage/go-service-dropbox/): [Dropbox](https://www.dropbox.com)
- [fs](https://github.com/beyondstorage/go-service-fs/): Local file system
- [gcs](https://github.com/beyondstorage/go-service-gcs/): [Google Cloud Storage](https://cloud.google.com/storage/)
- [kodo](https://github.com/beyondstorage/go-service-kodo/): [qiniu kodo](https://www.qiniu.com/products/kodo)
- [oss](https://github.com/beyondstorage/go-service-oss/): [Aliyun Object Storage](https://www.aliyun.com/product/oss)
- [qingstor](https://github.com/beyondstorage/go-service-qingstor/): [QingStor Object Storage](https://www.qingcloud.com/products/qingstor/)
- [s3](https://github.com/beyondstorage/go-service-s3/): [Amazon S3](https://aws.amazon.com/s3/)

**1** beta services that implemented required functions, but not passed [integration tests](https://github.com/beyondstorage/go-integration-test).

- [uss](https://github.com/beyondstorage/go-service-uss/): [UPYUN Storage Service](https://www.upyun.com/products/file-storage)

**11** alpha services that still under development.

- [ftp](https://github.com/beyondstorage/go-service-ftp/): FTP
- [gdrive](https://github.com/beyondstorage/go-service-gdrive): [Google Drive](https://www.google.com/drive/)
- [hdfs](https://github.com/beyondstorage/go-service-hdfs): [Hadoop Distributed File System](https://hadoop.apache.org/docs/r1.2.1/hdfs_design.html#Introduction)
- [ipfs](https://github.com/beyondstorage/go-service-ipfs): [InterPlanetary File System](https://ipfs.io)
- [memory](https://github.com/beyondstorage/go-service-memory): data that only in memory
- [minio](https://github.com/beyondstorage/go-service-minio): [MinIO](https://min.io)
- [onedrive](https://github.com/beyondstorage/go-service-onedrive): [Microsoft OneDrive](https://www.microsoft.com/en-ww/microsoft-365/onedrive/online-cloud-storage)
- [storj](https://github.com/beyondstorage/go-service-storj): [StorJ](https://www.storj.io/)
- [tar](https://github.com/beyondstorage/go-service-tar): tar files
- [webdav](https://github.com/beyondstorage/go-service-webdav): [WebDAV](http://www.webdav.org/)
- [zip](https://github.com/beyondstorage/go-service-zip): zip files

More service ideas could be found at [Service Integration Tracking](https://github.com/beyondstorage/go-storage/issues/536).

### Complete and easily extensible interface

Basic operations

- Metadata: get `Storager` metadata
```go
meta := store.Metadata()
_ := meta.GetWorkDir() // Get object WorkDir
_, ok := meta.GetWriteSizeMaximum() // Get the maximum size for write operation
```
- Read: read `Object` content
```go
// Read 2048 byte at the offset 1024 into the io.Writer.
n, err := store.Read("path", w, pairs.WithOffset(1024), pairs.WithSize(2048))
```
- Write: write content into `Object`
```go
// Write 2048 byte from io.Reader
n, err := store.Write("path", r, 2048)
```
- Stat: get `Object` metadata or check existences
```go
o, err := store.Stat("path")
if errors.Is(err, services.ErrObjectNotExist) {
	// object is not exist
}
length, ok := o.GetContentLength() // get the object content length.
```
- Delete: delete an `Object`
```go
err := store.Delete("path") // Delete the object "path"
```
- List: list `Object` in given prefix or dir
```go
it, err := store.List("path")
for {
	o, err := it.Next()
	if err != nil && errors.Is(err, types.IteratorDone) {
        // the list is over 
    }
    length, ok := o.GetContentLength() // get the object content length.
}
```

Extended operations

- Copy: copy a `Object` inside storager
```go
err := store.(Copier).Copy(src, dst) // Copy an object from src to dst.
```
- Move: move a `Object` inside storager
```go
err := store.(Mover).Move(src, dst) // Move an object from src to dst.
```
- Reach: generate a public accessible url to an `Object`
```go
url, err := store.(Reacher).Reach("path") // Generate an url to the object.
```
- Dir: Dir `Object` support
```go
o, err := store.(Direr).CreateDir("path") // Create a dir object.
```

Large file manipulation

- Multipart: allow doing multipart uploads
```go
ms := store.(Multiparter)

// Create a multipart object.
o, err := ms.CreateMultipart("path")
// Write 1024 bytes from io.Reader into a multipart at index 1
n, part, err := ms.WriteMultipart(o, r, 1024, 1)
// Complete a multipart object.
err := ms.CompleteMultipart(o, []*Part{part})
```
- Append: allow appending to an object
```go
as := store.(Appender)

// Create an appendable object.
o, err := as.CreateAppend("path")
// Write 1024 bytes from io.Reader.
n, err := as.WriteAppend(o, r, 1024)
// Commit an append object.
err = as.CommitAppend(o)
```
- Block: allow combining an object with block ids
```go
bs := store.(Blocker)

// Create a block object.
o, err := bs.CreateBlock("path")
// Write 1024 bytes from io.Reader with block id "id-abc"
n, err := bs.WriteBlock(o, r, 1024, "id-abc")
// Combine block via block ids.
err := bs.CombineBlock(o, []string{"id-abc"})
```
- Page: allow doing random writes
```go
ps := store.(Pager)

// Create a page object.
o, err := ps.CreatePage("path")
// Write 1024 bytes from io.Reader at offset 2048
n, err := ps.WritePage(o, r, 1024, 2048)
```

### Comprehensive metadata

Global object metadata

- `id`: unique key in service
- `name`: relative path towards service's work dir
- `mode`: object mode can be a combination of `read`, `dir`, `part` and [more](https://github.com/beyondstorage/go-storage/blob/master/types/object.go#L11) 
- `etag`: entity tag as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.19)
- `content-length`: object's content size.
- `content-md5`: md5 digest as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.15)
- `content-type`: media type as defined in [rfc2616](https://tools.ietf.org/html/rfc2616#section-14.17)
- `last-modified`: object's last updated time.

System object metadata

Service system object metadata like `storage-class` and so on.

```go
o, err := store.Stat("path")

// Get service system metadata via API provides by go-service-s3.
om := s3.GetObjectSystemMetadata(o)
_ = om.StorageClass // this object's storage class
_ = om.ServerSideEncryptionCustomerAlgorithm // this object's sse algorithm
```

### Strong Typing Everywhere

Self maintained codegen [definitions](https://github.com/beyondstorage/go-storage/tree/master/cmd/definitions) helps to generate all our APIs, pairs and metadata.

Generated pairs which can be used as API optional arguments.

```go
func WithContentMd5(v string) Pair {
    return Pair{
        Key:   "content_md5",
        Value: v,
    }
}
```

Generated object metadata which can be used to get content md5 from object.

```go
func (o *Object) GetContentMd5() (string, bool) {
    o.stat()
    
    if o.bit&objectIndexContentMd5 != 0 {
        return o.contentMd5, true
    }
    
    return "", false
}
```

### Server-Side Encrypt

Server-Side Encrypt supports via system pair and system metadata, and we can use [Default Pairs](https://beyondstorage.io/docs/go-storage/pairs/index#default-pairs) to simplify the job.

```go

func NewS3SseC(key []byte) (types.Storager, error) {
    defaultPairs := s3.DefaultStoragePairs{
        Write: []types.Pair{
            // Required, must be AES256
            s3.WithServerSideEncryptionCustomerAlgorithm(s3.ServerSideEncryptionAes256),
            // Required, your AES-256 key, a 32-byte binary value
            s3.WithServerSideEncryptionCustomerKey(key),
        },
        // Now you have to provide customer key to read encrypted data
        Read: []types.Pair{
            // Required, must be AES256
            s3.WithServerSideEncryptionCustomerAlgorithm(s3.ServerSideEncryptionAes256),
            // Required, your AES-256 key, a 32-byte binary value
            s3.WithServerSideEncryptionCustomerKey(key),
        }}
    
    return s3.NewStorager(..., s3.WithDefaultStoragePairs(defaultPairs))
}
```

## Sponsor

<a href="https://vercel.com?utm_source=beyondstorage&utm_campaign=oss">
    <img src="./docs/images/vercel.svg">
</a>
