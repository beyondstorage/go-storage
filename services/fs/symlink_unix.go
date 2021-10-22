//go:build !windows
// +build !windows

package fs

// evalSymlinks returns the path name after the evaluation of any symbolic
// links.
// The original implementation can be referenced to filepath.EvalSymlinks,
// but it will return the current path other then ENOENT error while walkSymlinks.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
// EvalSymlinks calls Clean on the result.
func evalSymlinks(path string) (string, error) {
	return walkSymlinks(path)
}
