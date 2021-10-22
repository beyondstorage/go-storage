//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package fs

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func fsReaddir(b *testing.B) {
	s, _ := newStorager()

	it, err := s.List("/usr/lib")
	if err != nil {
		b.Error(err)
	}

	for {
		_, err := it.Next()
		if err == types.IterateDone {
			break
		}
	}
}

func osReaddir(b *testing.B) {
	_, err := os.ReadDir("/usr/lib")
	if err != nil {
		b.Error(err)
	}
}

func TestGetFilesFs(t *testing.T) {
	s, _ := newStorager()

	it, err := s.List("/usr/lib")
	if err != nil {
		t.Error(err)
	}

	for {
		o, err := it.Next()
		if err == types.IterateDone {
			break
		}
		assert.NotNil(t, o)
	}
}

func TestIssue68(t *testing.T) {
	// There are fuzzy logic in testIssue68.
	// Run it 100 times to make sure everything is ok.
	for i := 0; i < 100; i++ {
		// We will create upto 1000 files, introduce rand for fuzzing.
		numbers := 225 + rand.Intn(800)

		t.Run(fmt.Sprintf("list %d files", numbers), func(t *testing.T) {
			testIssue68(t, numbers)
		})
	}
}

// This test case intends to reproduce issue #68.
//
// ref: https://github.com/beyondstorage/go-service-fs/issues/68
func testIssue68(t *testing.T, numbers int) {
	tmpDir := t.TempDir()

	store, err := newStorager(ps.WithWorkDir(tmpDir))
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	// Create enough files in a dir, the file name must be long enough.
	// So that the total size will bigger than 8196.
	for i := 0; i < numbers; i++ {
		// uuid's max size is 36.
		// We use rand here for fuzzing.
		filename := uuid.NewString()[:1+rand.Intn(35)]

		f, err := os.Create(path.Join(tmpDir, filename))
		if err != nil {
			t.Error(err)
		}
		err = f.Close()
		if err != nil {
			t.Error(err)
		}
	}

	expected := make(map[string]struct{})
	fi, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Error(err)
	}
	for _, v := range fi {
		expected[v.Name()] = struct{}{}
	}

	actual := make(map[string]struct{})
	it, err := store.List("")
	if err != nil {
		t.Error(err)
	}
	for {
		o, err := it.Next()
		if err == types.IterateDone {
			break
		}
		_, exist := actual[o.Path]
		if exist {
			t.Errorf("file %s exists already", o.Path)
			return
		}

		actual[o.Path] = struct{}{}
	}

	assert.Equal(t, expected, actual)
}

func BenchmarkGetFilesFs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fsReaddir(b)
	}
}

func BenchmarkGetFilesOs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osReaddir(b)
	}
}
