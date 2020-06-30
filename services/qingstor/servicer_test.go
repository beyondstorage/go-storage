package qingstor

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/qingstor/qingstor-sdk-go/v4/config"
	qerror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/stretchr/testify/assert"

	"github.com/aos-dev/go-storage/v2"
	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types"
	"github.com/aos-dev/go-storage/v2/types/pairs"
)

func TestService_String(t *testing.T) {
	accessKey := uuid.New().String()
	secretKey := uuid.New().String()

	srv, err := NewServicer(pairs.WithCredential(credential.MustNewHmac(accessKey, secretKey)))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, srv.String())
}

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("with location", func(t *testing.T) {
		mockService := NewMockService(ctrl)

		srv := Service{
			service: mockService,
		}

		name := uuid.New().String()
		location := uuid.New().String()

		mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).DoAndReturn(func(bucketName, inputLocation string) (*service.Bucket, error) {
			assert.Equal(t, name, bucketName)
			assert.Equal(t, location, inputLocation)
			return &service.Bucket{
				Properties: &service.Properties{
					BucketName: &name,
					Zone:       &location,
				},
			}, nil
		})

		s, err := srv.Get(name, pairs.WithLocation(location))
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("without location", func(t *testing.T) {
		mockService := NewMockService(ctrl)

		srv := &Service{}
		srv.service = mockService
		srv.config = &config.Config{
			AccessKeyID:     uuid.New().String(),
			SecretAccessKey: uuid.New().String(),
			Host:            uuid.New().String(),
			Port:            1234,
			Protocol:        uuid.New().String(),
		}

		name := uuid.New().String()
		location := uuid.New().String()

		// Patch http Head.
		fn := func(client *http.Client, url string) (*http.Response, error) {
			header := http.Header{}
			header.Set(
				"location",
				fmt.Sprintf("%s://%s.%s.%s",
					srv.config.Protocol, name, location, srv.config.Host),
			)
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header:     header,
			}, nil
		}
		monkey.PatchInstanceMethod(reflect.TypeOf(srv.client), "Head", fn)

		// Mock Bucket.
		mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).DoAndReturn(func(bucketName, inputLocation string) (*service.Bucket, error) {
			assert.Equal(t, name, bucketName)
			assert.Equal(t, location, inputLocation)
			return &service.Bucket{
				Properties: &service.Properties{
					BucketName: &name,
					Zone:       &location,
				},
			}, nil
		})

		s, err := srv.Get(name)
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("invalid bucket name", func(t *testing.T) {
		mockService := NewMockService(ctrl)

		srv := &Service{}
		srv.service = mockService
		srv.config = &config.Config{
			AccessKeyID:     uuid.New().String(),
			SecretAccessKey: uuid.New().String(),
			Host:            uuid.New().String(),
			Port:            1234,
			Protocol:        uuid.New().String(),
		}

		s, err := srv.Get("1234")
		assert.Error(t, err)
		assert.Nil(t, s)
	})
}

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	srv := Service{
		service: mockService,
	}

	// Test case1: without location
	path := uuid.New().String()
	_, err := srv.Create(path)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, services.ErrRestrictionDissatisfied))

	// Test case2: with valid location.
	path = uuid.New().String()
	location := uuid.New().String()

	// Monkey the bucket's Put method
	bucket := &service.Bucket{}
	fn := func(*service.Bucket) (*service.PutBucketOutput, error) {
		t.Log("Bucket put has been called")
		return &service.PutBucketOutput{}, nil
	}
	monkey.PatchInstanceMethod(reflect.TypeOf(bucket), "Put", fn)

	mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).Do(func(inputPath, inputLocation string) {
		assert.Equal(t, path, inputPath)
		assert.Equal(t, location, inputLocation)
	}).Return(bucket, nil)

	_, err = srv.Create(path, pairs.WithLocation(location))
	assert.NoError(t, err)
}

func TestService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	srv := Service{
		service: mockService,
	}

	{
		name := uuid.New().String()
		location := uuid.New().String()

		// Patch bucket.Delete
		bucket := &service.Bucket{}
		fn := func(*service.Bucket) (*service.DeleteBucketOutput, error) {
			t.Log("Bucket delete has been called")
			return nil, nil
		}
		monkey.PatchInstanceMethod(reflect.TypeOf(bucket), "Delete", fn)

		mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).DoAndReturn(func(bucketName, inputLocation string) (*service.Bucket, error) {
			assert.Equal(t, name, bucketName)
			assert.Equal(t, location, inputLocation)
			return bucket, nil
		})

		err := srv.Delete(name, pairs.WithLocation(location))
		assert.NoError(t, err)
	}
}

