package types

// Informer will be returned by stat, and should be a *File, *Dir or a *Stream.
type Informer interface{}

// File represents a seekable file or object.
type File struct {
	// Name must a complete path instead of basename in POSIX.
	Name string
	Size int64
	Type string

	CheckSum string
	Metadata map[string]string
}

// Stream represents a not seekable stream.
type Stream struct {
	// Name must a complete path instead of basename in POSIX.
	Name string
	Type string
}

// Dir represents a virtual directory which contains files or streams.
type Dir struct {
	// Name must a complete path instead of basename in POSIX.
	Name string
}
