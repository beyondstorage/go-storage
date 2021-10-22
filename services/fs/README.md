[![Services Test Fs](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-fs.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-fs.yml)

# fs

Local file system service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/fs/v4
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/fs/v4"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("fs:///path/to/workdir")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/fs) about go-service-fs.
