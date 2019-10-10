package types

// Option will store option for storage service.
//
//go:generate go run ../internal/cmd/options_gen/main.go
type Option struct {
	Key   string
	Value interface{}
}

// AvailableOptions are all available options for storage.
// This will be used to generate options.go
var AvailableOptions = map[string]string{
	"location":      "string",
	"checksum":      "string",
	"storage_class": "string",
}
