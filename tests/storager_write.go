package tests

import (
	"bytes"
	"io"
	"math/rand"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/pkg/randbytes"
)

type storageWriteSuite struct {
	suite.Suite

	p *StorageSuite

	size    int64
	content []byte
	path    string
}

func (s *storageWriteSuite) SetupTest() {
	var err error

	s.size = rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
	s.content, err = io.ReadAll(io.LimitReader(randbytes.NewRand(), s.size))
	s.NoError(err)

	s.path = uuid.New().String()
}

func (s *storageWriteSuite) TearDownTest() {
	err := s.p.store.Delete(s.path)
	s.NoError(err)
}

func (s *storageWriteSuite) TestWrite() {
	n, err := s.p.store.Write(s.path, bytes.NewReader(s.content), s.size)
	s.NoError(err)
	s.Equal(s.size, n)
}

func (s *storageWriteSuite) TestWriteWithSize() {
	size := rand.Int63n(s.size)

	n, err := s.p.store.Write(s.path, bytes.NewReader(s.content[:size]), size)
	s.NoError(err)
	s.Equal(size, n)
}

func (s *storageWriteSuite) TestWriteWithIoCallback() {
	curWrite := int64(0)
	writeFn := func(bs []byte) {
		curWrite += int64(len(bs))
	}

	n, err := s.p.store.Write(s.path, bytes.NewReader(s.content), s.size,
		ps.WithIoCallback(writeFn))
	s.NoError(err)
	s.Equal(s.size, n, "write size should be equal")
	s.Equal(s.size, curWrite, "io callback should be called")
}

func (s *storageWriteSuite) TestWriteViaNilReader() {
	_, err := s.p.store.Write(s.path, nil, 0)
	s.NoError(err)
}

func (s *storageWriteSuite) TestWriteViaValidReaderAndZeroSize() {
	n, err := s.p.store.Write(s.path, bytes.NewReader(s.content), 0)
	s.NoError(err)
	s.Equal(int64(0), n)
}

func (s *storageWriteSuite) TestWriteViaNilReaderAndValidSize() {
	_, err := s.p.store.Write(s.path, nil, s.size)
	s.Error(err)
}
