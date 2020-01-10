package storageclass

// Type is the type for storage class
type Type string

const (
	// Hot is default storage class for all service, suitable for frequently accessed data.
	//
	// Notes:
	//   - All service supports this kind of storage class
	//   - All services' default storage class
	//   - If a service doesn't have an idea for storage class, it provides `Hot` storage class
	Hot Type = "hot"
	// Warm usually represents storage class for infrequent access data, maybe accessed several times a month.
	//
	// Notes:
	//   - Higher latency and lower performance compared to `Hot`
	//   - Read warm data requires extra fees (depend on services)
	Warm Type = "warm"
	// Cold is a storage class for archiving data which maybe accessed one or two times a year.
	//
	// Notes:
	//   - Cold only support write operation
	//   - Depends on services' implementations, `Cold` storage class may need extra time (minutes to hours,
	//     except `gcs`) or extra API (`kodo`) even manual application (`cos`).
	//     So we will not add extra support to read `Cold` data, just return the error.
	Cold Type = "cold"
)
