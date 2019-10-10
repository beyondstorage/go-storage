package types

// Object will be returned by stat, and should be a *File, *Dir or a *Stream.
type Object interface{}

// File represents a seekable file or object.
type File struct {
	// Name must a complete path instead of basename in POSIX.
	Name string
	Size int64

	Metadata map[string]interface{}
}

// Stream represents a not seekable stream.
type Stream struct {
	// Name must a complete path instead of basename in POSIX.
	Name string

	Metadata map[string]interface{}
}

// Dir represents a virtual directory which contains files or streams.
type Dir struct {
	// Name must a complete path instead of basename in POSIX.
	Name string

	Metadata map[string]interface{}
}
