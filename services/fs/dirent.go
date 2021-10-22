//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package fs

import (
	"syscall"
	"unsafe"
)

// readInt returns the size-bytes unsigned integer in native byte order at offset off.
func readInt(b []byte, off, size uintptr) (u uint64, ok bool) {
	if len(b) < int(off+size) {
		return 0, false
	}
	if isBigEndian {
		return readIntBE(b[off:], size), true
	}
	return readIntLE(b[off:], size), true
}

func readIntBE(b []byte, size uintptr) uint64 {
	switch size {
	case 1:
		return uint64(b[0])
	case 2:
		_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[1]) | uint64(b[0])<<8
	case 4:
		_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[3]) | uint64(b[2])<<8 | uint64(b[1])<<16 | uint64(b[0])<<24
	case 8:
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
			uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	default:
		panic("syscall: readInt with unsupported size")
	}
}

func readIntLE(b []byte, size uintptr) uint64 {
	switch size {
	case 1:
		return uint64(b[0])
	case 2:
		_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8
	case 4:
		_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24
	case 8:
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	default:
		panic("syscall: readInt with unsupported size")
	}
}

/*
type Dirent struct {
	Ino    uint64      // 64-bit inode number
	Off    int64       // 64-bit offset to next structure
	Reclen uint16      // Size of this dirent
	Type   uint8       // File type
	Name   [256]int8   // Filename (null-terminated)
	_      [5]byte     // Zero padding byte
}
*/

const (
	direntOffsetIno    = unsafe.Offsetof(syscall.Dirent{}.Ino)
	direntSizeIno      = unsafe.Sizeof(syscall.Dirent{}.Ino)
	direntOffsetReclen = unsafe.Offsetof(syscall.Dirent{}.Reclen)
	direntSizeReclen   = unsafe.Sizeof(syscall.Dirent{}.Reclen)
	direntOffsetType   = unsafe.Offsetof(syscall.Dirent{}.Type)
	direntSizeType     = unsafe.Sizeof(syscall.Dirent{}.Type)
	direntOffsetName   = unsafe.Offsetof(syscall.Dirent{}.Name)
)

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, direntOffsetIno, direntSizeIno)
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, direntOffsetReclen, direntSizeReclen)
}

func direntType(buf []byte) (uint8, bool) {
	ty, ok := readInt(buf, direntOffsetType, direntSizeType)
	return uint8(ty), ok
}
