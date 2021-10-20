package ipfs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	ipfs "github.com/ipfs/go-ipfs-api"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	. "go.beyondstorage.io/v5/types"
)

// The src of `ipfs files cp` supports both `IPFS-path` and `MFS-path`
// After `s.getAbsPath(src)`, if the absolute path matches `IPFS-path`, it will take precedence
// This means that if the `workDir` is `/ipfs/`, there is a high probability that an error will be returned
// See https://github.com/beyondstorage/specs/pull/134#discussion_r663594807 for more details
func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {
	dst = s.getAbsPath(dst)
	stat, err := s.ipfs.FilesStat(ctx, dst)
	if err == nil {
		if stat.Type == "directory" {
			return services.ErrObjectModeInvalid
		} else {
			err = s.ipfs.FilesRm(ctx, dst, true)
			if err != nil {
				return err
			}
		}
	} else if !errors.Is(formatError(err), services.ErrObjectNotExist) {
		return err
	}
	return s.ipfs.FilesCp(ctx, s.getAbsPath(src), dst)
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		path += "/"
		o = NewObject(s, true)
		o.Mode = ModeDir
	} else {
		o = NewObject(s, false)
		o.Mode = ModeRead
	}
	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *Object, err error) {
	path = s.getAbsPath(path)
	err = s.ipfs.FilesMkdir(ctx, path, ipfs.FilesMkdir.Parents(true))
	if err != nil {
		return nil, err
	}
	o = NewObject(s, true)
	o.ID = path
	o.Path = path
	o.Mode = ModeDir
	return
}

// GSP-46: Idempotent Storager Delete Operation
// ref: https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md
func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	err = s.ipfs.FilesRm(ctx, s.getAbsPath(path), true)
	return
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *ObjectIterator, err error) {
	rp := s.getAbsPath(path)
	if !opt.HasListMode || opt.ListMode.IsDir() {
		nextFn := func(ctx context.Context, page *ObjectPage) error {
			dir, err := s.ipfs.FilesLs(ctx, rp, ipfs.FilesLs.Stat(true))
			if err != nil {
				return err
			}
			for _, f := range dir {
				o := NewObject(s, true)
				o.ID = f.Hash
				o.Path = f.Name
				switch f.Type {
				case ipfs.TFile:
					o.Mode |= ModeRead
				case ipfs.TDirectory:
					o.Mode |= ModeDir
				}
				o.SetContentLength(int64(f.Size))
				page.Data = append(page.Data, o)
			}
			return IterateDone
		}
		oi = NewObjectIterator(ctx, nextFn, nil)
		return
	} else {
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *StorageMeta) {
	meta = NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) move(ctx context.Context, src string, dst string, opt pairStorageMove) (err error) {
	dst = s.getAbsPath(dst)
	stat, err := s.ipfs.FilesStat(ctx, dst)
	if err == nil {
		if stat.Type == "directory" {
			return services.ErrObjectModeInvalid
		}
	} else if !errors.Is(formatError(err), services.ErrObjectNotExist) {
		return err
	}
	return s.ipfs.FilesMv(ctx, s.getAbsPath(src), s.getAbsPath(dst))
}

func (s *Storage) querySignHTTPDelete(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPDelete) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) querySignHTTPRead(ctx context.Context, path string, expire time.Duration, opt pairStorageQuerySignHTTPRead) (req *http.Request, err error) {
	rp := s.getAbsPath(path)
	stat, err := s.ipfs.FilesStat(ctx, rp, ipfs.FilesStat.WithLocal(true))
	if err != nil {
		return nil, err
	}
	if stat.Type != "file" {
		return nil, errors.New("path not a file")
	}

	return http.NewRequest(http.MethodGet, s.gateway+"/ipfs/"+stat.Hash, nil)
}

func (s *Storage) querySignHTTPWrite(ctx context.Context, path string, size int64, expire time.Duration, opt pairStorageQuerySignHTTPWrite) (req *http.Request, err error) {
	panic("not implemented")
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	fileOpts := make([]ipfs.FilesOpt, 0)
	if opt.HasOffset {
		fileOpts = append(fileOpts, ipfs.FilesRead.Offset(opt.Offset))
	}
	if opt.HasSize {
		fileOpts = append(fileOpts, ipfs.FilesRead.Count(opt.Size))
	}
	f, err := s.ipfs.FilesRead(ctx, s.getAbsPath(path), fileOpts...)
	if err != nil {
		return 0, err
	}
	if opt.HasIoCallback {
		f = iowrap.CallbackReadCloser(f, opt.IoCallback)
	}
	return io.Copy(w, f)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *Object, err error) {
	rp := s.getAbsPath(path)
	stat, err := s.ipfs.FilesStat(ctx, rp, ipfs.FilesStat.WithLocal(true))
	if err != nil {
		return nil, err
	}
	o = NewObject(s, true)
	o.ID = stat.Hash
	o.Path = path
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= ModeDir
	} else {
		o.Mode |= ModeRead
	}
	o.SetContentType(stat.Type)
	o.SetContentLength(int64(stat.Size))
	var sm ObjectSystemMetadata
	sm.Hash = stat.Hash
	sm.Blocks = stat.Blocks
	sm.Local = stat.Local
	sm.WithLocality = stat.WithLocality
	sm.CumulativeSize = stat.CumulativeSize
	sm.SizeLocal = stat.SizeLocal
	o.SetSystemMetadata(sm)
	return
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if r == nil {
		if size > 0 {
			return 0, errors.New("size is not 0 when io.Reader is nil")
		}
		r = bytes.NewReader([]byte{})
	}
	r = io.LimitReader(r, size)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	err = s.ipfs.FilesWrite(
		ctx, s.getAbsPath(path), r,
		ipfs.FilesWrite.Create(true),
		ipfs.FilesWrite.Parents(true),
		ipfs.FilesWrite.Truncate(true),
		ipfs.FilesWrite.Count(size),
	)
	if err != nil {
		return 0, err
	}
	return size, nil
}
