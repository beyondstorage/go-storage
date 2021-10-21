[![Services Test S3](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-s3.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-s3.yml)

# s3

AWS S3 service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/s3/v3
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/s3/v3"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("s3://bucket_name/path/to/workdir")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/s3) about go-service-s3. 

## Compatible Services

We can use go-service-s3 for the following services:

- [Aliyun OSS S3 Compatible API](https://help.aliyun.com/apsara/agile-data/v_2_5_0_20200506/oss/insight-developer-guide/s3-api-compatibility-instructions.html) (We also provide native support in [go-service-oss](https://github.com/beyondstorage/go-service-oss))
- [AWS S3](https://aws.amazon.com/s3/) (The native support service.)
- [DigitalOcean Space](https://www.digitalocean.com/products/spaces/)
- [ECloud (China Mobile Cloud) Object Storage](https://www.ctyun.cn/products/10020000)
- [GCS S3 Compatible API](https://cloud.google.com/storage/docs/interoperability) (We also provide native support in [go-service-gcs](https://github.com/beyondstorage/go-service-gcs))
- [IBM Cloud Storage Service](https://www.ibm.com/cloud/storage)
- [ksyun KS3](https://www.ksyun.com/nv/product/KS3.html)
- [JCloud Object Storage](https://www.jdcloud.com/cn/products/object-storage-service)
- [Minio](https://min.io/) (We also provide native support in [go-service-minio](https://github.com/beyondstorage/go-service-minio))
- [QingStor Object Storage S3 Compatible API](https://docs.qingcloud.com/qingstor/s3/) (We also provide native support in [go-service-qingstor](https://github.com/beyondstorage/go-service-qingstor))
- [Scaleway Object Storage](https://www.scaleway.com/en/object-storage/)
