package qingstor

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pengsrc/go-shared/convert"
	qerror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/stretchr/testify/assert"

	"github.com/aos-dev/go-storage/v2/pkg/segment"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/info"
	"github.com/aos-dev/go-storage/v2/types/pairs"
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

		m, err := client.Metadata()
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, name, m.Name)
		assert.Equal(t, location, m.MustGetLocation())
	}
}

func TestStorage_Statistical(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	{
		client := Storage{
			bucket: mockBucket,
		}

		name := uuid.New().String()
		location := uuid.New().String()
		size := int64(1234)
		count := int64(4321)

		mockBucket.EXPECT().GetStatistics().DoAndReturn(func() (*service.GetBucketStatisticsOutput, error) {
			return &service.GetBucketStatisticsOutput{
				Name:     &name,
				Location: &location,
				Size:     &size,
				Count:    &count,
			}, nil
		})
		m, err := client.Statistical()
		assert.NoError(t, err)
		assert.NotNil(t, m)
	}

	{
		client := Storage{
			bucket: mockBucket,
		}

		mockBucket.EXPECT().GetStatistics().DoAndReturn(func() (*service.GetBucketStatisticsOutput, error) {
			return nil, &qerror.QingStorError{}
		})
		_, err := client.Statistical()
		assert.Error(t, err)
	}
}

func TestStorage_AbortSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	client := Storage{
		bucket: mockBucket,
	}

	// Test valid segment.
	path := uuid.New().String()
	id := uuid.New().String()
	seg := segment.NewIndexBasedSegment(path, id)
	mockBucket.EXPECT().AbortMultipartUpload(gomock.Any(), gomock.Any()).Do(func(inputPath string, input *service.AbortMultipartUploadInput) {
		assert.Equal(t, path, inputPath)
		assert.Equal(t, id, *input.UploadID)
	})
	err := client.AbortSegment(seg)
	assert.NoError(t, err)
}

func TestStorage_CompleteSegment(t *testing.T) {
	t.Run("valid segment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBucket := NewMockBucket(ctrl)

		path, id := uuid.New().String(), uuid.New().String()
		seg := segment.NewIndexBasedSegment(path, id)

		client := Storage{
			bucket: mockBucket,
		}

		mockBucket.EXPECT().CompleteMultipartUpload(gomock.Any(), gomock.Any()).Do(func(inputPath string, input *service.CompleteMultipartUploadInput) {
			assert.Equal(t, path, inputPath)
			assert.Equal(t, id, *input.UploadID)
		})

		err := client.CompleteSegment(seg)
		assert.NoError(t, err)
	})
}

