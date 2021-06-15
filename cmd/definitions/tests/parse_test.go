package tests

import (
	"errors"
	"sort"
	"testing"

	"github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/services"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/stretchr/testify/assert"
)

type byPairName []Pair

func (a byPairName) Len() int           { return len(a) }
func (a byPairName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPairName) Less(i, j int) bool { return a[i].Key < a[j].Key }

func TestFromMap(t *testing.T) {
	cases := []struct {
		name  string
		m     map[string]string
		pairs []Pair
		err   error
	}{
		{
			"empty",
			map[string]string{},
			nil,
			nil,
		},
		{
			"global & service",
			map[string]string{
				"name":          "abc",
				"storage_class": "sc",
				"size":          "200",
				"expire":        "100",
			},
			[]Pair{
				pairs.WithName("abc"),
				pairs.WithExpire(100),
				pairs.WithSize(200),
				WithStorageClass("sc"),
			},
			nil,
		},
		{
			"not registered pair",
			map[string]string{
				"name":          "abc",
				"not_a_pair":    "",
				"storage_class": "sc",
			},
			nil,
			services.ErrPairNotRegistered,
		},
		{
			"not parseable pair",
			map[string]string{
				"name":          "abc",
				"io_callback":   "",
				"storage_class": "sc",
			},
			nil,
			services.ErrPairTypeNotParsable,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			servicer, err := services.NewServicerFromString("")
			service, ok := servicer.(*Service)

			if tt.err == nil {
				assert.Nil(t, err)
				assert.True(t, ok)
			} else {
				assert.True(t, errors.Is(err, tt.err))
				return
			}

			sort.Sort(byPairName(service.Pairs))
			sort.Sort(byPairName(tt.pairs))
			assert.Equal(t, tt.pairs, service.Pairs)
		})
	}
}
