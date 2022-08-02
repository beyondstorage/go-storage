package tests

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/beyondstorage/go-storage/v5/pkg/randbytes"
)

type storagePathSuite struct {
	suite.Suite

	p *StorageSuite

	base string
	path string
}

func (s *storagePathSuite) SetupTest() {
	s.base = uuid.NewString()
	s.path = uuid.NewString()
}

func (s *storagePathSuite) TearDownTest() {
	path := fmt.Sprintf("%s/%s", s.base, s.path)

	err := s.p.store.Delete(path)
	s.NoError(err)
}

func (s *storagePathSuite) TestAbsPath() {
	m := s.p.store.Metadata()

	path := fmt.Sprintf("%s%s/%s", m.WorkDir, s.base, s.path)

	size := rand.Int63n(4 * 1024 * 1024)
	content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
	s.NoError(err)

	_, err = s.p.store.Write(path, bytes.NewReader(content), size)
	s.NoError(err)

	var buf bytes.Buffer

	n, err := s.p.store.Read(path, &buf)
	s.NoError(err)
	s.Equal(size, n)
}

func (s *storagePathSuite) TestBackslash() {
	path := s.base + "\\" + s.path

	size := rand.Int63n(4 * 1024 * 1024)
	content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
	s.NoError(err)

	_, err = s.p.store.Write(path, bytes.NewReader(content), size)
	s.NoError(err)

	var buf bytes.Buffer

	n, err := s.p.store.Read(path, &buf)
	s.NoError(err)
	s.Equal(size, n)
}
