package headers

import (
	"strconv"
)

// FormatRange will format a valid http Range header
func FormatRange(offset, size int64) string {
	offs := strconv.FormatInt(offset, 10)
	if size == 0 {
		return "bytes=" + offs + "-"
	}
	end := offset + size - 1
	return "bytes=" + offs + "-" + strconv.FormatInt(end, 10)
}
