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

    srv := qingstor.New()
	err := srv.Init(
		pairs.WithAccessKey("test_access_key"),
		pairs.WithSecretKey("test_secret_key"),
	)
	if err != nil {
		log.Printf("service init failed: %v", err)
	}

2. Get a storage instance from an initiated service.

	store, err := srv.Get("test_bucket_name")
	if err != nil {
		log.Printf("service get bucket failed: %v", err)
	}

3. Init a storage.

	err := store.Init(pairs.WithWorkDir("/prefix"))
	if err != nil {
		log.Printf("storager init failed: %v", err)
	}

4. Use Storager API to maintain data.

	ch := make(chan *types.Object, 1)
	defer close(ch)

	err := store.ListDir("prefix", pairs.WithFileFunc(func(*types.Object){
		ch <- o
	}))
	if err != nil {
		log.Printf("storager listdir failed: %v", err)
	}

Notes

- Storage uses error wrapping added by go 1.13, go version before 1.13 could be behaved as unexpected.
*/
package storage
