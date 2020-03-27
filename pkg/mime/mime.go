package mime

import (
	"mime"
	"strings"
)

// TypeByFileName will get mime type via filename, and return "" if not detected.
func TypeByFileName(filename string) string {
	x := strings.Split(filename, ".")
	if len(x) < 2 {
		return ""
	}

	return mime.TypeByExtension("." + x[len(x)-1])
}
