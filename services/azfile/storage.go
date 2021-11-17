package azfile

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Azure/azure-storage-file-go/azfile"

	"go.beyondstorage.io/v5/pkg/iowrap"
	"go.beyondstorage.io/v5/types"
)

func (s *Storage) create(path string, opt pairStorageCreate) (o *types.Object) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o = s.newObject(true)
		o.Mode |= types.ModeDir
	} else {
		o = s.newObject(false)
		o.Mode |= types.ModeRead
	}

	o.ID = rp
	o.Path = path

	return o
}

func (s *Storage) createDir(ctx context.Context, path string, opt pairStorageCreateDir) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	attribute := azfile.FileAttributeNone

	properties := azfile.SMBProperties{
		FileAttributes: &attribute,
	}

	fi, err := s.client.NewDirectoryURL(path).GetProperties(ctx)
	if err == nil {
		// The directory exist, we should set the metadata.
		o = s.newObject(true)
		o.SetLastModified(fi.LastModified())
	} else if !checkError(err, fileNotFound) {
		// Something error other then file not found happened, return directly.
		return nil, err
	} else {
		// The directory not exists, we should create the directory.
		_, err = s.client.NewDirectoryURL(s.getRelativePath(path)).Create(ctx, azfile.Metadata{}, properties)
		if err != nil {
			return nil, err
		}

		o = s.newObject(false)
	}

	o.ID = rp
	o.Path = path
	o.Mode |= types.ModeDir

	return
}

func (s *Storage) delete(ctx context.Context, path string, opt pairStorageDelete) (err error) {
	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		_, err = s.client.NewDirectoryURL(s.getRelativePath(path)).Delete(ctx)
	} else {
		_, err = s.client.NewFileURL(s.getRelativePath(path)).Delete(ctx)
	}

	if err != nil {
		// azfile Delete is not idempotent, so we need to check file not found error.
		//
		// References
		// - [GSP-46](https://github.com/beyondstorage/specs/blob/master/rfcs/46-idempotent-delete.md)
		// - https://docs.microsoft.com/en-us/rest/api/storageservices/delete-file2#remarks
		if checkError(err, fileNotFound) {
			err = nil
		} else {
			return err
		}
	}

	return nil
}

func (s *Storage) list(ctx context.Context, path string, opt pairStorageList) (oi *types.ObjectIterator, err error) {
	input := &objectPageStatus{
		maxResults: 200,
		prefix:     s.getRelativePath(path),
	}

	return types.NewObjectIterator(ctx, s.nextObjectPage, input), nil
}

func (s *Storage) metadata(opt pairStorageMetadata) (meta *types.StorageMeta) {
	meta = types.NewStorageMeta()
	meta.WorkDir = s.workDir
	return meta
}

