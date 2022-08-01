[![Services Test Bos](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-bos.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-bos.yml)

# bos

BOS(Baidu Object Storage) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-storage/services/bos/v2
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-storage/services/bos/v2"
	"github.com/beyondstorage/go-storage/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("bos://bucket_name/path/to/workdir")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/bos) about go-service-bos.
