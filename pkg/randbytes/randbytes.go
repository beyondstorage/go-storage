package randbytes

import (
	"io"
	"math/rand"
	"time"
)

// Rand creates a stream of non-crypto quality random bytes
type Rand struct {
	rand.Source
}

// NewRand creates a new random reader with a time source.
func NewRand() io.Reader {
	return &Rand{rand.NewSource(time.Now().UnixNano())}
}

// Read satisfies io.Reader
func (r *Rand) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := r.Int63()
		for i := 0; i < 7; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}
}
