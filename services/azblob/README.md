[![Services Test Azblob](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-azblob.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-azblob.yml)

# azblob

Azure azblob service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-storage/services/azblob/v3
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-storage/services/azblob/v3"
	"github.com/beyondstorage/go-storage/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("azblob://container_name/path/to/workdir?credential=hmac:<account_name>:<account_key>&endpoint=https:<account_name>.<endpoint_suffix>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/azblob) about go-service-azblob.
