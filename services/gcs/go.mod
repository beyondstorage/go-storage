module go.beyondstorage.io/services/gcs/v3

go 1.16

require (
	cloud.google.com/go/storage v1.23.0
	github.com/google/uuid v1.3.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/v5 v5.0.0
	golang.org/x/oauth2 v0.0.0-20220608161450-d0670ef3b1eb
	google.golang.org/api v0.85.0
)

replace go.beyondstorage.io/v5 => ../../
