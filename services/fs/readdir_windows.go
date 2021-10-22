package fs

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/sys/windows"

	typ "go.beyondstorage.io/v5/types"
)

func (s *Storage) listDirNext(ctx context.Context, page *typ.ObjectPage) (err error) {
	input := page.Status.(*listDirInput)

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

	// Every list dir will fetch 128 files.
	limit := 128
	var data windows.Win32finddata

	for limit > 0 {
		err = windows.FindNextFile(windows.Handle(input.f.Fd()), &data)
		// Whole dir has been read, return IterateDone to mark this iteration is done
		if err != nil && err == windows.ERROR_NO_MORE_FILES {
			return typ.IterateDone
		}
		if err != nil {
			return
		}

		name := windows.UTF16ToString(data.FileName[0:])
		if name == "." || name == ".." {
			continue
		}

		o := s.newObject(true)
		// Always keep service original name as ID.
		o.SetID(filepath.Join(input.rp, name))
		// Object's name should always be separated by slash (/)
		o.SetPath(path.Join(input.dir, name))

		o.SetContentLength(int64(data.FileSizeHigh)<<32 + int64(data.FileSizeLow))

		switch {
		case data.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0:
			o.Mode |= typ.ModeDir
		case data.FileAttributes&windows.FILE_ATTRIBUTE_NORMAL != 0:
			o.Mode |= typ.ModeRead | typ.ModePage | typ.ModeAppend
		case data.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0:
			// FILE_ATTRIBUTE_REPARSE_POINT means this is a file or directory that has
			// an associated reparse point, or a file that is a symbolic link.
			o.Mode |= typ.ModeLink
		}
		page.Data = append(page.Data, o)

		limit--
	}
	return
}