func (s *Storage) nextObjectPage(ctx context.Context, page *types.ObjectPage) error {
	input := page.Status.(*objectPageStatus)

	options := azfile.ListFilesAndDirectoriesOptions{
		MaxResults: input.maxResults,
		Prefix:     input.prefix,
	}

	fmt.Println("option prefix: " + options.Prefix)

	output, err := s.client.ListFilesAndDirectoriesSegment(ctx, input.marker, options)
	if err != nil {
		return err
	}
	fmt.Print("Len of FileItems:")
	fmt.Println(len(output.FileItems))

	fmt.Print("Len of DirectoryItems:")
	fmt.Println(len(output.DirectoryItems))

	for _, v := range output.DirectoryItems {
		o, err := s.formatDirObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	for _, v := range output.FileItems {
		o, err := s.formatFileObject(v)
		if err != nil {
			return err
		}

		page.Data = append(page.Data, o)
	}

	if !output.NextMarker.NotDone() {
		return types.IterateDone
	}

	input.marker = output.NextMarker

	return nil
}

func (s *Storage) read(ctx context.Context, path string, w io.Writer, opt pairStorageRead) (n int64, err error) {
	offset := int64(0)
	if opt.HasOffset {
		offset = opt.Offset
	}

	count := int64(azfile.CountToEnd)
	if opt.HasSize {
		count = opt.Size
	}

	output, err := s.client.NewFileURL(s.getRelativePath(path)).Download(ctx, offset, count, false)
	if err != nil {
		return 0, err
	}
	defer func() {
		cErr := output.Response().Body.Close()
		if cErr != nil {
			err = cErr
		}
	}()

	rc := output.Response().Body
	if opt.HasIoCallback {
		rc = iowrap.CallbackReadCloser(rc, opt.IoCallback)
	}

	return io.Copy(w, rc)
}

func (s *Storage) stat(ctx context.Context, path string, opt pairStorageStat) (o *types.Object, err error) {
	rp := s.getAbsPath(path)

	var dirOutput *azfile.DirectoryGetPropertiesResponse
	var fileOutput *azfile.FileGetPropertiesResponse

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		dirOutput, err = s.client.NewDirectoryURL(s.getRelativePath(path)).GetProperties(ctx)
	} else {
		fileOutput, err = s.client.NewFileURL(s.getRelativePath(path)).GetProperties(ctx)
	}

	if err != nil {
		return nil, err
	}

	o = s.newObject(true)
	o.ID = rp
	o.Path = path

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		o.Mode |= types.ModeDir

		o.SetLastModified(dirOutput.LastModified())

		if v := string(dirOutput.ETag()); v != "" {
			o.SetEtag(v)
		}

		var sm ObjectSystemMetadata
		if v, err := strconv.ParseBool(dirOutput.IsServerEncrypted()); err == nil {
			sm.ServerEncrypted = v
		}
		o.SetSystemMetadata(sm)
	} else {
		o.Mode |= types.ModeRead

		o.SetContentLength(fileOutput.ContentLength())
		o.SetLastModified(fileOutput.LastModified())

		if v := string(fileOutput.ETag()); v != "" {
			o.SetEtag(v)
		}
		if v := fileOutput.ContentType(); v != "" {
			o.SetContentType(v)
		}
		if v := fileOutput.ContentMD5(); len(v) > 0 {
			o.SetContentMd5(base64.StdEncoding.EncodeToString(v))
		}

		var sm ObjectSystemMetadata
		if v, err := strconv.ParseBool(fileOutput.IsServerEncrypted()); err == nil {
			sm.ServerEncrypted = v
		}
		o.SetSystemMetadata(sm)
	}

	return o, nil
}

func (s *Storage) write(ctx context.Context, path string, r io.Reader, size int64, opt pairStorageWrite) (n int64, err error) {
	relativePath := s.getRelativePath(path)
	err = s.mkDirs(ctx, filepath.Dir(relativePath))
	if err != nil {
		return
	}

	if opt.HasIoCallback {
		r = iowrap.CallbackReader(r, opt.IoCallback)
	}

	// According to GSP-751, we should allow the user to pass in a nil io.Reader.
	// ref: https://github.com/beyondstorage/go-storage/blob/master/docs/rfcs/751-write-empty-file-behavior.md
	if r == nil && size == 0 {
		r = strings.NewReader("")
	} else if r == nil && size != 0 {
		return 0, fmt.Errorf("reader is nil but size is not 0")
	} else {
		r = io.LimitReader(r, size)
	}

	headers := azfile.FileHTTPHeaders{}

	if opt.HasContentType {
		headers.ContentType = opt.ContentType
	}

	// `Create` only initializes the file.
	// ref: https://docs.microsoft.com/en-us/rest/api/storageservices/create-file
	_, err = s.client.NewFileURL(relativePath).Create(ctx, size, headers, azfile.Metadata{})
	if err != nil {
		return 0, err
	}

	if size > 0 {
		body := iowrap.SizedReadSeekCloser(r, size)

		var transactionalMD5 []byte
		if opt.HasContentMd5 {
			transactionalMD5, err = base64.StdEncoding.DecodeString(opt.ContentMd5)
			if err != nil {
				return 0, err
			}
		}

		// Since `Create' only initializes the file, we need to call `UploadRange' to write the contents to the file.
		_, err = s.client.NewFileURL(relativePath).UploadRange(ctx, 0, body, transactionalMD5)
		if err != nil {
			return 0, err
		}
	}

	return size, nil
}
