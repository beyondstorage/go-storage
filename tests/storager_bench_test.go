package tests

import "testing"

func BenchmarkStorage_Stat(b *testing.B) {
	s, err := NewStorager()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = s.Stat("abc")
	}
}
