package dropbox

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the dropbox client.
type Storage struct {
	client files.Client

	workDir string
	loose   bool
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager dropbox {WorkDir: %s}",
		"/"+s.workDir,
	)
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.WorkDir = s.workDir
	m.Name = ""

	return m, nil
}

// ListDir implements Storager.ListDir
func (s *Storage) ListDir(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("list_dir", err, path)
	}()

	opt, err := parseStoragePairListDir(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	result, err := s.client.ListFolder(&files.ListFolderArg{
		Path: rp,
	})
	if err != nil {
		return err
	}

	for {
		for _, v := range result.Entries {
			switch meta := v.(type) {
			case *files.FileMetadata:
				o := &types.Object{
					ID:         meta.Id,
					Type:       types.ObjectTypeFile,
					Name:       filepath.Join(path, meta.Name),
					Size:       int64(meta.Size),
					UpdatedAt:  meta.ServerModified,
					ObjectMeta: metadata.NewObjectMeta(),
				}

				if meta.ContentHash != "" {
					o.SetETag(meta.ContentHash)
				}

				if opt.HasFileFunc {
					opt.FileFunc(o)
				}
			case *files.FolderMetadata:
				o := &types.Object{
					ID:         meta.Id,
					Type:       types.ObjectTypeDir,
					Name:       filepath.Join(path, meta.Name),
					ObjectMeta: metadata.NewObjectMeta(),
				}

				if opt.HasDirFunc {
					opt.DirFunc(o)
				}
			default:
				return ErrUnexpectedEntry
			}
		}
		if !result.HasMore {
			break
		}

		result, err = s.client.ListFolderContinue(&files.ListFolderContinueArg{
			Cursor: result.Cursor,
		})
		if err != nil {
			return err
		}
	}
	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	defer func() {
		err = s.formatError("read", err, path)
	}()

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)

	input := &files.DownloadArg{
		Path: rp,
	}

	_, r, err = s.client.Download(input)
	if err != nil {
		return nil, err
	}

	if opt.HasSize {
		r = iowrap.LimitReadCloser(r, opt.Size)
	}

	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReadCloser(r, opt.ReadCallbackFunc)
	}

	return
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("write", err, path)
	}()

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return err
	}

	rp := s.getAbsPath(path)

	if opt.HasSize {
		r = io.LimitReader(r, opt.Size)
	}
	if opt.HasReadCallbackFunc {
		r = iowrap.CallbackReader(r, opt.ReadCallbackFunc)
	}

	input := &files.CommitInfo{
		Path: rp,
		Mode: &files.WriteMode{
			Tagged: dropbox.Tagged{
				Tag: files.WriteModeAdd,
			},
		},
	}

	_, err = s.client.Upload(input, r)
	if err != nil {
		return err
	}

	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	defer func() {
		err = s.formatError("stat", err, path)
	}()

	rp := s.getAbsPath(path)

	input := &files.GetMetadataArg{
		Path: rp,
	}

	output, err := s.client.GetMetadata(input)
	if err != nil {
		return nil, err
	}

	switch meta := output.(type) {
	case *files.FileMetadata:
		o := &types.Object{
			ID:         meta.Id,
			Type:       types.ObjectTypeFile,
			Name:       filepath.Join(path, meta.Name),
			Size:       int64(meta.Size),
			UpdatedAt:  meta.ServerModified,
			ObjectMeta: metadata.NewObjectMeta(),
		}

		if meta.ContentHash != "" {
			o.SetETag(meta.ContentHash)
		}

		return o, nil
	case *files.FolderMetadata:
		o := &types.Object{
			ID:         meta.Id,
			Type:       types.ObjectTypeDir,
			Name:       filepath.Join(path, meta.Name),
			ObjectMeta: metadata.NewObjectMeta(),
		}

		return o, nil
	default:
		o := &types.Object{
			Type: types.ObjectTypeInvalid,
		}

		return o, nil
	}
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	defer func() {
		err = s.formatError("delete", err, path)
	}()

	rp := s.getAbsPath(path)

	input := &files.DeleteArg{
		Path: rp,
	}

	_, err = s.client.DeleteV2(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) getAbsPath(path string) string {
	return strings.TrimPrefix(s.workDir+"/"+path, "/")
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	if s.loose && errors.Is(err, services.ErrCapabilityInsufficient) {
		return nil
	}

	return &services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}
