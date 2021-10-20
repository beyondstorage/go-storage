[![Services Test Ipfs](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-ipfs.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-ipfs.yml)

# ipfs

[InterPlanetary File System(IPFS)](https://ipfs.io/) support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/ipfs
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/ipfs"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("ipfs:///path/to/workdir?endpoint=<ipfs_http_api_endpoint>&gateway=<ipfs_http_gateway>")
	if err != nil {
		log.Fatal(err)
	}
	
	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/ipfs) about go-service-ipfs.
