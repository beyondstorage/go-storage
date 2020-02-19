package qingstor

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/services"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingstor-sdk-go/v3/config"
	qerror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage/types/pairs"
)

func TestService_String(t *testing.T) {
	accessKey := uuid.New().String()
	secretKey := uuid.New().String()

	srv, _, err := New(pairs.WithCredential(credential.MustNewHmac(accessKey, secretKey)))
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
		monkey.PatchInstanceMethod(reflect.TypeOf(srv.noRedirectClient), "Head", fn)

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
	assert.True(t, errors.Is(err, services.ErrPairRequired))

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
