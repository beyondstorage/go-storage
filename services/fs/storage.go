package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/pkg/mime"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	err = s.osRemove(rp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	rp := s.getAbsPath(dir)

	fi, err := s.ioutilReadDir(rp)
	if err != nil {
		return err
	}

	for _, v := range fi {
		o := &types.Object{
			ID:         filepath.Join(rp, v.Name()),
			Name:       filepath.Join(dir, v.Name()),
			Size:       v.Size(),
			UpdatedAt:  v.ModTime(),
			ObjectMeta: info.NewObjectMeta(),
		}

		if v.IsDir() {
			o.Type = types.ObjectTypeDir
			if opt.HasDirFunc {
				opt.DirFunc(o)
			}
			continue
		}

		if v := mime.TypeByFileName(v.Name()); v != "" {
			o.SetContentType(v)
		}

		o.Type = types.ObjectTypeFile
		if opt.HasFileFunc {
			opt.FileFunc(o)
		}
	}
	return
}

func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta, nil
}

func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	// If path is "-", return stdin directly.
	if path == "-" {
		f := os.Stdin
		if opt.HasSize {
			return iowrap.LimitReadCloser(f, opt.Size), nil
		}
		return f, nil
	}

	rp := s.getAbsPath(path)

	f, err := s.osOpen(rp)
	if err != nil {
		return nil, err
	}
	if opt.HasOffset {
		_, err = f.Seek(opt.Offset, 0)
		if err != nil {
			return nil, err
		}
	}

	rc = f
	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}
	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}
	return rc, nil
}

func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
	if path == "-" {
		return &types.Object{
			ID:         "-",
			Name:       "-",
			Type:       types.ObjectTypeStream,
			Size:       0,
			ObjectMeta: info.NewObjectMeta(),
		}, nil
	}

	rp := s.getAbsPath(path)

	fi, err := s.osStat(rp)
	if err != nil {
		return nil, err
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Size:       fi.Size(),
		UpdatedAt:  fi.ModTime(),
		ObjectMeta: info.NewObjectMeta(),
	}

	if fi.IsDir() {
		o.Type = types.ObjectTypeDir
		return
	}
	if fi.Mode().IsRegular() {
		if v := mime.TypeByFileName(path); v != "" {
			o.SetContentType(v)
		}

		o.Type = types.ObjectTypeFile
		return
	}
	if fi.Mode()&StreamModeType != 0 {
		o.Type = types.ObjectTypeStream
		return
	}

	o.Type = types.ObjectTypeInvalid
	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
	var f io.WriteCloser
	// If path is "-", use stdout directly.
	if path == "-" {
		f = os.Stdout
	} else {
		// Create dir for path.
		err = s.createDir(path)
		if err != nil {
			return err
		}

		rp := s.getAbsPath(path)

		f, err = s.osCreate(rp)
		if err != nil {
			return err
		}
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	if opt.HasSize {
		_, err = s.ioCopyN(f, r, opt.Size)
	} else {
		_, err = s.ioCopyBuffer(f, r, make([]byte, 1024*1024))
	}
	if err != nil {
		return err
	}
	return
}

func (s *Storage) copy(ctx context.Context, src string, dst string, opt *pairStorageCopy) (err error) {
	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	srcFile, err := s.osOpen(rs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := s.osCreate(rd)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = s.ioCopyBuffer(dstFile, srcFile, make([]byte, 1024*1024))
	if err != nil {
		return err
	}
	return
}
func (s *Storage) move(ctx context.Context, src string, dst string, opt *pairStorageMove) (err error) {

	rs := s.getAbsPath(src)
	rd := s.getAbsPath(dst)

	// Create dir for dst path.
	err = s.createDir(dst)
	if err != nil {
		return err
	}

	err = s.osRename(rs, rd)
	if err != nil {
		return err
	}
	return
}
