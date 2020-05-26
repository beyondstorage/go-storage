package dropbox

import (
	"context"
	"io"
	"path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/info"
)

func (s *Storage) delete(ctx context.Context, path string, opt *pairStorageDelete) (err error) {
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
func (s *Storage) listDir(ctx context.Context, dir string, opt *pairStorageListDir) (err error) {
	rp := s.getAbsPath(dir)

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
					Name:       filepath.Join(dir, meta.Name),
					Size:       int64(meta.Size),
					UpdatedAt:  meta.ServerModified,
					ObjectMeta: info.NewObjectMeta(),
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
					Name:       filepath.Join(dir, meta.Name),
					ObjectMeta: info.NewObjectMeta(),
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
func (s *Storage) metadata(ctx context.Context, opt *pairStorageMetadata) (meta info.StorageMeta, err error) {
	meta = info.NewStorageMeta()
	meta.WorkDir = s.workDir
	meta.Name = ""

	return
}
func (s *Storage) read(ctx context.Context, path string, opt *pairStorageRead) (rc io.ReadCloser, err error) {
	rp := s.getAbsPath(path)

	input := &files.DownloadArg{
		Path: rp,
	}

	_, rc, err = s.client.Download(input)
	if err != nil {
		return nil, err
	}

	if opt.HasSize {
		rc = iowrap.LimitReadCloser(rc, opt.Size)
	}

	if opt.HasReadCallbackFunc {
		rc = iowrap.CallbackReadCloser(rc, opt.ReadCallbackFunc)
	}

	return
}
func (s *Storage) stat(ctx context.Context, path string, opt *pairStorageStat) (o *types.Object, err error) {
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
			ObjectMeta: info.NewObjectMeta(),
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
			ObjectMeta: info.NewObjectMeta(),
		}

		return o, nil
	default:
		o := &types.Object{
			Type: types.ObjectTypeInvalid,
		}

		return o, nil
	}
}
func (s *Storage) write(ctx context.Context, path string, r io.Reader, opt *pairStorageWrite) (err error) {
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
