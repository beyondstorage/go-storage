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

type storageReadSuite struct {
	suite.Suite

	p *StorageSuite

	size    int64
	content []byte
	path    string
}

func (s *storageReadSuite) SetupTest() {
	var err error

	s.size = rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
	s.content, err = io.ReadAll(io.LimitReader(randbytes.NewRand(), s.size))
	s.NoError(err)

	s.path = uuid.New().String()

	_, err = s.p.store.Write(s.path, bytes.NewReader(s.content), s.size)
	s.NoError(err)
}

func (s *storageReadSuite) TearDownTest() {
	err := s.p.store.Delete(s.path)
	s.NoError(err)
}

func (s *storageReadSuite) TestRead() {
	var buf bytes.Buffer

	n, err := s.p.store.Read(s.path, &buf)
	s.NoError(err)
	s.Equal(s.size, n, "size should equal")
	s.EqualValues(s.content, buf.Bytes(), "content should equal")
}

func (s *storageReadSuite) TestReadWithIoCallback() {
	curRead := int64(0)
	readFn := func(bs []byte) {
		curRead += int64(len(bs))
	}

	var buf bytes.Buffer

	n, err := s.p.store.Read(s.path, &buf, ps.WithIoCallback(readFn))
	s.NoError(err)
	s.Equal(s.size, n, "size should equal")
	s.Equal(s.size, curRead, "io callback should be called")
	s.EqualValues(s.content, buf.Bytes(), "content should equal")
}

func (s *storageReadSuite) TestReadWithOffset() {
	offset := rand.Int63n(s.size)

	var buf bytes.Buffer

	n, err := s.p.store.Read(s.path, &buf, ps.WithOffset(offset))
	s.NoError(err)
	s.Equal(s.size-offset, n, "size should equal")
	s.EqualValues(s.content[offset:], buf.Bytes(), "content should equal")
}

func (s *storageReadSuite) TestReadWithSize() {
	length := rand.Int63n(s.size)

	var buf bytes.Buffer

	n, err := s.p.store.Read(s.path, &buf, ps.WithSize(length))
	s.NoError(err)
	s.Equal(length, n, "size should equal")
	s.EqualValues(s.content[:length], buf.Bytes(), "content should equal")
}

func (s *storageReadSuite) TestReadWithSizeAndOffset() {
	offset := rand.Int63n(s.size)
	length := rand.Int63n(s.size - offset)

	var buf bytes.Buffer

	n, err := s.p.store.Read(s.path, &buf, ps.WithOffset(offset), ps.WithSize(length))
	s.NoError(err)
	s.Equal(length, n, "size should equal")
	s.EqualValues(s.content[offset:offset+length], buf.Bytes(), "content should equal")
}
