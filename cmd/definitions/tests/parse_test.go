package tests

import (
	"errors"
	"sort"
	"testing"

	"github.com/beyondstorage/go-storage/v4/pairs"
	. "github.com/beyondstorage/go-storage/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestParseServicePair(t *testing.T) {
	cases := []struct {
		name string
		k    string
		v    string
		pair Pair
		err  error
	}{
		// {"not parseable type", "interceptor", "", Pair{}, pairs.ErrPairTypeNotParsable},
		{"unknown pair", "not_a_pair", "", Pair{}, pairs.ErrPairNotRegistered},
		{"global pair", "name", "abc", pairs.WithName("abc"), nil},
		{"service pair", "storage_class", "sc", WithStorageClass("sc"), nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := pairs.Parse(Type, tt.k, tt.v)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}
			assert.Equal(t, tt.pair, p)
		})
	}
}

type ByPairName []Pair

func (a ByPairName) Len() int           { return len(a) }
func (a ByPairName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPairName) Less(i, j int) bool { return a[i].Key < a[j].Key }

func TestParseMap(t *testing.T) {
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
			},
			[]Pair{
				pairs.WithName("abc"),
				WithStorageClass("sc"),
			},
			nil,
		},
		{
			"include invalid",
			map[string]string{
				"name":          "abc",
				"not_a_pair":    "",
				"storage_class": "sc",
			},
			nil,
			pairs.ErrPairNotRegistered,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ps, err := pairs.ParseMap(Type, tt.m)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.err))
			}

			sort.Sort(ByPairName(ps))
			sort.Sort(ByPairName(tt.pairs))
			assert.Equal(t, tt.pairs, ps)
		})
	}
}
