package define

// All available options for storage.
//go:generate go run ../internal/cmd/options_gen/main.go
const (
	OptionStorageClass = "storage_class"
	OptionMd5          = "md5"
)

// Option will store option for storage service.
type Option struct {
	Key   string
	Value interface{}
}
