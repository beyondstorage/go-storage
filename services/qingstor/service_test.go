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
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	qerror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/types"
)

func TestService_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Missing required pair
	srv := Service{}
	err := srv.Init()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, types.ErrPairRequired))

	// Valid case
	srv = Service{}
	accessKey := uuid.New().String()
	secretKey := uuid.New().String()
	host := uuid.New().String()
	port := 1234
	protocol := uuid.New().String()
	err = srv.Init(
		types.WithAccessKey(accessKey),
		types.WithSecretKey(secretKey),
		types.WithHost(host),
		types.WithPort(port),
		types.WithProtocol(protocol),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv.service)
	assert.NotNil(t, srv.config)
	assert.Equal(t, srv.config.AccessKeyID, accessKey)
	assert.Equal(t, srv.config.SecretAccessKey, secretKey)
	assert.Equal(t, srv.config.Host, host)
	assert.Equal(t, srv.config.Port, port)
	assert.Equal(t, srv.config.Protocol, protocol)
}

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)

	{
		// Test case 1: with location
		srv := Service{
			service: mockService,
		}

		name := uuid.New().String()
		location := uuid.New().String()

		mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).Do(func(bucketName, inputLocation string) {
			assert.Equal(t, name, bucketName)
			assert.Equal(t, location, inputLocation)
		})

		s, err := srv.Get(name, types.WithLocation(location))
		assert.NoError(t, err)
		assert.NotNil(t, s)
	}

	{
		// Test case 2: without location
		srv := New()
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

		expectURL := fmt.Sprintf("%s://%s.%s:%d", srv.config.Protocol, name, srv.config.Host, srv.config.Port)

		// Patch http Head.
		fn := func(client *http.Client, url string) (*http.Response, error) {
			assert.Equal(t, expectURL, url)

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
		monkey.PatchInstanceMethod(reflect.TypeOf(srv.noRedirectClient), "Head", fn)

		// Mock Bucket.
		mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).Do(func(bucketName, inputLocation string) {
			assert.Equal(t, name, bucketName)
			assert.Equal(t, location, inputLocation)
		})

		s, err := srv.Get(name)
		assert.NoError(t, err)
		assert.NotNil(t, s)
	}
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
	assert.True(t, errors.Is(err, types.ErrPairRequired))

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

	_, err = srv.Create(path, types.WithLocation(location))
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

		err := srv.Delete(name, types.WithLocation(location))
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

	mockService.EXPECT().Bucket(gomock.Any(), gomock.Any()).AnyTimes()

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

		s, err := srv.List(types.WithLocation(location))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(s))
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

		s, err := srv.List()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(s))
	}

	{
		// Test request facing error.
		mockService.EXPECT().ListBuckets(gomock.Any()).DoAndReturn(func(input *service.ListBucketsInput) (*service.ListBucketsOutput, error) {
			return nil, &qerror.QingStorError{}
		})

		_, err := srv.List()
		assert.Error(t, err)
		assert.True(t, errors.Is(err, types.ErrUnhandledError))
	}
}

func ExampleService_Init() {
	srv := New()
	err := srv.Init(
		types.WithAccessKey("test_access_key"),
		types.WithSecretKey("test_secret_key"),
	)
	if err != nil {
		log.Printf("service init failed: %v", err)
	}
}

func ExampleService_Get() {
	srv := New()
	err := srv.Init(
		types.WithAccessKey("test_access_key"),
		types.WithSecretKey("test_secret_key"),
	)
	if err != nil {
		log.Printf("service init failed: %v", err)
	}

	store, err := srv.Get("bucket_name", types.WithLocation("location"))
	if err != nil {
		log.Printf("service get bucket failed: %v", err)
	}
	log.Printf("%v", store)
}
