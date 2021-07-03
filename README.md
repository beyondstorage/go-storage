# go-storage

[![Build Status](https://github.com/beyondstorage/go-storage/workflows/Unit%20Test/badge.svg?branch=master)](https://github.com/beyondstorage/go-storage/actions?query=workflow%3A%22Unit+Test%22)
[![Go dev](https://pkg.go.dev/badge/github.com/beyondstorage/go-storage/v4)](https://pkg.go.dev/github.com/beyondstorage/go-storage/v4)
[![License](https://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/beyondstorage/go-storage/blob/master/LICENSE)
[![go storage dev](https://img.shields.io/matrix/go-storage:aos.dev.svg?server_fqdn=chat.aos.dev&label=%23go-storage%3Aaos.dev&logo=matrix)](https://matrix.to/#/#go-storage:aos.dev) <!-- Need update after matrix updated -->

Storage abstraction that focus on neutral cross-cloud data operation.

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

## Goal

- Production ready
- High performance
- Vendor agnostic

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

### Complete and easily expandable interface

Basic operations

- Metadata: get `Storager` metadata
- Read: read `Object` content
- Write: write content into `Object`
- Stat: get `Object` metadata or check existences
- Delete: delete an `Object`
- List: list `Object` in given prefix or dir

Extended operations

- Copy: copy a `Object` inside storager
- Move: move a `Object` inside storager
- Reach: generate a public accessible url to an `Object`
- Link: Symlink `Object` support
- Dir: Dir `Object` support

Large file manipulation

- Multipart: allow doing multipart uploads
- Append: allow appending to an object
- Block: allow combining an object with block ids
- Page: allow doing random writes

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

### Strong Type Everywhere

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
