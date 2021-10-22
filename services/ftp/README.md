[![Services Test Ftp](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-ftp.yml/badge.svg)](https://github.com/beyondstorage/go-storage/actions/workflows/services-test-ftp.yml)

# ftp

[FTP](https://datatracker.ietf.org/doc/html/rfc959) service support for [go-storage](https://github.com/beyondstorage/go-storage).

## Install

```go
go get go.beyondstorage.io/services/ftp
```

## Usage

```go
import (
	"log"

	_ "go.beyondstorage.io/services/ftp"
	"go.beyondstorage.io/v5/services"
)

func main() {
	store, err := services.NewStoragerFromString("ftp:///path/to/workdir?credential=basic:<user>:<password>&endpoint=tcp:<host>:<port>")
	if err != nil {
		log.Fatal(err)
	}

	// Write data from io.Reader into hello.txt
	n, err := store.Write("hello.txt", r, length)
}
```

- See more examples in [go-storage-example](https://github.com/beyondstorage/go-storage-example).
- Read [more docs](https://beyondstorage.io/docs/go-storage/services/ftp) about go-service-ftp. 
