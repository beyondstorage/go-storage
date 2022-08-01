# uss

[UPYUN Storage Service](https://www.upyun.com/products/file-storage) support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get github.com/beyondstorage/go-storage/services/uss/v3
```

## Usage

```go
import (
	"log"

	_ "github.com/beyondstorage/go-storage/services/uss/v3"
	"github.com/beyondstorage/go-storage/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("uss://bucket_name/path/to/workdir?credential=hmac:<operator>:<password>&endpoint=https:<domain>")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/uss) about go-service-uss.
