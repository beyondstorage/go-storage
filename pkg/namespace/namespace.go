package namespace

import (
	"path/filepath"
	"strings"
)

// ParseObjectStorage will parse namespace into bucket name and prefix.
func ParseObjectStorage(namespace string) (name, prefix string) {
	x := strings.SplitN(namespace, "/", 2)
	if len(x) == 0 {
		return "", ""
	}
	if len(x) == 1 {
		return x[0], ""
	}
	return x[0], x[1]
}

// ParseLocalFS will parse namespace into path.
func ParseLocalFS(namespace string) (path string) {
	if strings.Contains(namespace, ":") {
		return namespace
	}
	return filepath.Join("/", namespace)
}
