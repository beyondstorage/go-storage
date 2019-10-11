package qingstor

import (
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
