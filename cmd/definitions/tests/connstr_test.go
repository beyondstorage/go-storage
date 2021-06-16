package tests

import (
	"errors"
	"testing"

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
			"both <name> and <work_dir> missing",
			"tests://",
			nil,
			services.ErrConnectionStringInvalid,
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
				pairs.WithSize(200),
				WithStringPair("a=b:/c?d"),
			},
			nil,
		},
		{
			"full format",
			"tests://abc/tmp?size=200&expire=100&storage_class=sc",
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithWorkDir("/tmp"),
				pairs.WithExpire(100),
				pairs.WithSize(200),
				WithStorageClass("sc"),
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
			"not parseable pair",
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

			assert.ElementsMatch(t, service.Pairs, tt.pairs)
		})
	}
}
