package ftp

import (
	"bytes"
	"context"
	"errors"
	"github.com/jlaffaye/ftp"
	mime "github.com/qingstor/go-mime"
	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
	"io"
	"net/textproto"
	"path/filepath"
)

type listDirInput struct {
	rp                string
	dir               string
	started           bool
	continuationToken string
	objList           []*ftp.Entry
	counter           int
}

func (input *listDirInput) ContinuationToken() string {
	return input.continuationToken
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(false)
		o.Mode = ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = ModeRead
	}
	path = filepath.ToSlash(path)
	o.ID = filepath.Join(s.workDir, path)
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string) (o *Object, err error) {
	rp := s.getAbsPath(path)
	err = s.connection.MakeDir(rp)
	if err != nil {
		return nil, err
	}
	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	o.Mode |= ModeDir
	return
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)
	err = s.connection.Delete(rp)
	if err != nil {
		var txtErr *textproto.Error
		// ignore error with code ftp.StatusFileUnavailable, to make delete idempotent
		if errors.As(err, &txtErr) && (txtErr.Code == ftp.StatusFileUnavailable || txtErr.Code == ftp.StatusFileActionIgnored) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	if !opt.HasListMode || opt.ListMode.IsDir() {
		input := listDirInput{
			// Always keep service original name as rp.
			rp: s.getAbsPath(path),
			// Then convert the dir to slash separator.
			dir: filepath.ToSlash(path),
			// if HasContinuationToken, we should start after we scanned this token.
			// else, we can start directly.
			started:           !opt.HasContinuationToken,
			continuationToken: opt.ContinuationToken,
			counter:           0,
		}
		return NewObjectIterator(ctx, s.listDirNext, &input), nil
	} else {
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)
	var offset uint64 = 0
	if opt.HasOffset {
		offset = uint64(opt.Offset)
	}
	r, err := s.connection.RetrFrom(rp, offset)
	if err != nil {
		return n, err
	}
	defer func() {
		closeErr := r.Close()
		if err == nil {
			err = closeErr
		}
	}()

	var rc io.ReadCloser
	rc = r
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}
	if opt.HasSize {
		return io.CopyN(w, rc, opt.Size)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)
	fl, err := s.connection.List(rp)
	if err != nil {
		flst, err := s.connection.List(filepath.Dir(rp))
		if err != nil {
			return nil, err
		}
		for i := range flst {
			if filepath.Base(rp) == flst[i].Name {
				fl = []*ftp.Entry{flst[i]}
				break
			}
		}
	}
	if len(fl) == 0 {
		return nil, services.ErrObjectNotExist
	}
	var fe *ftp.Entry = fl[0]
	if fe == nil {
		return nil, services.ErrObjectNotExist
	}
	o = s.newObject(true)
	o.ID = rp
	o.Path = path
	switch fe.Type {
	case ftp.EntryTypeFolder:
		o.Mode |= ModeDir
		return
	case ftp.EntryTypeLink:
		o.Mode |= ModeLink
		target := fe.Target
		if err != nil {
			return nil, err
		}
		o.SetLinkTarget(target)
	default:
		o.Mode |= ModeRead | ModePage | ModeAppend
		o.SetContentLength(int64(fe.Size))
		o.SetLastModified(fe.Time)
		if v := mime.DetectFilePath(path); v != "" {
			o.SetContentType(v)
		}
	}
	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	rp := s.getAbsPath(path)
	err = s.makeDir(filepath.Dir(rp))
	if err != nil {
		return
	}
	if size == 0 {
		r = bytes.NewReader([]byte{})
	}
	if r == nil {
		return 0, services.ErrObjectNotExist
	}
	lr := io.LimitReader(r, size)
	if opt.HasIoCallback {
		lr = iowrap.CallbackReader(lr, opt.IoCallback)
	}
	err = s.connection.Stor(rp, lr)
	if err != nil {
		return
	}
	return
}
