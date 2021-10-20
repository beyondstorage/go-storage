package memory

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/uuid"

	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func setup(b *testing.B, size int64) (store *Storage, path string) {
	root := newObject("", nil, types.ModeDir)
	root.parent = root

	store = &Storage{
		root:    root,
		workDir: "/",
	}

	path = uuid.NewString()
	content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
	if err != nil {
		b.Fatal(err)
	}

	_, err = store.Write(path, bytes.NewReader(content), size)
	if err != nil {
		b.Fatal(err)
	}
	return
}

func BenchmarkStorage_Read(b *testing.B) {
	cases := []struct {
		name string
		size int64
	}{
		{"64B", 64},
		{"4k", 4 * 1024},
		{"64M", 64 * 1024 * 1024},
	}
	for _, v := range cases {
		b.Run(v.name, func(b *testing.B) {
			store, path := setup(b, v.size)

			b.SetBytes(v.size)
			for i := 0; i < b.N; i++ {
				_, _ = store.Read(path, io.Discard)
			}
		})
	}
}

func BenchmarkStorage_Write(b *testing.B) {
	cases := []struct {
		name string
		size int64
	}{
		{"64B", 64},
		{"4k", 4 * 1024},
		{"64M", 64 * 1024 * 1024},
	}
	for _, v := range cases {
		b.Run(v.name, func(b *testing.B) {
			store, _ := setup(b, v.size)

			path := uuid.NewString()
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), v.size))
			if err != nil {
				b.Fatal(err)
			}

			b.SetBytes(v.size)
			for i := 0; i < b.N; i++ {
				_, _ = store.Write(path, bytes.NewReader(content), v.size)
			}
		})
	}
}
