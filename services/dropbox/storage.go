package dropbox

import (
	"context"
	"fmt"
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func (s *Storage) commitAppend(ctx context.Context, o *types.Object, opt pairStorageCommitAppend) (err error) {
	rp := o.GetID()

	offset, _ := o.GetAppendOffset()

	sessionId := GetObjectSystemMetadata(o).UploadSessionID

	cursor := &files.UploadSessionCursor{
		SessionId: sessionId,
		Offset:    uint64(offset),
	}

	input := &files.CommitInfo{
		Path: rp,
		Mode: &files.WriteMode{
			Tagged: dropbox.Tagged{
				Tag: files.WriteModeOverwrite,
			},
		},
	}

	finishArg := &files.UploadSessionFinishArg{
		Cursor: cursor,
		Commit: input,
	}

	fileMetadata, err := s.client.UploadSessionFinish(finishArg, nil)
	if err != nil {
		return
	}

	o.Mode &= ^types.ModeAppend
	if fileMetadata != nil && fileMetadata.IsDownloadable {
		o.Mode |= types.ModeRead
	}

	return nil
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(true)
		o.Mode = types.ModeDir
	} else {
		o = s.newObject(false)
		o.Mode = types.ModeRead
	}

	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) createAppend(ctx context.Context, path string, opt pairStorageCreateAppend) (o *types.Object, err error) {
	startArg := &files.UploadSessionStartArg{
		Close: false,
		SessionType: &files.UploadSessionType{
			Tagged: dropbox.Tagged{
				Tag: files.UploadSessionTypeSequential,
			},
		},
	}

	res, err := s.client.UploadSessionStart(startArg, nil)
	if err != nil {
		return
	}

	sm := ObjectSystemMetadata{
		UploadSessionID: res.SessionId,
	}

	o = s.newObject(true)
	o.Mode = types.ModeAppend
	o.ID = s.getAbsPath(path)
	o.Path = path
	o.SetAppendOffset(0)
	o.SetSystemMetadata(sm)
	return o, nil
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	res, err := s.client.CreateFolderV2(&files.CreateFolderArg{
		Path: rp,
	})
	if err != nil && checkError(err, files.CreateFolderErrorPath, files.WriteErrorConflict, files.WriteConflictErrorFolder) {
		// Omit `path/conflict/folder` (`dir exists` related) error here.
		err = nil
	}
	if err != nil {
		return
	}

	if res != nil {
		// A successful response indicates that the folder is created and the returned `res` is the corresponding `FolderMetadata`.
		o = s.newObject(true)
		o.Mode = types.ModeDir
		o.ID = res.Metadata.Id
		o.Path = path
	} else {
		// `res` is nil when the given path is an existing folder.
		o = s.newObject(false)
		o.Mode = types.ModeDir
		o.ID = rp
		o.Path = path
	}

	return o, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	rp := s.getAbsPath(path)

	input := &files.DeleteArg{
		Path: rp,
	}

	// If the path is a folder, all its contents will be deleted too.
	// ref: https://www.dropbox.com/developers/documentation/http/documentation#files-delete
	_, err = s.client.DeleteV2(input)
	if err != nil && checkError(err, files.DeleteErrorPathLookup, files.LookupErrorNotFound) {
		// Omit `path_lookup/not_found` error here.
		// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	input := &objectPageStatus{
		limit: 200,
		path:  s.getAbsPath(path),
	}

	if opt.ListMode.IsPrefix() {
		input.recursive = true
	}

	return types.NewObjectIterator(ctx, s.nextObjectPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.WorkDir = s.workDir
	meta.Name = ""
	// set write restriction
	meta.SetWriteSizeMaximum(writeSizeMaximum)
	// set append restrictions
	meta.SetAppendTotalSizeMaximum(appendTotalSizeMaximum)
	return
}

func (s *Storage) nextObjectPage(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	var err error
	var output *files.ListFolderResult

	if input.cursor == "" {
		output, err = s.client.ListFolder(&files.ListFolderArg{
			Path: input.path,
		})
	} else {
		output, err = s.client.ListFolderContinue(&files.ListFolderContinueArg{
			Cursor: input.cursor,
		})
	}
	if err != nil {
		return err
	}

	for _, v := range output.Entries {
		var o *types.Object
		switch meta := v.(type) {
		case *files.FolderMetadata:
			o = s.formatFolderObject(meta.Name, meta)
		case *files.FileMetadata:
			o = s.formatFileObject(meta.Name, meta)
		}

		page.Data = append(page.Data, o)
	}

	if !output.HasMore {
		return types.IterateDone
	}

	input.cursor = output.Cursor
	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	rp := s.getAbsPath(path)

	input := &files.DownloadArg{
		Path: rp,
	}

	input.ExtraHeaders = make(map[string]string)
	if opt.HasOffset && !opt.HasSize {
		input.ExtraHeaders["Range"] = fmt.Sprintf("bytes=%d-", opt.Offset)
	} else if !opt.HasOffset && opt.HasSize {
		input.ExtraHeaders["Range"] = fmt.Sprintf("bytes=0-%d", opt.Size-1)
	} else if opt.HasOffset && opt.HasSize {
		input.ExtraHeaders["Range"] = fmt.Sprintf("bytes=%d-%d", opt.Offset, opt.Offset+opt.Size-1)
	}

	_, rc, err := s.client.Download(input)
	if err != nil {
		return 0, err
	}

	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	input := &files.GetMetadataArg{
		Path: rp,
	}

	output, err := s.client.GetMetadata(input)
	if err != nil {
		return nil, err
	}

	switch meta := output.(type) {
	case *files.FolderMetadata:
		o = s.formatFolderObject(path, meta)
	case *files.FileMetadata:
		o = s.formatFileObject(path, meta)
	}

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	if size > writeSizeMaximum {
		err = fmt.Errorf("size limit exceeded: %w", services.ErrRestrictionDissatisfied)
		return
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	}

	rp := s.getAbsPath(path)

	r = io.LimitReader(r, size)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	input := &files.CommitInfo{
		Path: rp,
		Mode: &files.WriteMode{
			Tagged: dropbox.Tagged{
				Tag: files.WriteModeOverwrite,
			},
		},
	}

	_, err = s.client.Upload(input, r)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (s *Storage) writeAppend(ctx context.Context, o *types.Object, r io.Reader, size int64, opt pairStorageWriteAppend) (n int64, err error) {
	sessionId := GetObjectSystemMetadata(o).UploadSessionID

	offset := o.MustGetAppendOffset()

	cursor := &files.UploadSessionCursor{
		SessionId: sessionId,
		Offset:    uint64(offset),
	}

	appendArg := &files.UploadSessionAppendArg{
		Cursor: cursor,
		Close:  false,
	}

	err = s.client.UploadSessionAppendV2(appendArg, r)
	if err != nil {
		return
	}

	offset += size
	o.SetAppendOffset(offset)

	return size, nil
}
