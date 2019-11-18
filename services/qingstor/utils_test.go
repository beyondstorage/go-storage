package qingstor

import (
	"errors"
	"testing"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/stretchr/testify/assert"
	qserror "github.com/yunify/qingstor-sdk-go/v3/request/errors"
)

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
		{"under sub dir", "/root", "abc", "root/abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{}
			err := client.Init(pairs.WithWorkDir(tt.base))
			if err != nil {
				t.Error(err)
			}

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
		{"under sub dir", "/root", "root/abc", "abc"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			err := client.Init(pairs.WithWorkDir(tt.base))
			if err != nil {
				t.Error(err)
			}

			gotPath := client.getRelPath(tt.path)
			assert.Equal(t, tt.expectedPath, gotPath)
		})
	}
}

func TestHandleQingStorError(t *testing.T) {
	t.Run("nil error will panic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = handleQingStorError(nil)
		})
	})

	t.Run("non-qingstor error will return a ErrUnhandledError", func(t *testing.T) {
		err := handleQingStorError(errors.New("test"))
		assert.True(t, errors.Is(err, types.ErrUnhandledError))
	})

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
				types.ErrObjectNotExist,
			},
			{
				"invalid status code",
				&qserror.QingStorError{
					StatusCode:   444,
					Code:         "",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				types.ErrUnhandledError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.True(t, errors.Is(handleQingStorError(tt.input), tt.expected))
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
				types.ErrPermissionDenied,
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
				types.ErrObjectNotExist,
			},
			{
				"invalid_access_key_id",
				&qserror.QingStorError{
					StatusCode:   400,
					Code:         "invalid_access_key_id",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				types.ErrConfigIncorrect,
			},
			{
				"not handled",
				&qserror.QingStorError{
					StatusCode:   400,
					Code:         "xxxxxx",
					Message:      "",
					RequestID:    "",
					ReferenceURL: "",
				},
				types.ErrUnhandledError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.True(t, errors.Is(handleQingStorError(tt.input), tt.expected))
			})
		}
	}
}
