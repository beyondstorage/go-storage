package fs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/qingstor/go-mime"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = os.Remove(rp)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// Omit `file not exist` error here
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}

type listDirInput struct {
	rp  string
	dir string

	started           bool
	continuationToken string

	f    *os.File
	buf  *[]byte
	bufp int
}

func (input *listDirInput) ContinuationToken() string {
	return input.continuationToken
}

func (s *Storage) commitAppend(ctx context.Context, o *Object, opt pairStorageCommitAppend) (err error) {
	return
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	srcFile, needClose, err := s.openFile(rs, os.O_RDONLY)
	if err != nil {
		return err
	}
	if needClose {
		defer srcFile.Close()
	}

	dstFile, needClose, err := s.createFile(rd)
	if err != nil {
		return err
	}
	if needClose {
		defer dstFile.Close()
	}

	_, err = io.CopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return err
	}
	return
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(false)
		o.Mode = ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = ModeRead
	}

	o.ID = filepath.Join(s.workDir, path)
	o.Path = path
	return o
}

func (s *Storage) createAppend(ctx context.Context, path string, opt pairStorageCreateAppend) (o *Object, err error) {
	rp := s.getAbsPath(path)

	f, needClose, err := s.createFile(rp)
	if err != nil {
		return
	}
	if needClose {
		err = f.Close()
		if err != nil {
			return
		}
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode = ModeRead | ModeAppend
	o.SetAppendOffset(0)

	return o, nil
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	rp := s.getAbsPath(path)

	err = os.MkdirAll(rp, 0755)
	if err != nil {
		return
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeDir
	return
}

func (s *Storage) createLink(ctx context.Context, path string, target string, opt pairStorageCreateLink) (o *Object, err error) {
	rt := s.getAbsPath(target)
	rp := s.getAbsPath(path)

	fi, err := os.Lstat(rp)
	if err == nil {
		// File exists. If the file is a symlink, then we remove it.
		if fi.Mode()&os.ModeSymlink != 0 {
			err = os.Remove(rp)
			if err != nil {
				return nil, err
			}
		} else {
			// File exists, but is not a symlink.
			return nil, services.ErrObjectModeInvalid
		}
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Something error other than ErrNotExist happened, return directly.
		return nil, err
	}

	// Set stat error to nil
	err = nil

	// The file is not exist, we should create the dir and create the file
	if fi == nil {
		err = os.MkdirAll(filepath.Dir(rp), 0755)
		if err != nil {
			return nil, err
		}
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.SetLinkTarget(rt)

	o.Mode |= ModeLink

	err = os.Symlink(rt, rp)
	if err != nil {
		return nil, err
	}

	return
}

func (s *Storage) fetch(ctx context.Context, path string, url string, opt pairStorageFetch) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// TODO: Use go-storage http client instead
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 403:
			err = os.ErrPermission
		case 404:
			err = os.ErrNotExist
		default:
			return fmt.Errorf("fetch from url %s expected %d, but got %d", url, http.StatusOK, resp.StatusCode)
		}
		return fmt.Errorf("%w: fetch from url %s expected %d, but got %d", err, url, http.StatusOK, resp.StatusCode)
	}
	_, err = s.WriteWithContext(ctx, path, resp.Body, resp.ContentLength)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	buf := make([]byte, 8192)

	input := listDirInput{
		// Always keep service original name as rp.
		rp: s.getAbsPath(path),
		// Then convert the dir to slash separator.
		dir: filepath.ToSlash(path),

		// if HasContinuationToken, we should start after we scanned this token.
		// else, we can start directly.
		started:           !opt.HasContinuationToken,
		continuationToken: opt.ContinuationToken,

		buf: &buf,
	}

	return NewObjectIterator(ctx, s.listDirNext, &input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	fi, err := os.Lstat(rd)
	if err == nil {
		// File is exist, let's check if the file is a dir or a symlink.
		if fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 {
			return services.ErrObjectModeInvalid
		}
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Something error other than ErrNotExist happened, return directly.
		return
	}
	// Set stat error to nil.
	err = nil

	// The file is not exist, we should create the dir and create the file.
	if fi == nil {
		err = os.MkdirAll(filepath.Dir(rd), 0755)
		if err != nil {
			return err
		}
	}

	err = os.Rename(rs, rd)
	if err != nil {
		return err
	}
	return
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	var rc io.ReadCloser

	rp := s.getAbsPath(path)

	f, needClose, err := s.openFile(rp, os.O_RDONLY)
	if err != nil {
		return
	}
	if needClose {
		defer func() {
			closeErr := f.Close()
			// Only return close error while copy without error
			if err == nil {
				err = closeErr
			}
		}()
	}

	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return n, err
		}
	}

	rc = f

	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)

	fi, err := s.statFile(rp)
	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if fi.IsDir() {
		o.Mode |= ModeDir
		return
	}

	if fi.Mode().IsRegular() {
		o.Mode |= ModeRead | ModePage | ModeAppend

		o.SetContentLength(fi.Size())
		o.SetLastModified(fi.ModTime())

		if v := mime.DetectFilePath(path); v != "" {
			o.SetContentType(v)
		}
	}

	// Check if this file is a link.
	if fi.Mode()&os.ModeSymlink != 0 {
		o.Mode |= ModeLink

		target, err := evalSymlinks(rp)
		if err != nil {
			return nil, err
		}
		o.SetLinkTarget(target)
	}

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	var f io.WriteCloser

	rp := s.getAbsPath(path)

	f, needClose, err := s.createFile(rp)
	if err != nil {
		return
	}
	if needClose {
		defer f.Close()
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	return io.CopyN(f, r, size)
}

func (s *Storage) writeAppend(ctx context.Context, o *Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
	f, needClose, err := s.createFileWithFlag(o.ID, os.O_RDWR|os.O_CREATE|os.O_APPEND)
	if err != nil {
		return
	}
	if needClose {
		defer f.Close()
	}

	return io.CopyN(f, r, size)
}
