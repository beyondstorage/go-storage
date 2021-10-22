package fs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// normVolumeName is like VolumeName, but makes drive letter upper case.
// result of EvalSymlinks must be unique, so we have
// EvalSymlinks(`c:\a`) == EvalSymlinks(`C:\a`).
func normVolumeName(path string) string {
	volume := filepath.VolumeName(path)

	if len(volume) > 2 { // isUNC
		return volume
	}

	return strings.ToUpper(volume)
}

// normBase returns the last element of path with correct case.
func normBase(path string) (string, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	var data syscall.Win32finddata

	h, err := syscall.FindFirstFile(p, &data)
	if err != nil {
		return "", err
	}
	syscall.FindClose(h)

	return syscall.UTF16ToString(data.FileName[:]), nil
}

// baseIsDotDot reports whether the last element of path is "..".
// The given path should be 'Clean'-ed in advance.
func baseIsDotDot(path string) bool {
	i := strings.LastIndexByte(path, filepath.Separator)
	return path[i+1:] == ".."
}

// toNorm returns the normalized path that is guaranteed to be unique.
// It should accept the following formats:
//   * UNC paths                              (e.g \\server\share\foo\bar)
//   * absolute paths                         (e.g C:\foo\bar)
//   * relative paths begin with drive letter (e.g C:foo\bar, C:..\foo\bar, C:.., C:.)
//   * relative paths begin with '\'          (e.g \foo\bar)
//   * relative paths begin without '\'       (e.g foo\bar, ..\foo\bar, .., .)
// The returned normalized path will be in the same form (of 5 listed above) as the input path.
// If two paths A and B are indicating the same file with the same format, toNorm(A) should be equal to toNorm(B).
// The normBase parameter should be equal to the normBase func, except for in tests.  See docs on the normBase func.
func toNorm(path string, normBase func(string) (string, error)) (string, error) {
	if path == "" {
		return path, nil
	}

	path = filepath.Clean(path)

	volume := normVolumeName(path)
	path = path[len(volume):]

	// skip special cases
	if path == "" || path == "." || path == `\` {
		return volume + path, nil
	}

	var normPath string

	for {
		if baseIsDotDot(path) {
			normPath = path + `\` + normPath

			break
		}

		name, err := normBase(volume + path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return volume + path, nil
			}
			return "", err
		}

		normPath = name + `\` + normPath

		i := strings.LastIndexByte(path, filepath.Separator)
		if i == -1 {
			break
		}
		if i == 0 { // `\Go` or `C:\Go`
			normPath = `\` + normPath

			break
		}

		path = path[:i]
	}

	normPath = normPath[:len(normPath)-1] // remove trailing '\'

	return volume + normPath, nil
}

// evalSymlinks returns the path name after the evaluation of any symbolic
// links.
// The original implementation can be referenced to filepath.EvalSymlinks,
// but it will return the current path other then ENOENT error while walkSymlinks.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
// EvalSymlinks calls Clean on the result.
func evalSymlinks(path string) (string, error) {
	newpath, err := walkSymlinks(path)
	if err != nil {
		return "", err
	}
	newpath, err = toNorm(newpath, normBase)
	if err != nil {
		return "", err
	}
	return newpath, nil
}
