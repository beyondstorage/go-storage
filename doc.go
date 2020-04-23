/*
Package storage intend to provide a unified storage layer for Golang.

Goals

- Production ready: high test coverage, enterprise storage software adaptation, semantic versioning, well documented.

- High performance: more code generation, less runtime reflect.

- Vendor lock free: more generic abstraction, less internal details.

Details

There two main public interfaces: Servicer and Storager. Storager is a fully functional storage client, and Servicer is a
manager of Storager instances, which will be useful for services like object storage. For any service, Storager is
required to implement and Servicer is optional.

Examples

The most common case to use a Storager service could be following:

1. Init a service.

    _, store, err := coreutils.Open("fs", pairs.WithWorkDir("/tmp"))
	if err != nil {
		log.Fatalf("service init failed: %v", err)
	}

2. Use Storager API to maintain data.

	r, err := store.Read("path/to/file")
	if err != nil {
		log.Printf("storager read: %v", err)
	}

Notes

- Storage uses error wrapping added by go 1.13, go version before 1.13 could be behaved as unexpected.
*/
package storage
