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
		name string
		cfg  string
		t    string
		opt  []*types.Pair
		err  error
	}{
		{
			"invalid config",
			"xxx://",
			"",
			nil,
			ErrInvalidConfig,
		},
		{
			"no credential, endpoint and options",
			"posixfs:///",
			"posixfs",
			nil,
			nil,
		},
		{
			"no credential, endpoint, but with options",
			"posixfs:///?work_dir=/path&test_key=test_value",
			"posixfs",
			[]*types.Pair{
				pairs.WithWorkDir("/path"),
				{
					Key:   "test_key",
					Value: "test_value",
				},
			},
			nil,
		},
		{
			"no endpoint, but with credential and options",
			"qingstor://hmac:ak:sk/?name=path",
			"qingstor",
			[]*types.Pair{
				pairs.WithCredential(credential.MustNewHmac("ak", "sk")),
				pairs.WithName("path"),
			},
			nil,
		},
		{
			"all elements available",
			"qingstor://hmac:ak:sk@https:qingstor.com:443/?name=path",
			"qingstor",
			[]*types.Pair{
				pairs.WithEndpoint(endpoint.NewHTTPS("qingstor.com", 443)),
				pairs.WithCredential(credential.MustNewHmac("ak", "sk")),
				pairs.WithName("path"),
			},
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			typ, opt, err := Parse(tt.cfg)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.Equal(t, tt.t, typ)
			assert.ElementsMatch(t, tt.opt, opt)
		})
	}
}
