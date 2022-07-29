package tests

import (
	"bytes"
	"io"
	"math/rand"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/beyondstorage/go-storage/v5/pkg/randbytes"
)

type storageStatSuite struct {
	suite.Suite

	p *StorageSuite

	size    int64
	content []byte
	path    string
}

func (s *storageStatSuite) SetupTest() {
	var err error

	s.size = rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
	s.content, err = io.ReadAll(io.LimitReader(randbytes.NewRand(), s.size))
	s.NoError(err)

	s.path = uuid.New().String()

	_, err = s.p.store.Write(s.path, bytes.NewReader(s.content), s.size)
	s.NoError(err)
}

func (s *storageStatSuite) TearDownTest() {
	err := s.p.store.Delete(s.path)
	s.NoError(err)
}

func (s *storageStatSuite) TestStat() {
	o, err := s.p.store.Stat(s.path)
	s.NoError(err)
	s.NotNil(o)

	osize, ok := o.GetContentLength()
	s.True(ok)
	s.Equal(osize, s.size)
}
