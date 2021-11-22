package gdrive

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"google.golang.org/api/drive/v3"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

const directoryMimeType = "application/vnd.google-apps.folder"

func (s *Storage) copy(ctx context.Context, src string, dst string, opt pairStorageCopy) (err error) {

	var dstFile *drive.File

	srcFileId, err := s.pathToId(ctx, src)
	if err != nil {
		return err
	}

	dstFileId, err := s.pathToId(ctx, dst)
	if err != nil {
		return err
	}

	// FIXME: I don't know how to directly copy a file into an existing one
	if dstFileId != "" {
		err = s.service.Files.Delete(dstFileId).Do()
		if err != nil {
			return err
		}
	}
	dstFile = &drive.File{
		Name: s.getFileName(dst),
	}
	_, err = s.service.Files.Copy(srcFileId, dstFile).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	o = s.newObject(false)
	o.ID = s.getAbsPath(path)
	o.Path = path
	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *types.Object, err error) {

	_, err = s.createDirs(ctx, path)

	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = s.getAbsPath(path)
	o.Path = path
	o.Mode = types.ModeDir

	return o, nil

}

// This function is very similar to `createDir` but has different uses. Unlike `creatDir`, it
// is mainly responsible for communicating with gdrive API
func (s *Storage) createDirs(ctx context.Context, path string) (parentsId string, err error) {
	pathUnits := strings.Split(path, "/")
	parentsId = "root"

	for _, v := range pathUnits {
		// TODO: use `strings.Split` to split path is not perfect, maybe
		// we should add a helper function to do this.
		if v != "" {
			parentsId, err = s.mkDir(ctx, parentsId, v)
			if err != nil {
				return "", err
			}
		}
	}

	return parentsId, nil
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	var fileId string
	fileId, err = s.pathToId(ctx, path)
	if err != nil {
		return err
	}
	err = s.service.Files.Delete(fileId).Do()

	// Omit `path_lookup/not_found` error here.
	// ref: [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
	if err != nil && strings.Contains(err.Error(), "404") {
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

	if !opt.HasListMode || opt.ListMode.IsDir() {
		return types.NewObjectIterator(ctx, s.nextObjectPage, input), nil
	} else {
		return nil, services.ListModeInvalidError{Actual: opt.ListMode}
	}
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.Name = s.name
	meta.WorkDir = s.workDir
	return meta
}

// Create a directory by passing it's name and the parents' fileId.
// It will return the fileId of the directory whether it exist or not.
// If error occurs, it will return an empty string and error.
func (s *Storage) mkDir(ctx context.Context, parents string, dirName string) (string, error) {
	id, err := s.searchContentInDir(ctx, parents, dirName)
	if err != nil {
		return "", err
	}
	// Simply return the fileId if the directory already exist
	if id != "" {
		return id, nil
	}

	// create a directory if not exist
	dir := &drive.File{
		Name:     dirName,
		Parents:  []string{parents},
		MimeType: directoryMimeType,
	}
	f, err := s.service.Files.Create(dir).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	return f.Id, nil
}

func (s *Storage) nextObjectPage(ctx context.Context, page *types.ObjectPage) (err error) {
	input := page.Status.(*objectPageStatus)

	var dirId string
	dirId, err = s.pathToId(ctx, input.path)
	if err != nil {
		return err
	}

	// When dirId is empty, the path is an empty dir. We should directly return IterateDone
	if dirId == "" {
		return types.IterateDone
	}

	q := s.service.Files.List().Q(fmt.Sprintf("parents='%s'", dirId)).Fields("*")

	if input.pageToken != "" {
		q = q.PageToken(input.pageToken)
	}
	r, err := q.Do()

	if err != nil {
		return err
	}

	for _, f := range r.Files {
		o := s.newObject(true)
		o.SetContentLength(f.Size)
		o.Path = s.getRelativePath(input.path, f.Name)
		switch f.MimeType {
		case directoryMimeType:
			o.Mode = types.ModeDir
		default:
			o.Mode = types.ModeRead
		}
		page.Data = append(page.Data, o)
	}

	if r.IncompleteSearch == false {
		return types.IterateDone
	}

	input.pageToken = r.NextPageToken
	return nil
}

// pathToId converts path to fileId, as we discussed in RFC-14.
// Ref: https://github.com/beyondstorage/go-service-gdrive/blob/master/docs/rfcs/14-gdrive-for-go-storage-design.md
// Behavior:
// err represents the error handled in pathToId
// fileId represents the results: fileId empty means the path is not exist, otherwise it's the fileId of input path
func (s *Storage) pathToId(ctx context.Context, path string) (fileId string, err error) {
	path = s.getAbsPath(path)

	fileId, found := s.getCache(path)
	if found {
		return fileId, nil
	}

	pathUnits := strings.Split(path, "/")
	fileId = "root"
	cacheCurrentPath := ""
	// Traverse the whole path, break the loop if we fails at one search
	for _, v := range pathUnits {
		fileId, err = s.searchContentInDir(ctx, fileId, v)

		if fileId == "" || err != nil {
			break
		}

		if cacheCurrentPath == "" {
			cacheCurrentPath = v
		} else {
			cacheCurrentPath += "/" + v
		}

		s.setCache(cacheCurrentPath, fileId)
	}

	if err != nil {
		return "", err
	}

	return fileId, nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	fileId, err := s.pathToId(ctx, path)
	if err != nil {
		return 0, err
	}
	fileGetCall := s.service.Files.Get(fileId)
	if opt.HasOffset && !opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=%d-", opt.Offset)
		fileGetCall.Header().Add("Range", rangeBytes)
	} else if !opt.HasOffset && opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=0-%d", opt.Size-1)
		fileGetCall.Header().Add("Range", rangeBytes)
	} else if opt.HasOffset && opt.HasSize {
		rangeBytes := fmt.Sprintf("bytes=%d-%d", opt.Offset, opt.Offset+opt.Size-1)
		fileGetCall.Header().Add("Range", rangeBytes)
	}
	f, err := fileGetCall.Context(ctx).Download()
	if err != nil {
		return 0, err
	}

	var rc io.ReadCloser
	rc = f.Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

