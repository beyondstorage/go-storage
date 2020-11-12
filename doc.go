/*
Package storage intend to provide a unified storage layer for Golang.

Goals

- Production ready: high test coverage, enterprise storage software adaptation, semantic versioning, well documented.

- High performance: more code generation, less runtime reflect.

- Vendor agnostic: more generic abstraction, less internal details.

Examples

The most common case to use a Storager service could be following:

1. Init a storager.

    store, err := fs.NewStorager(pairs.WithWorkDir("/tmp"))
	if err != nil {
		log.Fatalf("service init failed: %v", err)
	}

2. Use Storager API to maintain data.

	var buf bytes.Buffer

	n, err := store.Read("path/to/file", &buf)
	if err != nil {
		log.Printf("storager read: %v", err)
	}

*/
package storage

//go:generate go run -tags tools ./cmd/definitions
//go:generate go run -tags tools ./internal/cmd/iterator