func TestService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	srv := &Service{
		service: mockService,
	}
	listFunc := pairs.WithStoragerFunc(func(storager storage.Storager) {})

	mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).DoAndReturn(func(inputName, inputLocation string) (*service.Bucket, error) {
		return &service.Bucket{
			Config: &config.Config{},
		}, nil
	}).AnyTimes()

	{
		// Test request with location.
		name := uuid.New().String()
		location := uuid.New().String()

		mockService.EXPECT().ListBuckets(gomock.Any()).DoAndReturn(func(input *service.ListBucketsInput) (*service.ListBucketsOutput, error) {
			assert.Equal(t, location, *input.Location)
			return &service.ListBucketsOutput{
				Buckets: []*service.BucketType{
					{Name: &name, Location: &location},
				},
			}, nil
		})

		err := srv.List(pairs.WithLocation(location), listFunc)
		assert.NoError(t, err)
		// assert.Equal(t, 1, len(s))
	}

	{
		// Test request without location.
		name := uuid.New().String()
		location := uuid.New().String()

		mockService.EXPECT().ListBuckets(gomock.Any()).DoAndReturn(func(input *service.ListBucketsInput) (*service.ListBucketsOutput, error) {
			assert.Nil(t, input.Location)
			return &service.ListBucketsOutput{
				Buckets: []*service.BucketType{
					{Name: &name, Location: &location},
				},
			}, nil
		})

		err := srv.List(listFunc)
		assert.NoError(t, err)
		// assert.Equal(t, 1, len(s))
	}

	{
		// Test request facing error.
		mockService.EXPECT().ListBuckets(gomock.Any()).DoAndReturn(func(input *service.ListBucketsInput) (*service.ListBucketsOutput, error) {
			return nil, &qerror.QingStorError{}
		})

		err := srv.List(listFunc)
		t.Log(err)
		assert.Error(t, err)
	}
}

func ExampleNew() {
	_, _, err := New(
		pairs.WithCredential(
			credential.MustNewHmac("test_access_key", "test_secret_key"),
		),
	)
	if err != nil {
		log.Printf("service init failed: %v", err)
	}
}

func ExampleService_Get() {
	srv, _, err := New(
		pairs.WithCredential(
			credential.MustNewHmac("test_access_key", "test_secret_key"),
		),
	)
	if err != nil {
		log.Printf("service init failed: %v", err)
	}

	store, err := srv.Get("bucket_name", pairs.WithLocation("location"))
	if err != nil {
		log.Printf("service get bucket failed: %v", err)
	}
	log.Printf("%v", store)
}

func Test_isWorkDirValid(t *testing.T) {
	type args struct {
		wd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not a dir",
			args: args{wd: "/path/to/file"},
			want: false,
		},
		{
			name: "not a abs dir",
			args: args{wd: "path/to/file/"},
			want: false,
		},
		{
			name: "start with multi-slash",
			args: args{wd: "///multi-slash/to/file/"},
			want: false,
		},
		{
			name: "end with multi-slash",
			args: args{wd: "/path/to/multi-slash///"},
			want: false,
		},
		{
			name: "root path",
			args: args{wd: "/"},
			want: true,
		},
		{
			name: "normal abs dir",
			args: args{wd: "/path/to/dir/"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWorkDirValid(tt.args.wd)
			if got != tt.want {
				t.Errorf("isWorkDirValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_newStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blankBucket := service.Bucket{
		Config: &config.Config{},
	}
	validWorkDir, invalidWorkDir := "/valid/dir/", "invalid/dir"
	type args struct {
		pairs []*types.Pair
	}
	tests := []struct {
		name       string
		wd         string
		args       args
		wantBucket *service.Bucket
		targetErr  error
		wantErr    bool
	}{
		{
			name: "normal case",
			wd:   validWorkDir,
			args: args{[]*types.Pair{
				{Key: pairs.Location, Value: uuid.New().String()},
				{Key: pairs.Name, Value: uuid.New().String()},
				{Key: pairs.WorkDir, Value: validWorkDir},
			}},
			wantBucket: &blankBucket,
			targetErr:  nil,
			wantErr:    false,
		},
		{
			name: "invalid work dir",
			wd:   invalidWorkDir,
			args: args{[]*types.Pair{
				{Key: pairs.Location, Value: uuid.New().String()},
				{Key: pairs.Name, Value: uuid.New().String()},
				{Key: pairs.WorkDir, Value: invalidWorkDir},
			}},
			wantBucket: nil,
			targetErr:  ErrInvalidWorkDir,
			wantErr:    true,
		},
		{
			name: "blank work dir",
			wd:   "/",
			args: args{[]*types.Pair{
				{Key: pairs.Location, Value: uuid.New().String()},
				{Key: pairs.Name, Value: uuid.New().String()},
			}},
			wantBucket: &blankBucket,
			targetErr:  nil,
			wantErr:    false,
		},
		{
			name: "invalid bucket name",
			wd:   validWorkDir,
			args: args{[]*types.Pair{
				{Key: pairs.Location, Value: uuid.New().String()},
				{Key: pairs.Name, Value: "invalid bucket name"},
				{Key: pairs.WorkDir, Value: validWorkDir},
			}},
			wantBucket: nil,
			targetErr:  ErrInvalidBucketName,
			wantErr:    true,
		},
		{
			name:       "no pairs, fail when parse",
			args:       args{},
			wantBucket: nil,
			targetErr:  services.ErrRestrictionDissatisfied,
			wantErr:    true,
		},
	}

	mockService := NewMockService(ctrl)
	mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).DoAndReturn(func(inputName, inputLocation string) (*service.Bucket, error) {
		return &blankBucket, nil
	}).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{service: mockService}
			gotStore, err := s.newStorage(tt.args.pairs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr || err != nil {
				assert.Nil(t, gotStore, tt.name)
				assert.True(t, errors.Is(err, tt.targetErr), tt.name)
			} else {
				assert.Equal(t, tt.wantBucket, gotStore.bucket, tt.name)
				assert.Equal(t, tt.wd, gotStore.workDir, tt.name)
			}
		})
	}
}
