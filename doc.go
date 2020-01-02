/*
Package storage intend to provide a unified storage layer for Golang.

Goals

- Production ready: high test coverage, enterprise storage software adaptation, semantic versioning, well documented.

- High performance: more code generation, less runtime reflect.

- Vendor lock free: more generic abstraction, less internal details.

Details

There two public interfaces: Servicer and Storager. Storager is a fully functional storage client, and Servicer is a
manager of Storager instances, which will be useful for services like object storage. For any service, Storager is
required to implement and Servicer is optional.

Examples

The most common case to use a Storager service could be following:

1. Init a service.

    srv, store, err := coreutils.Open("qingstor://hmac:test_access_key:test_secret_key@https:qingstor.com:443/test_bucket_name")
	if err != nil {
		log.Fatalf("service init failed: %v", err)
	}

2. Use Storager API to maintain data.

	ch := make(chan *types.Object, 1)
	defer close(ch)

	err := store.List("prefix", pairs.WithFileFunc(func(*types.Object){
		ch <- o
	}))
	if err != nil {
		log.Printf("storager listdir failed: %v", err)
	}

Notes

- Storage uses error wrapping added by go 1.13, go version before 1.13 could be behaved as unexpected.
*/
package storage
