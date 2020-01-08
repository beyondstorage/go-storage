package uss

import (
	"fmt"
	"io"
	"strings"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/metadata"
	"github.com/upyun/go-sdk/upyun"
)

// Storage is the uss service.
//
//go:generate ../../internal/bin/meta
//go:generate ../../internal/bin/context
type Storage struct {
	bucket *upyun.UpYun

	name    string
	workDir string
}

// New will create a new uss service.
func New(name string, pairs ...*types.Pair) (s *Storage, err error) {
	const errorMessage = "%s New: %w"

	s = &Storage{}

	opt, err := parseStoragePairNew(pairs...)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, err)
	}

	credProtocol, cred := opt.Credential.Protocol(), opt.Credential.Value()
	if credProtocol != credential.ProtocolHmac {
		return nil, fmt.Errorf(errorMessage, s, credential.ErrUnsupportedProtocol)
	}

	cfg := &upyun.UpYunConfig{
		Bucket:   name,
		Operator: cred[0],
		Password: cred[1],
	}
	s.bucket = upyun.NewUpYun(cfg)
	s.name = name
	return
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager uss {Name: %s, WorkDir: %s}",
		s.name, "/"+s.workDir)
}

// Init implements Storager.Init
func (s *Storage) Init(pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Init: %w"

	opt, err := parseStoragePairInit(pairs...)
	if err != nil {
		return fmt.Errorf(errorMessage, s, err)
	}

	if opt.HasWorkDir {
		// TODO: we should validate workDir
		s.workDir = strings.TrimLeft(opt.WorkDir, "/")
	}

	return nil
}

// Metadata implements Storager.Metadata
func (s *Storage) Metadata(pairs ...*types.Pair) (m metadata.StorageMeta, err error) {
	m = metadata.NewStorageMeta()
	m.Name = s.name
	m.WorkDir = s.workDir
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

	ch := make(chan *upyun.FileInfo, 200)
	defer close(ch)

	go func() {
		for v := range ch {
			o := &types.Object{
				ID:         v.Name,
				Name:       s.getRelPath(v.Name),
				Type:       types.ObjectTypeFile,
				Size:       v.Size,
				UpdatedAt:  v.Time,
				ObjectMeta: metadata.NewObjectMeta(),
			}
			o.SetETag(v.ETag)

			opt.FileFunc(o)
		}
	}()

	err = s.bucket.List(&upyun.GetObjectsConfig{
		Path:        rp,
		ObjectsChan: ch,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return
}

// Read implements Storager.Read
func (s *Storage) Read(path string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
	const errorMessage = "%s Read [%s]: %w"

	rp := s.getAbsPath(path)

	r, w := io.Pipe()

	_, err = s.bucket.Get(&upyun.GetObjectConfig{
		Path:   rp,
		Writer: w,
	})
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}
	return r, nil
}

// Write implements Storager.Write
func (s *Storage) Write(path string, r io.Reader, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Write [%s]: %w"

	rp := s.getAbsPath(path)

	cfg := &upyun.PutObjectConfig{
		Path:   rp,
		Reader: r,
	}

	err = s.bucket.Put(cfg)
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return
}

// Stat implements Storager.Stat
func (s *Storage) Stat(path string, pairs ...*types.Pair) (o *types.Object, err error) {
	const errorMessage = "%s Stat [%s]: %w"

	rp := s.getAbsPath(path)

	output, err := s.bucket.GetInfo(rp)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, s, path, err)
	}

	o = &types.Object{
		ID:         rp,
		Name:       path,
		Type:       types.ObjectTypeFile,
		Size:       output.Size,
		UpdatedAt:  output.Time,
		ObjectMeta: metadata.NewObjectMeta(),
	}
	return o, nil
}

// Delete implements Storager.Delete
func (s *Storage) Delete(path string, pairs ...*types.Pair) (err error) {
	const errorMessage = "%s Delete [%s]: %w"

	rp := s.getAbsPath(path)

	err = s.bucket.Delete(&upyun.DeleteObjectConfig{
		Path: rp,
	})
	if err != nil {
		return fmt.Errorf(errorMessage, s, path, err)
	}
	return
}
