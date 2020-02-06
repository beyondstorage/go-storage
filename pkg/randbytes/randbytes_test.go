package randbytes

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"
)

func TestRand(t *testing.T) {
	buf := make([]byte, 16)
	n, err := NewRand().Read(buf)
	if err != nil {
		t.Fatalf("Error reading: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("Short read: %v", n)
	}
	t.Logf("Read %x", buf)
}

const toCopy = 1024 * 1024

func BenchmarkRand(b *testing.B) {
	b.SetBytes(toCopy)
	r := NewRand()
	for i := 0; i < b.N; i++ {
		_, _ = io.CopyN(ioutil.Discard, r, toCopy)
	}
}

func BenchmarkCrypto(b *testing.B) {
	b.SetBytes(toCopy)
	for i := 0; i < b.N; i++ {
		_, _ = io.CopyN(ioutil.Discard, rand.Reader, toCopy)
	}
}
