[![Services Test Oss](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-oss.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-oss.yml)

# oss

[Aliyun Object Storage](https://cn.aliyun.com/product/oss) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/oss/v3
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/oss/v3"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("oss://bucket_name/path/to/workdir?credential=hmac:<access_key>:<secret_key>&endpoint=https:<location>.aliyuncs.com")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/oss) about go-service-oss. 
