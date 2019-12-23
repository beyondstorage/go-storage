package config

import (
	"errors"
	"testing"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name      string
		cfg       string
		t         string
		namespace string
		opt       []*types.Pair
		err       error
	}{
		{
			"invalid config",
			"xxx://",
			"",
			"",
			nil,
			ErrInvalidConfig,
		},
		{
			"no credential, endpoint and options",
			"posixfs:///path",
			"posixfs",
			"path",
			nil,
			nil,
		},
		{
			"no credential, endpoint, but with options",
			"posixfs:///path?test_key=test_value",
			"posixfs",
			"path",
			nil,
			nil,
		},
		{
			"no credential, but with endpoint and options",
			"qingstor://https:qingstor.com:443/path",
			"qingstor",
			"path",
			[]*types.Pair{
				pairs.WithEndpoint(endpoint.NewHTTPS("qingstor.com", 443)),
			},
			nil,
		},
		{
			"all elements available",
			"qingstor://static:ak:sk@https:qingstor.com:443/path",
			"qingstor",
			"path",
			[]*types.Pair{
				pairs.WithEndpoint(endpoint.NewHTTPS("qingstor.com", 443)),
				pairs.WithCredential(credential.NewStatic("ak", "sk")),
			},
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			typ, id, opt, err := Parse(tt.cfg)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.Equal(t, tt.t, typ)
			assert.Equal(t, tt.namespace, id)
			assert.ElementsMatch(t, tt.opt, opt)
		})
	}
}
