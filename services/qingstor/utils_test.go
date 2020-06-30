package qingstor

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pengsrc/go-shared/convert"
	qserror "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/stretchr/testify/assert"

	"github.com/aos-dev/go-storage/v2/pkg/credential"
	"github.com/aos-dev/go-storage/v2/pkg/endpoint"
	"github.com/aos-dev/go-storage/v2/services"
	"github.com/aos-dev/go-storage/v2/types/pairs"
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
	name := uuid.New().String()
	port := 1234
	srv, store, err := New(
		pairs.WithCredential(credential.MustNewHmac(accessKey, secretKey)),
		pairs.WithEndpoint(endpoint.NewHTTP(host, port)),
		pairs.WithLocation("test"),
		pairs.WithName(name),
	)
	assert.NoError(t, err)
	assert.NotNil(t, srv)
	assert.NotNil(t, store)
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
		{"under prefix", "/root", "/abc", "root/abc"},
		{"under prefix ending with /", "/root/", "abc", "root/abc"},
		{"under unexpected prefix", "//abc", "/def", "/abc/def"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := Storage{workDir: tt.base}

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
		{"under prefix", "/root", "root/abc", "/abc"},
		{"under prefix ending with /", "/root/", "root/abc", "abc"},
		{"under unexpected prefix", "//abc", "/abc/def", "/def"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := &Storage{workDir: tt.base}

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

func Test_isObjectDirectory(t *testing.T) {
	tests := []struct {
		name    string
		keyType string
		want    bool
	}{
		{
			name:    "false",
			keyType: "application/octet-stream",
			want:    false,
		},
		{
			name:    "true",
			keyType: "application/x-directory",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &service.KeyType{
				MimeType: convert.String(tt.keyType),
			}
			if got := isObjectDirectory(o); got != tt.want {
				t.Errorf("isObjectDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
