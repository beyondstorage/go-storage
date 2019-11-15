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
	"github.com/stretchr/testify/assert"
	qerror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/pkg/iterator"

	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
)

func TestClient_Init(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		client := Client{}
		err := client.Init()
		assert.NoError(t, err)
		assert.Equal(t, "", client.workDir)
	})

	t.Run("with workDir", func(t *testing.T) {
		client := Client{}
		err := client.Init(types.WithWorkDir("test"))
		assert.NoError(t, err)
		assert.Equal(t, "test", client.workDir)
	})
}

func TestClient_String(t *testing.T) {
	c := Client{}
	err := c.Init(types.WithWorkDir("/test"))
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, c.String())
}

func TestClient_Metadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	{
		client := Client{
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
		m, err := client.Metadata()
		assert.NoError(t, err)
		assert.NotNil(t, m)
		gotName, ok := m.GetName()
		assert.True(t, ok)
		assert.Equal(t, name, gotName)
	}

	{
		client := Client{
			bucket: mockBucket,
		}

		mockBucket.EXPECT().GetStatistics().DoAndReturn(func() (*service.GetBucketStatisticsOutput, error) {
			return nil, &qerror.QingStorError{}
		})
		_, err := client.Metadata()
		assert.Error(t, err)
		assert.True(t, errors.Is(err, types.ErrUnhandledError))
	}
}

func TestClient_AbortSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	client := Client{
		bucket:   mockBucket,
		segments: make(map[string]*segment.Segment),
	}

	// Test valid segment.
	path := uuid.New().String()
	id := uuid.New().String()
	client.segments[id] = &segment.Segment{
		Path: path,
		ID:   id,
	}
	mockBucket.EXPECT().AbortMultipartUpload(gomock.Any(), gomock.Any()).Do(func(inputPath string, input *service.AbortMultipartUploadInput) {
		assert.Equal(t, path, inputPath)
		assert.Equal(t, id, *input.UploadID)
	})
	err := client.AbortSegment(id)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(client.segments))

	// Test not exist segment.
	id = uuid.New().String()
	err = client.AbortSegment(id)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, segment.ErrSegmentNotInitiated))
}

func TestClient_CompleteSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		id       string
		segments map[string]*segment.Segment
		hasCall  bool
		mockFn   func(string, *service.CompleteMultipartUploadInput)
		hasError bool
		wantErr  error
	}{
		{
			"not initiated segment",
			"", map[string]*segment.Segment{},
			false, nil, true,
			segment.ErrSegmentNotInitiated,
		},
		{
			"segment part empty",
			"test_id",
			map[string]*segment.Segment{
				"test_id": {
					ID:    "test_id",
					Path:  "test_path",
					Parts: nil,
				},
			},
			false, nil,
			true, segment.ErrSegmentPartsEmpty,
		},
		{
			"valid segment",
			"test_id",
			map[string]*segment.Segment{
				"test_id": {
					ID:   "test_id",
					Path: "test_path",
					Parts: map[int64]*segment.Part{
						0: {Offset: 0, Size: 1},
					},
				},
			},
			true,
			func(inputPath string, input *service.CompleteMultipartUploadInput) {
				assert.Equal(t, "test_path", inputPath)
				assert.Equal(t, "test_id", *input.UploadID)
			},
			false, nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if v.hasCall {
				mockBucket.EXPECT().CompleteMultipartUpload(gomock.Any(), gomock.Any()).Do(v.mockFn)
			}

			client := Client{
				bucket:   mockBucket,
				segments: v.segments,
			}

			err := client.CompleteSegment(v.id)
			if v.hasError {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, v.wantErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 0, len(client.segments))
			}
		})
	}
}

