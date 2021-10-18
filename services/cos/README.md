[![Services Test Cos](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-cos.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-cos.yml)

# cos

[COS(Cloud Object Storage)](https://cloud.tencent.com/product/cos) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/cos/v3
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/v5/services/cos/v3"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("cos://bucket_name/path/to/workdir?credential=hmac:<account_name>:<account_key>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/cos) about go-service-cos.
