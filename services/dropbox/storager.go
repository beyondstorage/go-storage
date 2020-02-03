package dropbox

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

	"github.com/Xuanwo/storage/pkg/iowrap"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
)

// Storage is the dropbox client.
type Storage struct {
	client files.Client

	workDir string
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

// List implements Storager.List
func (s *Storage) List(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s List [%s]: %w"

	opt, err := parseStoragePairList(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	result, err := s.client.ListFolder(&files.ListFolderArg{
		Path: rp,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
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
				// TODO: manage ContentHash

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
				return fmt.Errorf(errorMessage, s, path, ErrUnexpectedEntry)
			}
		}
		if !result.HasMore {
			break
		}

		result, err = s.client.ListFolderContinue(&files.ListFolderContinueArg{
			Cursor: result.Cursor,
		})
		if err != nil {
			return fmt.Errorf(errorMessage, s, path, err)
		}
	}
	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	opt, err := parseStoragePairRead(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	input := &files.DownloadArg{
		Path: rp,
	}

	_, r, err = s.client.Download(input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	if opt.HasSize {
		return iowrap.LimitReadCloser(r, opt.Size), nil
	}

	return
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	opt, err := parseStoragePairWrite(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	rp := s.getAbsPath(path)

	if opt.HasSize {
		r = io.LimitReader(r, opt.Size)
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
		return fmt.Errorf(errorMessage, s, path, err)
	}

	return nil
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	input := &files.GetMetadataArg{
		Path: rp,
	}

	output, err := s.client.GetMetadata(input)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
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
		// TODO: manage ContentHash

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
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	input := &files.DeleteArg{
		Path: rp,
	}

	_, err = s.client.DeleteV2(input)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}

	return nil
}