func TestClient_Copy(t *testing.T) {
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

		client := Client{
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

func TestClient_Delete(t *testing.T) {
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

		client := Client{
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

func TestClient_InitSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		path     string
		segments map[string]*segment.Segment
		hasCall  bool
		mockFn   func(string, *service.InitiateMultipartUploadInput) (*service.InitiateMultipartUploadOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"valid init segment",
			"test", map[string]*segment.Segment{},
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
			map[string]*segment.Segment{
				"test": {
					ID: "test_segment_id",
				},
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

		client := Client{
			bucket:   mockBucket,
			segments: v.segments,
		}

		_, err := client.InitSegment(v.path, types.WithPartSize(10))
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestClient_ListDir(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	keys := make([]string, 7)
	for idx := range keys {
		keys[idx] = uuid.New().String()
	}

	tests := []struct {
		name   string
		pairs  []*types.Pair
		output *service.ListObjectsOutput
		items  []*types.Object
		err    error
	}{
		{
			"list without delimiter",
			nil,
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
					Name:     keys[0],
					Type:     types.ObjectTypeFile,
					Metadata: make(types.Metadata),
				},
			},
			nil,
		},
		{
			"list with delimiter",
			[]*types.Pair{
				types.WithRecursive(true),
			},
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				CommonPrefixes: []*string{
					service.String(keys[1]),
				},
			},
			[]*types.Object{
				{
					Name:     keys[1],
					Type:     types.ObjectTypeDir,
					Metadata: make(types.Metadata),
				},
			},
			nil,
		},
		{
			"list with return next marker",
			[]*types.Pair{
				types.WithRecursive(true),
			},
			&service.ListObjectsOutput{
				NextMarker: service.String("test_marker"),
				HasMore:    service.Bool(false),
				CommonPrefixes: []*string{
					service.String(keys[2]),
				},
			},
			[]*types.Object{
				{
					Name:     keys[2],
					Type:     types.ObjectTypeDir,
					Metadata: make(types.Metadata),
				},
			},
			nil,
		},
		{
			"list with return empty keys",
			[]*types.Pair{
				types.WithRecursive(true),
			},
			&service.ListObjectsOutput{
				NextMarker: service.String("test_marker"),
				HasMore:    service.Bool(true),
			},
			nil,
			nil,
		},
		{
			"list with error return",
			nil,
			nil,
			nil,
			&qerror.QingStorError{
				StatusCode: 401,
			},
		},
		{
			"list with all data returned",
			nil,
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key:          service.String(keys[5]),
						MimeType:     service.String("application/json"),
						StorageClass: service.String("cool"),
						Etag:         service.String("xxxxx"),
						Size:         service.Int64(1233),
						Modified:     service.Int(1233),
					},
				},
			},
			[]*types.Object{
				{
					Name: keys[5],
					Type: types.ObjectTypeFile,
					Metadata: types.Metadata{
						types.Type:         "application/json",
						types.StorageClass: "cool",
						types.Checksum:     "xxxxx",
						types.Size:         int64(1233),
						types.UpdatedAt:    time.Unix(1233, 0),
					},
				},
			},
			nil,
		},
		{
			"list with return a dir MIME type",
			nil,
			&service.ListObjectsOutput{
				HasMore: service.Bool(false),
				Keys: []*service.KeyType{
					{
						Key:      service.String(keys[6]),
						MimeType: convert.String(DirectoryMIMEType),
					},
				},
			},
			[]*types.Object{
				{
					Name: keys[6],
					Type: types.ObjectTypeDir,
					Metadata: types.Metadata{
						types.Type: DirectoryMIMEType,
					},
				},
			},
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := uuid.New().String()

			mockBucket.EXPECT().ListObjects(gomock.Any()).DoAndReturn(func(input *service.ListObjectsInput) (*service.ListObjectsOutput, error) {
				assert.Equal(t, path, *input.Prefix)
				assert.Equal(t, 200, *input.Limit)
				return v.output, v.err
			})

			client := Client{
				bucket: mockBucket,
			}

			x := client.ListDir(path, v.pairs...)
			for _, expectItem := range v.items {
				item, err := x.Next()
				if v.err != nil {
					assert.Error(t, err)
					assert.True(t, errors.Is(err, v.err))
				}
				assert.NotNil(t, item)
				assert.EqualValues(t, expectItem, item)
			}
			if len(v.items) == 0 {
				item, err := x.Next()
				if v.err != nil {
					assert.True(t, errors.Is(err, types.ErrUnhandledError))
				} else {
					assert.True(t, errors.Is(err, iterator.ErrDone))
				}
				assert.Nil(t, item)
			}
		})
	}
}

