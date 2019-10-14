package types

// ObjectType is the type for object, under layer type is string.
type ObjectType string

// All available type for object.
const (
	ObjectTypeFile    ObjectType = "file"
	ObjectTypeStream             = "stream"
	ObjectTypeDir                = "dir"
	ObjectTypeInvalid            = "invalid"
)

// Object may be a *File, *Dir or a *Stream.
type Object struct {
	// name must a complete path instead of basename in POSIX.
	Name string
	// type should be one of "file", "stream", "dir" or "invalid".
	Type ObjectType

	// metadata is the metadata of the object.
	Metadata
}

// Metadata is the metadata used in object.
type Metadata map[string]interface{}

// Pair will store option for storage service.
//
//go:generate go run ../internal/cmd/pairs_gen/main.go
//go:generate go run ../internal/cmd/metadata_gen/main.go
type Pair struct {
	Key   string
	Value interface{}
}

// AvailablePairs are all available options for storage.
// This will be used to generate options.go
var AvailablePairs = map[string]string{
	"access_key":    "string",
	"checksum":      "string",
	"count":         "int64",
	"delimiter":     "string",
	"expire":        "int",
	"host":          "string",
	"location":      "string",
	"name":          "string",
	"port":          "int",
	"protocol":      "string",
	"recursive":     "bool",
	"secret_key":    "string",
	"size":          "int64",
	"storage_class": "string",
	"type":          "string",
	"updated_at":    "time.Time",
}
