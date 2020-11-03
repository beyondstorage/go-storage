package storageclass

const (
	// Hot is the most frequently visited storage class for this service,
	// suitable for frequently accessed data.
	Hot = "hot"
	// Warm usually represents storage class for infrequent access data, maybe accessed several times a month.
	//
	// Notes:
	//   - Higher latency and lower performance compared to `Hot`
	//   - Read warm data requires extra fees (depend on services)
	Warn = "warm"
	// Cold is a storage class for archiving data which maybe accessed one or two times a year.
	// Notes:
	//   - Cold only support write operation
	//   - Depends on services' implementations, `Cold` storage class may need extra time (minutes to hours,
	//     except `gcs`) or extra API (`kodo`) even manual application (`cos`).
	//     So we will not add extra support to read `Cold` data, just return the error.
	Cold = "cold"
)
