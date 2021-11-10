package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func TestStorager(t *testing.T, store types.Storager) {
	suite.Run(t, &StorageSuite{store: store})
}

type StorageSuite struct {
	suite.Suite
	store types.Storager
}

func (s *StorageSuite) TestString() {
	v := s.store.String()
	s.NotEmpty(v, "String() should not be empty.")
}

func (s *StorageSuite) TestMetadata() {
	m := s.store.Metadata()
	s.NotNil(m, "Metadata() should not return nil.")
}

func (s *StorageSuite) TestRead() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Read {
		s.T().Skipf("store doesn't support Delete and Read, skip TestRead.")
	}

	suite.Run(s.T(), &storageReadSuite{p: s})
}

func (s *StorageSuite) TestWrite() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Write {
		s.T().Skipf("store doesn't support Delete and Write, skip TestWrite.")
	}

	suite.Run(s.T(), &storageWriteSuite{p: s})
}

func (s *StorageSuite) TestStat() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Write || !fe.Stat {
		s.T().Skipf("store doesn't support Delete, Write and Stat, skip TestStat.")
	}

	suite.Run(s.T(), &storageStatSuite{p: s})
}

func (s *StorageSuite) TestDelete() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Write {
		s.T().Skipf("store doesn't support Delete, Write, skip TestDelete.")
	}

	suite.Run(s.T(), &storageDeleteSuite{p: s})
}

func (s *StorageSuite) TestList() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Write || !fe.List {
		s.T().Skipf("store doesn't support Delete, Write and List, skip TestList.")
	}

	suite.Run(s.T(), &storageListSuite{p: s})
}

func (s *StorageSuite) TestPath() {
	fe := s.store.Features()

	if !fe.Delete || !fe.Write || !fe.Read {
		s.T().Skipf("store doesn't support Delete, Write and Read, skip TestPath.")
	}

	suite.Run(s.T(), &storagePathSuite{p: s})
}

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

type storageDeleteSuite struct {
	suite.Suite

	p *StorageSuite

	size    int64
	content []byte
	path    string
}

func (s *storageDeleteSuite) SetupTest() {
	var err error

	s.size = rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
	s.content, err = io.ReadAll(io.LimitReader(randbytes.NewRand(), s.size))
	s.NoError(err)

	s.path = uuid.New().String()

	_, err = s.p.store.Write(s.path, bytes.NewReader(s.content), s.size)
	s.NoError(err)
}

func (s *storageDeleteSuite) TestDelete() {
	err := s.p.store.Delete(s.path)
	s.NoError(err)
}

type storageListSuite struct {
	suite.Suite

	p *StorageSuite

	base   string
	length int
	paths  []string
}

func (s *storageListSuite) SetupTest() {
	size := rand.Int63n(256)

	s.length = rand.Intn(16)
	s.base = uuid.NewString()
	s.paths = make([]string, s.length)

	for i := 0; i < s.length; i++ {
		s.paths[i] = fmt.Sprintf("%s/%s", s.base, uuid.NewString())

		_, err := s.p.store.Write(s.paths[i],
			io.LimitReader(randbytes.NewRand(), size), size)
		s.NoError(err)
	}
}

func (s *storageListSuite) TearDownTest() {
	for i := 0; i < s.length; i++ {
		err := s.p.store.Delete(s.paths[i])
		s.NoError(err)
	}
}

func (s *storageListSuite) TestList() {
	it, err := s.p.store.List(s.base, ps.WithListMode(types.ListModeDir))
	s.NoError(err)
	s.NotNil(it)

	paths := make([]string, 0)
	for {
		o, err := it.Next()
		if errors.Is(err, types.IterateDone) {
			break
		}
		s.NoError(err)

		paths = append(paths, o.Path)
	}
	s.ElementsMatch(s.paths, paths)
}

func (s *storageListSuite) TestListWithoutListMode() {
	it, err := s.p.store.List(s.base)
	s.NoError(err)
	s.NotNil(it)

	paths := make([]string, 0)
	for {
		o, err := it.Next()
		if errors.Is(err, types.IterateDone) {
			break
		}
		s.NoError(err)

		paths = append(paths, o.Path)
	}
	s.ElementsMatch(s.paths, paths)
}

func (s *storageListSuite) TestListEmptyDir() {
	path := uuid.New().String()

	it, err := s.p.store.List(path, ps.WithListMode(types.ListModeDir))
	s.NoError(err)
	s.NotNil(it)

	o, err := it.Next()
	s.ErrorIs(err, types.IterateDone)
	s.Nil(o)
}

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
