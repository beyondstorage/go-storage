package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	cases := []struct {
		name    string
		connStr string
		pairs   []Pair
		err     error
	}{
		{
			"empty",
			"",
			nil,
			services.ErrConnectionStringInvalid,
		},
		{
			"simplest",
			"tests://",
			nil,
			nil,
		},
		{
			"only options",
			"tests://?size=200",
			[]Pair{
				pairs.WithSize(200),
			},
			nil,
		},
		{
			"only root dir",
			"tests:///",
			[]Pair{
				pairs.WithWorkDir("/"),
			},
			nil,
		},
		{
			"end with ?",
			"tests:///?",
			[]Pair{
				pairs.WithWorkDir("/"),
			},
			nil,
		},
		{
			"stupid, but valid (ignored)",
			"tests:///?&??&&&",
			[]Pair{
				pairs.WithWorkDir("/"),
			},
			nil,
		},
		{
			"value can contain all characters except &",
			"tests:///?string_pair=a=b:/c?d&size=200",
			[]Pair{
				pairs.WithWorkDir("/"),
				WithStringPair("a=b:/c?d"),
				pairs.WithSize(200),
			},
			nil,
		},
		{
			"full format",
			"tests://abc/tmp/tmp1?size=200&expire=100&storage_class=sc",
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithWorkDir("/tmp/tmp1"),
				pairs.WithSize(200),
				pairs.WithExpire(time.Duration(100)),
				WithStorageClass("sc"),
			},
			nil,
		},
		{
			"with default pair",
			"tests://abc/tmp/tmp1?default_content_type=application/text&storage_class=sc",
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithWorkDir("/tmp/tmp1"),
				WithDefaultContentType("application/text"),
				WithStorageClass("sc"),
			},
			nil,
		},
		//{
		//	"with feature",
		//	"tests://abc/tmp/tmp1?enable_loose_pair=true&storage_class=sc",
		//	[]Pair{
		//		pairs.WithName("abc"),
		//		pairs.WithWorkDir("/tmp/tmp1"),
		//		WithServiceFeatures(ServiceFeatures{
		//			LoosePair: true,
		//		}),
		//		WithStorageClass("sc"),
		//	},
		//	nil,
		//},
		{
			"duplicate key, appear in order (finally, first will be picked)",
			"tests://abc/tmp/tmp1?size=200&name=def&size=300",
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithWorkDir("/tmp/tmp1"),
				pairs.WithSize(200),
				pairs.WithName("def"),
				pairs.WithSize(300),
			},
			nil,
		},
		{
			"not registered pair",
			"tests://abc/tmp?not_a_pair=a",
			nil,
			services.ErrConnectionStringInvalid,
		},
		{
			"key without value is ignored (even not registered pair)",
			"tests://abc/tmp?not_a_pair&&",
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithWorkDir("/tmp"),
			},
			nil,
		},
		{
			"not parsable pair",
			"tests://abc/tmp?io_call_back=a",
			nil,
			services.ErrConnectionStringInvalid,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			servicer, err := services.NewServicerFromString(tt.connStr)
			service, ok := servicer.(*Service)

			if tt.err == nil {
				assert.Nil(t, err)
				assert.True(t, ok)
			} else {
				assert.True(t, errors.Is(err, tt.err))
				return
			}

			assert.Equal(t, service.Pairs, tt.pairs)
		})
	}
}
