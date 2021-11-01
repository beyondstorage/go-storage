package qingstor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pengsrc/go-shared/convert"
	qerror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/stretchr/testify/assert"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func TestStorage_String(t *testing.T) {
	bucketName := "test_bucket"
	zone := "test_zone"
	c := Storage{
		workDir: "/test",
		properties: &service.Properties{
			BucketName: &bucketName,
			Zone:       &zone,
		},
	}
	assert.NotEmpty(t, c.String())
}

func TestStorage_Metadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	{
		name := uuid.New().String()
		location := uuid.New().String()

		client := Storage{
			bucket: mockBucket,
			properties: &service.Properties{
				BucketName: &name,
				Zone:       &location,
			},
		}

		m := client.Metadata()
		assert.NotNil(t, m)
		assert.Equal(t, name, m.Name)
		assert.Equal(t, location, m.MustGetLocation())
	}
}

func TestStorage_Copy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	name := uuid.New().String()
	location := uuid.New().String()

	tests := []struct {
		name     string
		src      string
		dst      string
		mockFn   func(context.Context, string, *service.PutObjectInput)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src", "test_dst",
			func(ctx context.Context, inputObjectKey string, input *service.PutObjectInput) {
				assert.Equal(t, "test_dst", inputObjectKey)
				assert.Equal(t, "/"+name+"/"+"test_src", *input.XQSCopySource)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).Do(v.mockFn)

		client := Storage{
			bucket: mockBucket,
			properties: &service.Properties{
				BucketName: &name,
				Zone:       &location,
			},
		}

		err := client.Copy(v.src, v.dst)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStorage_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		src      string
		mockFn   func(context.Context, string)
		hasError bool
		wantErr  error
	}{
		{
			"valid delete",
			"test_src",
			func(ctx context.Context, inputObjectKey string) {
				assert.Equal(t, "test_src", inputObjectKey)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().DeleteObjectWithContext(gomock.Eq(context.Background()), gomock.Any()).Do(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		err := client.Delete(v.src)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStorage_ListPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	path := uuid.New().String()
	key := uuid.New().String()

	mockBucket.EXPECT().ListObjectsWithContext(gomock.Eq(context.Background()), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *service.ListObjectsInput) (*service.ListObjectsOutput, error) {
			assert.Equal(t, path, *input.Prefix)
			assert.Equal(t, 200, *input.Limit)
			return &service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key: service.String(key),
					},
				},
			}, nil
		})

	client := Storage{
		bucket: mockBucket,
	}

	it, err := client.List(path, pairs.WithListMode(types.ListModePrefix))
	if err != nil {
		t.Error(err)
	}
	object, err := it.Next()
	assert.Equal(t, object.ID, key)
	assert.Nil(t, err)
}

func TestStorage_Move(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	name := uuid.New().String()
	location := uuid.New().String()

	tests := []struct {
		name     string
		src      string
		dst      string
		mockFn   func(context.Context, string, *service.PutObjectInput)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src", "test_dst",
			func(ctx context.Context, inputObjectKey string, input *service.PutObjectInput) {
				assert.Equal(t, "test_dst", inputObjectKey)
				assert.Equal(t, "/"+name+"/"+"test_src", *input.XQSMoveSource)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).Do(v.mockFn)

		client := Storage{
			bucket: mockBucket,
			properties: &service.Properties{
				BucketName: &name,
				Zone:       &location,
			},
		}

		err := client.Move(v.src, v.dst)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStorage_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		path     string
		mockFn   func(context.Context, string, *service.GetObjectInput) (*service.GetObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src",
			func(ctx context.Context, inputPath string, input *service.GetObjectInput) (*service.GetObjectOutput, error) {
				assert.Equal(t, "test_src", inputPath)
				return &service.GetObjectOutput{
					Body: io.NopCloser(bytes.NewBuffer([]byte("content"))),
				}, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().GetObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		var buf bytes.Buffer
		n, err := client.Read(v.path, &buf)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.Equal(t, "content", buf.String())
			assert.Equal(t, int64(buf.Len()), n)
		}
	}
}

func TestStorage_Stat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		src      string
		mockFn   func(context.Context, string, *service.HeadObjectInput) (*service.HeadObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid file",
			"test_src",
			func(ctx context.Context, objectKey string, input *service.HeadObjectInput) (*service.HeadObjectOutput, error) {
				assert.Equal(t, "test_src", objectKey)
				length := int64(100)
				return &service.HeadObjectOutput{
					ContentLength:   &length,
					ContentType:     convert.String("test_content_type"),
					ETag:            convert.String("test_etag"),
					XQSStorageClass: convert.String("STANDARD"),
				}, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().HeadObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		o, err := client.Stat(v.src)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, o)
			assert.Equal(t, types.ModeRead, o.Mode)
			assert.Equal(t, int64(100), o.MustGetContentLength())
			contentType, ok := o.GetContentType()
			assert.True(t, ok)
			assert.Equal(t, "test_content_type", contentType)
			checkSum, ok := o.GetEtag()
			assert.True(t, ok)
			assert.Equal(t, "test_etag", checkSum)

			om := GetObjectSystemMetadata(o)
			assert.Equal(t, StorageClassStandard, om.StorageClass)
		}
	}
}

