package iowrap

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/beyondstorage/go-storage/v4/pkg/randbytes"
)

func BenchmarkStdPipe(b *testing.B) {
	cases := []struct {
		name string
		size int
	}{
		{"4k", 4 * 1024},
		{"64k", 64 * 1024},
		{"4m", 4 * 1024 * 1024},
	}

	for _, v := range cases {
		b.Run(v.name, func(b *testing.B) {
			content := make([]byte, v.size)
			_, err := randbytes.NewRand().Read(content)
			if err != nil {
				b.Error(err)
			}

			r, w := io.Pipe()

			go func() {
				_, _ = io.Copy(ioutil.Discard, r)
			}()

			b.SetBytes(int64(v.size))
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				_, _ = w.Write(content)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkIowrapPipe(b *testing.B) {
	cases := []struct {
		name string
		size int
	}{
		{"4k", 4 * 1024},
		{"64k", 64 * 1024},
		{"4m", 4 * 1024 * 1024},
	}

	for _, v := range cases {
		b.Run(v.name, func(b *testing.B) {
			content := make([]byte, v.size)
			_, err := randbytes.NewRand().Read(content)
			if err != nil {
				b.Error(err)
			}

			r, w := Pipe()

			go func() {
				_, _ = io.Copy(ioutil.Discard, r)
			}()

			b.SetBytes(int64(v.size))
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				_, _ = w.Write(content)
			}
			b.StopTimer()
		})
	}
}
