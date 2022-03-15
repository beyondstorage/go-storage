module go.beyondstorage.io/services/gdrive

go 1.16

require (
	github.com/dgraph-io/ristretto v0.1.0
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/v5 v5.0.0
	golang.org/x/oauth2 v0.0.0-20220309155454-6242fa91716a
	google.golang.org/api v0.72.0
)

replace go.beyondstorage.io/v5 => ../../
