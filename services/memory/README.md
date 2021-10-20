[![Services Test Memory](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-memory.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-memory.yml)

# memory

memory service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/memory
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/memory"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("memory:///path/to/workdir")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
