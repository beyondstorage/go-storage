package tests

import (
	"context"
	"testing"

	"github.com/beyondstorage/go-storage/v4/types"
)

func BenchmarkStorage_Stat(b *testing.B) {
	s, err := NewStorager()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = s.Stat("abc")
	}
}

func BenchmarkStorage_List(b *testing.B) {
	ctx := context.TODO()
	
	var ob []*types.Object
	for i := 0; i < 1024; i++ {
		ob = append(ob, &types.Object{})
	}

	s := &Storage{
		objects: ob,
	}

	it, err := s.list(ctx, "", pairStorageList{})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = it.Next()
	}
}