func TestStorage_Copy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		src      string
		dst      string
		mockFn   func(string, *service.PutObjectInput)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src", "test_dst",
			func(inputObjectKey string, input *service.PutObjectInput) {
				assert.Equal(t, "test_dst", inputObjectKey)
				assert.Equal(t, "test_src", *input.XQSCopySource)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObject(gomock.Any(), gomock.Any()).Do(v.mockFn)

		client := Storage{
			bucket: mockBucket,
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
		mockFn   func(string)
		hasError bool
		wantErr  error
	}{
		{
			"valid delete",
			"test_src",
			func(inputObjectKey string) {
				assert.Equal(t, "test_src", inputObjectKey)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().DeleteObject(gomock.Any()).Do(v.mockFn)

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

func TestStorage_InitSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		path     string
		segments map[string]*segment.IndexBasedSegment
		hasCall  bool
		mockFn   func(string, *service.InitiateMultipartUploadInput) (*service.InitiateMultipartUploadOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid init segment",
			"test", map[string]*segment.IndexBasedSegment{},
			true,
			func(inputPath string, input *service.InitiateMultipartUploadInput) (*service.InitiateMultipartUploadOutput, error) {
				assert.Equal(t, "test", inputPath)

				uploadID := "test"
				return &service.InitiateMultipartUploadOutput{
					UploadID: &uploadID,
				}, nil
			},
			false, nil,
		},
		{
			"segment already exist",
			"test",
			map[string]*segment.IndexBasedSegment{
				"test": segment.NewIndexBasedSegment("xxx", "test_segment_id"),
			},
			true,
			func(inputPath string, input *service.InitiateMultipartUploadInput) (*service.InitiateMultipartUploadOutput, error) {
				assert.Equal(t, "test", inputPath)

				uploadID := "test"
				return &service.InitiateMultipartUploadOutput{
					UploadID: &uploadID,
				}, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		if v.hasCall {
			mockBucket.EXPECT().InitiateMultipartUpload(gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)
		}

		client := Storage{
			bucket: mockBucket,
		}

		_, err := client.InitIndexSegment(v.path)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStorage_ListDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	keys := make([]string, 7)
	for idx := range keys {
		keys[idx] = uuid.New().String()
	}

	tests := []struct {
		name   string
		path   string
		output *service.ListObjectsOutput
		items  []*types.Object
		err    error
	}{
		{
			"list without delimiter",
			uuid.New().String(),
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key: service.String(keys[0]),
					},
				},
			},
			[]*types.Object{
				{
					ID:         keys[0],
					Name:       keys[0],
					Type:       types.ObjectTypeFile,
					ObjectMeta: info.NewObjectMeta(),
				},
			},
			nil,
		},
		{
			"list with return next marker",
			uuid.New().String(),
			&service.ListObjectsOutput{
				NextMarker: service.String("test_marker"),
				HasMore:    service.Bool(false),
				CommonPrefixes: []*string{
					service.String(keys[2]),
				},
			},
			[]*types.Object{
				{
					ID:         keys[2],
					Name:       keys[2],
					Type:       types.ObjectTypeDir,
					ObjectMeta: info.NewObjectMeta(),
				},
			},
			nil,
		},
		{
			"list with return empty keys",
			uuid.New().String(),
			&service.ListObjectsOutput{
				NextMarker: service.String("test_marker"),
				HasMore:    service.Bool(true),
			},
			[]*types.Object{},
			nil,
		},
		{
			"list with error return",
			uuid.New().String(),
			nil,
			[]*types.Object{},
			&qerror.QingStorError{
				StatusCode: 401,
			},
		},
		{
			"list with all data returned",
			uuid.New().String(),
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key:      service.String(keys[5]),
						MimeType: service.String("application/json"),
						Etag:     service.String("xxxxx"),
						Size:     service.Int64(1233),
						Modified: service.Int(1233),
					},
				},
			},
			[]*types.Object{
				{
					ID:        keys[5],
					Name:      keys[5],
					Type:      types.ObjectTypeFile,
					Size:      1233,
					UpdatedAt: time.Unix(1233, 0),
					ObjectMeta: info.NewObjectMeta().
						SetContentType("application/json").
						SetETag("xxxxx"),
				},
			},
			nil,
		},
		{
			"list with dir and objects",
			keys[6],
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key:      service.String(keys[6]),
						MimeType: service.String("application/x-directory"),
						Etag:     service.String("xxxxx"),
						Size:     service.Int64(0),
						Modified: service.Int(1233),
					},
				},
			},
			[]*types.Object{},
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := v.path

			mockBucket.EXPECT().ListObjects(gomock.Any()).DoAndReturn(func(input *service.ListObjectsInput) (*service.ListObjectsOutput, error) {
				assert.Equal(t, path, *input.Prefix)
				assert.Equal(t, 200, *input.Limit)
				return v.output, v.err
			})

			client := Storage{
				bucket: mockBucket,
			}

			items := make([]*types.Object, 0)

			err := client.ListDir(path, pairs.WithDirFunc(func(object *types.Object) {
				items = append(items, object)
			}), pairs.WithFileFunc(func(object *types.Object) {
				items = append(items, object)
			}))
			assert.Equal(t, v.err == nil, err == nil)
			assert.EqualValues(t, v.items, items)
		})
	}
}

