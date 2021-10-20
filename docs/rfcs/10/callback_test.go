package main

import (
	"bytes"
	"io"
	"testing"
)

// Use 4K read for benchmark
var n int64 = 4 * 1024
var content = bytes.Repeat([]byte{'x'}, int(n))

func BenchmarkPlainReader(b *testing.B) {
	b.SetBytes(n)
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(content)
		_, _ = io.ReadAll(r)
	}
}

type IntCallbackReader struct {
	r  io.Reader
	fn func(int)
}

func (r IntCallbackReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.fn(n)
	return n, err
}

func BenchmarkIntCallbackReader(b *testing.B) {
	b.SetBytes(n)
	for i := 0; i < b.N; i++ {
		x := 0
		r := IntCallbackReader{
			r: bytes.NewReader(content),
			fn: func(i int) {
				x += i
			},
		}
		_, _ = io.ReadAll(r)
	}
}

type BytesCallbackReader struct {
	r  io.Reader
	fn func([]byte)
}

func (r BytesCallbackReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.fn(p[:n])
	return n, err
}

func BenchmarkBytesCallbackReader(b *testing.B) {
	b.SetBytes(n)
	for i := 0; i < b.N; i++ {
		x := 0
		r := BytesCallbackReader{
			r: bytes.NewReader(content),
			fn: func(i []byte) {
				x += len(i)
			},
		}
		_, _ = io.ReadAll(r)
	}
}
