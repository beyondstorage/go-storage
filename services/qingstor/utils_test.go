package qingstor

import (
	"errors"
	"testing"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/pkg/storageclass"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Missing required pair
	_, _, err := New()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, services.ErrRestrictionDissatisfied))

	// Valid case
	accessKey := uuid.New().String()
	secretKey := uuid.New().String()
	host := uuid.New().String()
	port := 1234
	srv, _, err := New(
		pairs.WithCredential(credential.MustNewHmac(accessKey, secretKey)),
		pairs.WithEndpoint(endpoint.NewHTTP(host, port)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv)
}

func TestIsBucketNameValid(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"start with letter", "a-bucket-test", true},
		{"start with digit", "0-bucket-test", true},
		{"start with strike", "-bucket-test", false},
		{"end with strike", "bucket-test-", false},
		{"too short", "abcd", false},
		{"too long (64)", "abcdefghijklmnopqrstuvwxyz123456abcdefghijklmnopqrstuvwxyz123456", false},
		{"contains illegal char", "abcdefg_1234", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBucketNameValid(tt.args); got != tt.want {
				t.Errorf("IsBucketNameValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAbsPath(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		expectedPath string
	}{
		{"under root", "/", "abc", "abc"},
		{"under sub dir", "/root/", "abc", "root/abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := Storage{workDir: tt.base[1:]}

			gotPath := client.getAbsPath(tt.path)
			assert.Equal(t, tt.expectedPath, gotPath)
		})
	}
}

func TestGetRelPath(t *testing.T) {
	cases := []struct {
		name         string
		base         string
		path         string
		expectedPath string
	}{
		{"under root", "/", "abc", "abc"},
		{"under sub dir", "/root/", "root/abc", "abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := &Storage{workDir: tt.base[1:]}

			gotPath := client.getRelPath(tt.path)
			assert.Equal(t, tt.expectedPath, gotPath)
		})
	}
}

func TestHandleQingStorError(t *testing.T) {
	{
		tests := []struct {
			name     string
			input    *qserror.QingStorError
			expected error
		}{
			{
				"not found",
				&qserror.QingStorError{
					StatusCode:   404,
					Code:         "",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				services.ErrObjectNotExist,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.True(t, errors.Is(formatError(tt.input), tt.expected))
			})
		}
	}

	{
		tests := []struct {
			name     string
			input    *qserror.QingStorError
			expected error
		}{
			{
				"permission_denied",
				&qserror.QingStorError{
					StatusCode:   403,
					Code:         "permission_denied",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				services.ErrPermissionDenied,
			},
			{
				"object_not_exists",
				&qserror.QingStorError{
					StatusCode:   404,
					Code:         "object_not_exists",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				services.ErrObjectNotExist,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.True(t, errors.Is(formatError(tt.input), tt.expected))
			})
		}
	}
}

func TestParseStorageClass(t *testing.T) {
	tests := []struct {
		name        string
		input       storageclass.Type
		expected    string
		expectedErr error
	}{
		{"hot", storageclass.Hot, storageClassStandard, nil},
		{"warm", storageclass.Warm, storageClassStandardIA, nil},
		{"xxxxx", "xxxx", "", services.ErrCapabilityInsufficient},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStorageClass(tt.input)
			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestFormatStorageClass(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected storageclass.Type
	}{
		{"hot", storageClassStandard, storageclass.Hot},
		{"warm", storageClassStandardIA, storageclass.Warm},
		{"xxxxx", "xxxxx", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStorageClass(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
