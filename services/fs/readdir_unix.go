//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package fs

import (
	"context"
	"golang.org/x/sys/unix"
	"os"
	"path"
	"path/filepath"

	typ "go.beyondstorage.io/v5/types"
)

// Available value for Dirent Type
//
// Copied from linux kernel <dirent.h>
// #define	DT_UNKNOWN	 0
// #define	DT_FIFO		 1
// #define	DT_CHR		 2
// #define	DT_DIR		 4
// #define	DT_BLK		 6
// #define	DT_REG		 8
// #define	DT_LNK		10
// #define	DT_SOCK		12
// #define	DT_WHT		14
const (
	// The file type is unknown.
	DirentTypeUnknown = 0
	// This is a named pipe (FIFO).
	DirentTypeFIFO = 1
	// This is a character device.
	DirentTypeCharDevice = 2
	// This is a directory.
	DirentTypeDirectory = 4
	// This is a block device.
	DirentTypeBlockDevice = 6
	// This is a regular file.
	DirentTypeRegular = 8
	// This is a symbolic link.
	DirentTypeLink = 10
	// This is a UNIX domain socket.
	DirentTypeSocket = 12
	// WhiteOut from BSD, don't know what's it mean.
	DirentTypeWhiteOut = 14
)

func (s *Storage) listDirNext(ctx context.Context, page *typ.ObjectPage) (err error) {
	input := page.Status.(*listDirInput)

	defer func() {
		err = s.formatError("list_dir_next", err, input.rp)
	}()

	defer func() {
		// Make sure file has been close every time we return an error
		if err != nil && input.f != nil {
			_ = input.f.Close()
			input.f = nil
		}
	}()

	// Open dir before we read it.
	if input.f == nil {
		input.f, err = os.Open(input.rp)
		if err != nil {
			return
		}
	}

	// Reset bufp before refill buf.
	input.bufp = 0
	n, err := unix.ReadDirent(int(input.f.Fd()), *input.buf)
	if err != nil {
		return err
	}
	if n <= 0 {
		return typ.IterateDone
	}

	for input.bufp < n {
		// Drain the buffer
		buf := (*input.buf)[input.bufp:n]

		// Get and check reclen
		reclen, ok := direntReclen(buf)
		if !ok || reclen > uint64(len(buf)) {
			return
		}
		input.bufp += int(reclen)
		// current dirent
		rec := buf[:reclen]

		// Get and check inode
		ino, ok := direntIno(rec)
		if !ok {
			break
		}
		if ino == 0 { // File absent in directory.
			continue
		}

		// Get and check type
		ty, ok := direntType(rec)
		if !ok {
			continue
		}

		// Get and check name
		name := rec[direntOffsetName:reclen]
		for i, c := range name {
			if c == 0 {
				name = name[:i]
				break
			}
		}
		// Format object
		fname := string(name)
		// Check for useless names before allocating a string.
		if fname == "." || fname == ".." {
			continue
		}

		if !input.started {
			if fname != input.continuationToken {
				continue
			}
			// ContinuationToken is the next file, we should include this file.
			input.started = true
		}

		o := s.newObject(false)
		// FIXME: filepath.Join and path.Join is really slow here, we need handle this.
		// Always keep service original name as ID.
		o.ID = filepath.Join(input.rp, fname)
		// Object's name should always be separated by slash (/)
		o.Path = path.Join(input.dir, fname)

		switch ty {
		case DirentTypeDirectory:
			o.Mode |= typ.ModeDir
		case DirentTypeRegular:
			o.Mode |= typ.ModeRead | typ.ModeAppend | typ.ModePage
		case DirentTypeLink:
			o.Mode |= typ.ModeLink
		}

		// Set update name here.
		input.continuationToken = o.Path
		page.Data = append(page.Data, o)
	}

	return
}