// Search something in directory by passing it's name and the fileId of the folder.
// It will return the fileId of the content we want, and nil for sure.
// If nothing is found, we will return an empty string and nil.
// We will only return non nil if error occurs.
func (s *Storage) searchContentInDir(ctx context.Context, dirId string, contentName string) (fileId string, err error) {
	searchArg := fmt.Sprintf("name = '%s' and parents = '%s'", contentName, dirId)
	fileList, err := s.service.Files.List().Context(ctx).Q(searchArg).Fields("*").Do()
	if err != nil {
		return "", err
	}
	// Because we assume that the path is unique, so there would be only two results: One file matches or none
	if len(fileList.Files) == 0 {
		return "", nil
	}
	return fileList.Files[0].Id, nil

}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {

	content, err := s.pathToId(ctx, path)

	if content == "" {
		return nil, services.ErrObjectNotExist
	}

	if err != nil {
		return nil, err
	}

	rp := s.getAbsPath(path)
	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	//TODO: Just a temporary hack, maybe add a helper function to do this?
	file, _ := s.service.Files.Get(content).Context(ctx).Fields("*").Do()

	if file.MimeType == directoryMimeType {
		o.Mode |= types.ModeDir
	}

	o.SetContentLength(file.Size)

	return o, nil
}

// First we need make sure this file is not exist.
// If it is, then we upload it, or we will overwrite it.
func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not nil")
	}

	// Parent directory of the file
	var parentsId string

	r = io.LimitReader(r, size)

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	fileId, err := s.pathToId(ctx, path)

	if err != nil {
		return 0, err
	}

	// fileId can be empty when err is nil
	if fileId == "" {
		// upload
		dirs, fileName := filepath.Split(s.getAbsPath(path))

		if dirs != "" {
			parentsId, err = s.createDirs(ctx, dirs)
			if err != nil {
				return 0, err
			}

		}

		file := &drive.File{
			Name:    fileName,
			Parents: []string{parentsId},
		}
		_, err = s.service.Files.Create(file).Context(ctx).Media(r).Do()

		if err != nil {
			return 0, err
		}

	} else {
		// update
		newFile := &drive.File{Name: s.getFileName(path)}
		_, err = s.service.Files.Update(fileId, newFile).Context(ctx).Media(r).Do()

		if err != nil {
			return 0, err
		}
	}

	return size, nil
}