func TestClient_Move(t *testing.T) {
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

		client := Client{
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

func TestClient_Read(t *testing.T) {
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

		client := Client{
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

func TestClient_Stat(t *testing.T) {
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
			"valid delete",
			"test_src",
			func(objectKey string, input *service.HeadObjectInput) (*service.HeadObjectOutput, error) {
				assert.Equal(t, "test_src", objectKey)
				length := int64(100)
				return &service.HeadObjectOutput{
					ContentLength:   &length,
					ContentType:     convert.String("test_content_type"),
					ETag:            convert.String("test_etag"),
					XQSStorageClass: convert.String("test_storage_class"),
				}, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		mockBucket.EXPECT().HeadObject(gomock.Any(), gomock.Any()).DoAndReturn(v.mockFn)

		client := Client{
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
			size, ok := o.GetSize()
			assert.True(t, ok)
			assert.Equal(t, int64(100), size)
			contentType, ok := o.GetType()
			assert.True(t, ok)
			assert.Equal(t, "test_content_type", contentType)
			checkSum, ok := o.GetChecksum()
			assert.True(t, ok)
			assert.Equal(t, "test_etag", checkSum)
			storageClass, ok := o.GetStorageClass()
			assert.True(t, ok)
			assert.Equal(t, "test_storage_class", storageClass)
		}
	}
}

func TestClient_Write(t *testing.T) {
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

		client := Client{
			bucket: mockBucket,
		}

		err := client.Write(v.path, nil, types.WithSize(v.size))
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestClient_WriteSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	tests := []struct {
		name     string
		id       string
		segments map[string]*segment.Segment
		offset   int64
		size     int64
		hasCall  bool
		mockFn   func(string, *service.UploadMultipartInput) (*service.UploadMultipartOutput, error)
		hasError bool
		wantErr  error
	}{
		{
			"not initiated segment",
			"", map[string]*segment.Segment{},
			0, 1, false, nil, true,
			segment.ErrSegmentNotInitiated,
		},
		{
			"valid segment",
			"test_id",
			map[string]*segment.Segment{
				"test_id": segment.NewSegment("test_path", "test_id", 1),
			}, 0, 1,
			true,
			func(objectKey string, input *service.UploadMultipartInput) (*service.UploadMultipartOutput, error) {
				assert.Equal(t, "test_path", objectKey)
				assert.Equal(t, "test_id", *input.UploadID)

				return nil, nil
			},
			false, nil,
		},
	}

	for _, v := range tests {
		if v.hasCall {
			mockBucket.EXPECT().UploadMultipart(gomock.Any(), gomock.Any()).Do(v.mockFn)
		}

		client := Client{
			bucket:   mockBucket,
			segments: v.segments,
		}

		err := client.WriteSegment(v.id, v.offset, v.size, nil)
		if v.hasError {
			assert.Error(t, err)
			assert.True(t, errors.Is(err, v.wantErr))
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestClient_ListSegments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBucket := NewMockBucket(ctrl)

	keys := make([]string, 100)
	for idx := range keys {
		keys[idx] = uuid.New().String()
	}

	tests := []struct {
		name   string
		pairs  []*types.Pair
		output *service.ListMultipartUploadsOutput
		items  []*segment.Segment
		err    error
	}{
		{
			"list without delimiter",
			nil,
			&service.ListMultipartUploadsOutput{
				HasMore: service.Bool(false),
				Uploads: []*service.UploadsType{
					{
						Key:      service.String(keys[0]),
						UploadID: service.String(keys[1]),
					},
				},
			},
			[]*segment.Segment{
				segment.NewSegment(keys[0], keys[1], 0),
			},
			nil,
		},
		{
			"list with return next marker",
			[]*types.Pair{
				types.WithRecursive(true),
			},
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
			[]*segment.Segment{
				segment.NewSegment(keys[1], keys[2], 0),
			},
			nil,
		},
		{
			"list with error return",
			nil,
			nil,
			nil,
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

			client := Client{
				bucket:   mockBucket,
				segments: make(map[string]*segment.Segment),
			}

			x := client.ListSegments(path, v.pairs...)
			for _, expectItem := range v.items {
				item, err := x.Next()
				if v.err != nil {
					assert.Error(t, err)
					assert.True(t, errors.Is(err, v.err))
				}
				assert.NotNil(t, item)
				assert.EqualValues(t, expectItem, item)
			}
			if len(v.items) == 0 {
				item, err := x.Next()
				if v.err != nil {
					assert.True(t, errors.Is(err, types.ErrUnhandledError))
				} else {
					assert.True(t, errors.Is(err, iterator.ErrDone))
				}
				assert.Nil(t, item)
			}
		})
	}
}
