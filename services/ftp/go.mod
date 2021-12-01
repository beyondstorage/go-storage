module go.beyondstorage.io/services/ftp

go 1.16

require (
	github.com/google/uuid v1.3.0
	github.com/jlaffaye/ftp v0.0.0-20210307004419-5d4190119067
	github.com/kr/pretty v0.1.0 // indirect
	github.com/qingstor/go-mime v0.1.0
	go.beyondstorage.io/credential v1.0.0
	go.beyondstorage.io/endpoint v1.2.0
	go.beyondstorage.io/v5 v5.0.0
	golang.org/x/sys v0.0.0-20211109065445-02f5c0300f6e // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace go.beyondstorage.io/v5 => ../../