func TestStorage_ListPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	path := uuid.New().String()
	key := uuid.New().String()

	mockBucket.EXPECT().ListObjects(gomock.Any()).DoAndReturn(func(input *service.ListObjectsInput) (*service.ListObjectsOutput, error) {
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

	err := client.ListPrefix(path, pairs.WithObjectFunc(func(object *types.Object) {
		assert.Equal(t, object.ID, key)
	}))
	assert.Nil(t, err)
}

func TestStorage_Move(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		src      string
		dst      string
		mockFn   func(string, *service.PutObjectInput)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src", "test_dst",
			func(inputObjectKey string, input *service.PutObjectInput) {
				assert.Equal(t, "test_dst", inputObjectKey)
				assert.Equal(t, "test_src", *input.XQSMoveSource)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObject(gomock.Any(), gomock.Any()).Do(v.mockFn)

		client := Storage{
			bucket: mockBucket,
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
		mockFn   func(string, *service.GetObjectInput) (*service.GetObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src",
			func(inputPath string, input *service.GetObjectInput) (*service.GetObjectOutput, error) {
				assert.Equal(t, "test_src", inputPath)
				return &service.GetObjectOutput{
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte("content"))),
				}, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().GetObject(gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		r, err := client.Read(v.path)
		if v.hasError {
			assert.Error(t, err)
			assert.Nil(t, r)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NotNil(t, r)
			content, rerr := ioutil.ReadAll(r)
			assert.NoError(t, rerr)
			assert.Equal(t, "content", string(content))
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
		mockFn   func(objectKey string, input *service.HeadObjectInput) (*service.HeadObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid file",
			"test_src",
			func(objectKey string, input *service.HeadObjectInput) (*service.HeadObjectOutput, error) {
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
		mockBucket.EXPECT().HeadObject(gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

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
			assert.Equal(t, types.ObjectTypeFile, o.Type)
			assert.Equal(t, int64(100), o.Size)
			contentType, ok := o.GetContentType()
			assert.True(t, ok)
			assert.Equal(t, "test_content_type", contentType)
			checkSum, ok := o.GetETag()
			assert.True(t, ok)
			assert.Equal(t, "test_etag", checkSum)
			storageClass, ok := GetStorageClass(o.ObjectMeta)
			assert.True(t, ok)
			assert.Equal(t, StorageClassStandard, storageClass)
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
		mockFn   func(string, *service.PutObjectInput) (*service.PutObjectOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid copy",
			"test_src",
			100,
			func(inputPath string, input *service.PutObjectInput) (*service.PutObjectOutput, error) {
				assert.Equal(t, "test_src", inputPath)
				return nil, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().PutObject(gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Storage{
			bucket: mockBucket,
		}

		err := client.Write(v.path, nil, pairs.WithSize(v.size))
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestStorage_WriteSegment(t *testing.T) {
	t.Run("valid segment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBucket := NewMockBucket(ctrl)

		path, id := uuid.New().String(), uuid.New().String()
		seg := segment.NewIndexBasedSegment(path, id)

		client := Storage{
			bucket: mockBucket,
		}

		mockBucket.EXPECT().UploadMultipart(gomock.Any(), gomock.Any()).Do(func(inputPath string, input *service.UploadMultipartInput) {
			assert.Equal(t, path, inputPath)
			assert.Equal(t, id, *input.UploadID)
		})

		err := client.WriteIndexSegment(seg, nil, 0, 100)
		assert.NoError(t, err)
	})
}

func TestStorage_ListPrefixSegments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	keys := make([]string, 100)
	for idx := range keys {
		keys[idx] = uuid.New().String()
	}

	tests := []struct {
		name   string
		output *service.ListMultipartUploadsOutput
		items  []segment.Segment
		err    error
	}{
		{
			"list without delimiter",
			&service.ListMultipartUploadsOutput{
				HasMore: service.Bool(false),
				Uploads: []*service.UploadsType{
					{
						Key:      service.String(keys[0]),
						UploadID: service.String(keys[1]),
					},
				},
			},
			[]segment.Segment{
				segment.NewIndexBasedSegment(keys[0], keys[1]),
			},
			nil,
		},
		{
			"list with return next marker",
			&service.ListMultipartUploadsOutput{
				NextKeyMarker: service.String("test_marker"),
				HasMore:       service.Bool(false),
				Uploads: []*service.UploadsType{
					{
						Key:      service.String(keys[1]),
						UploadID: service.String(keys[2]),
					},
				},
			},
			[]segment.Segment{
				segment.NewIndexBasedSegment(keys[1], keys[2]),
			},
			nil,
		},
		{
			"list with error return",
			nil,
			[]segment.Segment{},
			&qerror.QingStorError{
				StatusCode: 401,
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := uuid.New().String()

			mockBucket.EXPECT().ListMultipartUploads(gomock.Any()).DoAndReturn(func(input *service.ListMultipartUploadsInput) (*service.ListMultipartUploadsOutput, error) {
				assert.Equal(t, path, *input.Prefix)
				assert.Equal(t, 200, *input.Limit)
				return v.output, v.err
			})

			client := Storage{
				bucket: mockBucket,
			}

			items := make([]segment.Segment, 0)

			err := client.ListPrefixSegments(path,
				pairs.WithSegmentFunc(func(segment segment.Segment) {
					items = append(items, segment)
				}),
			)
			assert.Equal(t, v.err == nil, err == nil)
			assert.Equal(t, v.items, items)
		})
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
			targetErr: errCasual,
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

		var storageErr *services.StorageError
		assert.True(t, errors.As(err, &storageErr), tt.name)
		assert.Equal(t, tt.op, storageErr.Op, tt.name)

		assert.Equal(t, tt.targetEq, errors.Is(err, tt.targetErr), tt.name)
	}
}
