package fs

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func volumeNameLen(path string) int {
	if len(path) < 2 {
		return 0
	}
	// with drive letter
	c := path[0]
	if path[1] == ':' && ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') {
		return 2
	}
	// is it UNC? https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx
	if l := len(path); l >= 5 && isSlash(path[0]) && isSlash(path[1]) &&
		!isSlash(path[2]) && path[2] != '.' {
		// first, leading `\\` and next shouldn't be `\`. its server name.
		for n := 3; n < l-1; n++ {
			// second, next '\' shouldn't be repeated.
			if isSlash(path[n]) {
				n++
				// third, following something characters. its share name.
				if !isSlash(path[n]) {
					if path[n] == '.' {
						break
					}
					for ; n < l; n++ {
						if isSlash(path[n]) {
							break
						}
					}
					return n
				}
				break
			}
		}
	}
	return 0
}

func walkSymlinks(path string) (string, error) {
	volLen := volumeNameLen(path)
	pathSeparator := string(os.PathSeparator)

	if volLen < len(path) && os.IsPathSeparator(path[volLen]) {
		volLen++
	}
	vol := path[:volLen]
	dest := vol
	linksWalked := 0
	var i = 0
	for start, end := volLen, volLen; start < len(path); start = end {
		for start < len(path) && os.IsPathSeparator(path[start]) {
			start++
		}
		end = start
		for end < len(path) && !os.IsPathSeparator(path[end]) {
			end++
		}

		// On Windows, "." can be a symlink.
		// We look it up, and use the value if it is absolute.
		// If not, we just return ".".
		isWindowsDot := runtime.GOOS == "windows" && path[volumeNameLen(path):] == "."

		// The next path component is in path[start:end].
		if end == start {
			// No more path components.
			break
		} else if path[start:end] == "." && !isWindowsDot {
			// Ignore path component ".".
			continue
		} else if path[start:end] == ".." {
			// Back up to previous component if possible.
			// Note that volLen includes any leading slash.

			// Set r to the index of the last slash in dest,
			// after the volume.
			var r int
			for r = len(dest) - 1; r >= volLen; r-- {
				if os.IsPathSeparator(dest[r]) {
					break
				}
			}
			if r < volLen || dest[r+1:] == ".." {
				// Either path has no slashes
				// (it's empty or just "C:")
				// or it ends in a ".." we had to keep.
				// Either way, keep this "..".
				if len(dest) > volLen {
					dest += pathSeparator
				}
				dest += ".."
			} else {
				// Discard everything since the last slash.
				dest = dest[:r]
			}
			continue
		}

		// Ordinary path component. Add it to result.

		if len(dest) > volumeNameLen(dest) && !os.IsPathSeparator(dest[len(dest)-1]) {
			dest += pathSeparator
		}

		dest += path[start:end]

		// Resolve symlink.

		i++

		fi, err := os.Lstat(dest)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return filepath.Clean(dest + path[end:]), nil
			}
			return "", err
		}

		if fi.Mode()&os.ModeSymlink == 0 {
			if !fi.Mode().IsDir() && end < len(path) {
				return "", syscall.ENOTDIR
			}
			continue
		}

		// Found symlink.

		linksWalked++
		if linksWalked > 255 {
			return "", errors.New("EvalSymlinks: too many links")
		}

		link, err := os.Readlink(dest)
		if err != nil {
			return "", err
		}

		if isWindowsDot && !filepath.IsAbs(link) {
			// On Windows, if "." is a relative symlink,
			// just return ".".
			break
		}

		path = link + path[end:]

		v := volumeNameLen(link)
		if v > 0 {
			// Symlink to drive name is an absolute path.
			if v < len(link) && os.IsPathSeparator(link[v]) {
				v++
			}
			vol = link[:v]
			dest = vol
			end = len(vol)
		} else if len(link) > 0 && os.IsPathSeparator(link[0]) {
			// Symlink to absolute path.
			dest = link[:1]
			end = 1
		} else {
			// Symlink to relative path; replace last
			// path component in dest.
			var r int
			for r = len(dest) - 1; r >= volLen; r-- {
				if os.IsPathSeparator(dest[r]) {
					break
				}
			}
			if r < volLen {
				dest = vol
			} else {
				dest = dest[:r]
			}
			end = 0
		}
	}
	return filepath.Clean(dest), nil
}