func TestStorage_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		path     string
		size     int64
		r        io.Reader
		mockFn   func(context.Context, string, *service.PutObjectInput) (*service.PutObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src",
			100,
			io.LimitReader(randbytes.NewRand(), 100),
			func(ctx context.Context, inputPath string, input *service.PutObjectInput) (*service.PutObjectOutput, error) {
				assert.Equal(t, "test_src", inputPath)
				return nil, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		n, err := client.Write(v.path, v.r, v.size)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
			assert.Equal(t, v.size, n)
		}
	}
}

func TestStorage_formatError(t *testing.T) {
	s := &Storage{}
	errCasual := errors.New("casual error")
	cases := []struct {
		name      string
		op        string
		err       error
		path      string
		targetErr error
		targetEq  bool
	}{
		{
			name: "nil error",
			err:  nil,
		},
		{
			name:      "casual error",
			op:        "stat",
			err:       errCasual,
			targetErr: services.ErrUnexpected,
			targetEq:  true,
		},
		{
			name: "not found with blank code",
			op:   "stat",
			err: &qerror.QingStorError{
				StatusCode: 404,
				Code:       "",
				Message:    "msg",
			},
			targetErr: services.ErrObjectNotExist,
			targetEq:  true,
		},
		{
			name: "not found by code",
			op:   "stat",
			err: &qerror.QingStorError{
				StatusCode: 404,
				Code:       "object_not_exists",
				Message:    "msg",
			},
			targetErr: services.ErrObjectNotExist,
			targetEq:  true,
		},
		{
			name: "permission denied by code",
			op:   "stat",
			err: &qerror.QingStorError{
				StatusCode: 403,
				Code:       "permission_denied",
				Message:    "msg",
			},
			targetErr: services.ErrPermissionDenied,
			targetEq:  true,
		},
		{
			name: "not found by code error not eq",
			op:   "stat",
			err: qerror.QingStorError{
				StatusCode: 404,
				Code:       "object_not_exists",
				Message:    "msg",
			},
			targetErr: services.ErrObjectNotExist,
			targetEq:  false,
		},
	}
	for _, tt := range cases {
		err := s.formatError(tt.op, tt.err, tt.path)
		if tt.err == nil {
			assert.Nil(t, err, tt.name)
			continue
		}

		var storageErr services.StorageError
		assert.True(t, errors.As(err, &storageErr), tt.name)
		assert.Equal(t, tt.op, storageErr.Op, tt.name)

		assert.Equal(t, tt.targetEq, errors.Is(err, tt.targetErr), tt.name)
	}
}

func TestStorage_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	{
		client := Storage{
			bucket: mockBucket,
		}

		name := uuid.New().String()
		url := uuid.New().String()

		mockBucket.EXPECT().PutObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, objectKey string, input *service.PutObjectInput) (*service.PutObjectOutput, error) {
				assert.Equal(t, name, objectKey)
				assert.Equal(t, *input.XQSFetchSource, url)
				return &service.PutObjectOutput{}, nil
			})
		err := client.Fetch(name, url)
		assert.NoError(t, err)
	}

	{
		client := Storage{
			bucket: mockBucket,
		}

		name := uuid.New().String()
		url := uuid.New().String()

		mockBucket.EXPECT().PutObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, objectKey string, input *service.PutObjectInput) (*service.PutObjectOutput, error) {
				return nil, &qerror.QingStorError{}
			})
		err := client.Fetch(name, url)
		assert.Error(t, err)
	}
}

func TestStorage_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	wd := "/test/"
	c := Storage{
		workDir: wd,
		bucket:  mockBucket,
	}

	cases := []struct {
		name        string
		path        string
		multipartID string
	}{
		{
			name:        "normal object",
			path:        uuid.NewString(),
			multipartID: "",
		},
		{
			name:        "multipart object",
			path:        uuid.NewString(),
			multipartID: uuid.NewString(),
		},
	}

	mockBucket.EXPECT().HeadObjectWithContext(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, objectKey string, input *service.HeadObjectInput) (*service.HeadObjectOutput, error) {
			return &service.HeadObjectOutput{}, nil
		}).Times(1)

	for _, tt := range cases {
		ps := make([]types.Pair, 0)
		if tt.multipartID != "" {
			ps = append(ps, pairs.WithMultipartID(tt.multipartID))
		}
		obj := c.Create(tt.path, ps...)
		assert.NotNil(t, obj)
		assert.Equal(t, c.getAbsPath(tt.path), obj.ID)
		assert.Equal(t, tt.path, obj.Path)
		if tt.multipartID != "" {
			assert.Equal(t, tt.multipartID, obj.MustGetMultipartID())
			assert.Equal(t, types.ModePart, obj.Mode)
		} else {
			assert.Equal(t, types.ModeRead, obj.Mode)
			assert.Panics(t, func() {
				obj.MustGetMultipartID()
			})
		}
	}
}
