package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/Xuanwo/storage/types/metadata"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
)

func TestNewClient(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func TestClient_String(t *testing.T) {
	c := Storage{}
	err := c.Init(pairs.WithWorkDir("/test"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "posixfs Storager {WorkDir: /test}", c.String())
}

func TestClient_Init(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		client := Storage{}
		err := client.Init()
		assert.Error(t, err)
		assert.Equal(t, "", client.workDir)
	})

	t.Run("with workDir", func(t *testing.T) {
		client := Storage{}
		err := client.Init(pairs.WithWorkDir("test"))
		assert.NoError(t, err)
		assert.Equal(t, "test", client.workDir)
	})
}

func TestClient_Metadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	{
		client := Storage{workDir: "/test"}

		m, err := client.Metadata()
		assert.NoError(t, err)
		gotBase, ok := m.GetWorkDir()
		assert.True(t, true, ok)
		assert.Equal(t, "/test", gotBase)
	}
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() os.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.modTime
}

func (f fileInfo) IsDir() bool {
	return f.mode.IsDir()
}

func (f fileInfo) Sys() interface{} {
	return f
}

func TestClient_Stat(t *testing.T) {
	nowTime := time.Now()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
		file fileInfo

		object *types.Object
	}{
		{
			"regular file",
			nil,
			fileInfo{
				name:    "regular file",
				size:    1234,
				mode:    0777,
				modTime: nowTime,
			},
			&types.Object{
				Name:      "regular file",
				Type:      types.ObjectTypeFile,
				Size:      1234,
				UpdatedAt: nowTime,
				Metadata:  metadata.Metadata{},
			},
		},
		{
			"dir",
			nil,
			fileInfo{
				name:    "dir",
				size:    0,
				mode:    os.ModeDir | 0777,
				modTime: nowTime,
			},
			&types.Object{
				Name:      "dir",
				Type:      types.ObjectTypeDir,
				UpdatedAt: nowTime,
				Metadata:  make(metadata.Metadata),
			},
		},
		{
			"stream",
			nil,
			fileInfo{
				name:    "stream",
				size:    0,
				mode:    os.ModeDevice | 0777,
				modTime: nowTime,
			},
			&types.Object{
				Name:      "stream",
				Type:      types.ObjectTypeStream,
				UpdatedAt: nowTime,
				Metadata:  make(metadata.Metadata),
			},
		},
		{
			"-",
			nil,
			fileInfo{},
			&types.Object{
				Name:     "-",
				Type:     types.ObjectTypeStream,
				Metadata: make(metadata.Metadata),
			},
		},
		{
			"invalid",
			nil,
			fileInfo{
				name:    "invalid",
				size:    0,
				mode:    os.ModeIrregular | 0777,
				modTime: nowTime,
			},
			&types.Object{
				Name:      "invalid",
				Type:      types.ObjectTypeInvalid,
				UpdatedAt: nowTime,
				Metadata:  make(metadata.Metadata),
			},
		},
		{
			"error",
			&os.PathError{
				Op:   "stat",
				Path: "invalid",
				Err:  os.ErrPermission,
			},
			fileInfo{}, nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				osStat: func(name string) (os.FileInfo, error) {
					assert.Equal(t, v.name, name)
					return v.file, v.err
				},
			}
			o, err := client.Stat(v.name)
			assert.Equal(t, v.err == nil, err == nil)
			if v.object != nil {
				assert.NotNil(t, o)
				assert.EqualValues(t, v.object, o)
			} else {
				assert.Nil(t, o)
			}
		})
	}
}

func TestClient_WriteStream(t *testing.T) {
	err := os.Remove("/tmp/test")
	var e *os.PathError
	if errors.As(err, &e) {
		t.Logf("%#v", e)
	}
}

func TestClient_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		err  error
	}{
		{"delete file", nil},
		{"delete nonempty dir", &os.PathError{
			Op:   "remove",
			Path: "delete nonempty dir",
			Err:  errors.New("remove fail"),
		}},
	}

	for _, v := range tests {
		v := v

		t.Run(v.name, func(t *testing.T) {

			client := Storage{
				osRemove: func(name string) error {
					assert.Equal(t, v.name, name)
					return v.err
				},
			}
			err := client.Delete(v.name)
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestClient_Copy(t *testing.T) {
	t.Run("Failed at open source file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: errors.New("path error"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at open dst file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: errors.New("open fail"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at io.CopyBuffer", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return nil, nil
			},
			ioCopyBuffer: func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, io.ErrShortWrite
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("All successful", func(t *testing.T) {
		fakeFile := &os.File{}
		// Monkey patch the file's Close.
		monkey.PatchInstanceMethod(reflect.TypeOf(fakeFile), "Close",
			func(f *os.File) error {
				return nil
			})

		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Storage{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return fakeFile, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return fakeFile, nil
			},
			ioCopyBuffer: func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, nil
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Copy(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestClient_Move(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Storage{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return &os.LinkError{
					Op:  "rename",
					Old: oldpath,
					New: newpath,
					Err: errors.New("rename fail"),
				}
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Move(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Storage{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return nil
			},
			osMkdirAll: func(path string, perm os.FileMode) error {
				return nil
			},
		}

		err := client.Move(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestClient_Reach(t *testing.T) {
	client := Storage{}

	assert.Panics(t, func() {
		_, _ = client.Reach(uuid.New().String())
	})
}

func TestClient_ListDir(t *testing.T) {
	paths := make([]string, 100)
	for k := range paths {
		paths[k] = uuid.New().String()
	}

	tests := []struct {
		name  string
		fi    []os.FileInfo
		items []*types.Object
		err   error
	}{
		{
			"success file",
			[]os.FileInfo{
				fileInfo{
					name:    "test_file",
					size:    1234,
					mode:    0644,
					modTime: time.Unix(1, 0),
				},
			},
			[]*types.Object{
				{
					Name:      filepath.Join(paths[0], "test_file"),
					Type:      types.ObjectTypeFile,
					Size:      1234,
					UpdatedAt: time.Unix(1, 0),
					Metadata:  metadata.Metadata{},
				},
			},
			nil,
		},
		{
			"success file recursively",
			[]os.FileInfo{
				fileInfo{
					name:    "test_file",
					size:    1234,
					mode:    0644,
					modTime: time.Unix(1, 0),
				},
			},
			[]*types.Object{
				{
					Name:      filepath.Join(paths[1], "test_file"),
					Type:      types.ObjectTypeFile,
					Size:      1234,
					UpdatedAt: time.Unix(1, 0),
					Metadata:  metadata.Metadata{},
				},
			},
			nil,
		},
		{
			"success dir",
			[]os.FileInfo{
				fileInfo{
					name:    "test_dir",
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: time.Unix(1, 0),
				},
			},
			[]*types.Object{
				{
					Name:      filepath.Join(paths[2], "test_dir"),
					Type:      types.ObjectTypeDir,
					Size:      0,
					UpdatedAt: time.Unix(1, 0),
					Metadata:  metadata.Metadata{},
				},
			},
			nil,
		},
		{
			"success dir recursively",
			[]os.FileInfo{
				fileInfo{
					name:    "test_dir",
					size:    0,
					mode:    os.ModeDir | 0755,
					modTime: time.Unix(1, 0),
				},
			},
			[]*types.Object{
				{
					Name:      filepath.Join(paths[3], "test_dir"),
					Type:      types.ObjectTypeDir,
					Size:      0,
					UpdatedAt: time.Unix(1, 0),
					Metadata:  metadata.Metadata{},
				},
			},
			nil,
		},
		{
			"os error",
			nil,
			[]*types.Object{},
			&os.PathError{Op: "readdir", Path: "", Err: errors.New("readdir fail")},
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				ioutilReadDir: func(dirname string) (infos []os.FileInfo, e error) {
					assert.Equal(t, paths[k], dirname)
					return v.fi, v.err
				},
			}

			items := make([]*types.Object, 0)

			err := client.ListDir(paths[k], pairs.WithDirFunc(func(object *types.Object) {
				items = append(items, object)
			}), pairs.WithFileFunc(func(object *types.Object) {
				items = append(items, object)
			}))
			assert.Equal(t, v.err == nil, err == nil)
			assert.EqualValues(t, v.items, items)
		})
	}
}

func TestClient_Read(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pairs   []*types.Pair
		isNil   bool
		openErr error
		seekErr error
	}{
		{
			"success",
			"test_success",
			nil,
			false,
			nil,
			nil,
		},
		{
			"error",
			"test_error",
			nil,
			true,
			&os.PathError{Op: "readdir", Path: "", Err: errors.New("readdir fail")},
			nil,
		},
		{
			"stdin",
			"-",
			nil,
			false,
			nil,
			nil,
		},
		{
			"stdin with size",
			"-",
			[]*types.Pair{
				pairs.WithSize(100),
			},
			false,
			nil,
			nil,
		},
		{
			"success with size",
			"test_success",
			[]*types.Pair{
				pairs.WithSize(100),
			},
			false,
			nil,
			nil,
		},
		{
			"success with offset",
			"test_success",
			[]*types.Pair{
				pairs.WithOffset(10),
			},
			false,
			nil,
			nil,
		},
		{
			"error with offset",
			"test_success",
			[]*types.Pair{
				pairs.WithOffset(10),
			},
			true,
			nil,
			io.ErrUnexpectedEOF,
		},
		{
			"success with and size offset",
			"test_success",
			[]*types.Pair{
				pairs.WithSize(100),
				pairs.WithOffset(10),
			},
			false,
			nil,
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			fakeFile := &os.File{}
			monkey.PatchInstanceMethod(reflect.TypeOf(fakeFile), "Seek", func(f *os.File, offset int64, whence int) (ret int64, err error) {
				t.Logf("Seek has been called.")
				assert.Equal(t, int64(10), offset)
				assert.Equal(t, 0, whence)
				return 0, v.seekErr
			})

			client := Storage{
				osOpen: func(name string) (file *os.File, e error) {
					assert.Equal(t, v.path, name)
					return fakeFile, v.openErr
				},
			}

			o, err := client.Read(v.path, v.pairs...)
			assert.Equal(t, v.openErr == nil && v.seekErr == nil, err == nil)
			assert.Equal(t, v.isNil, o == nil)
		})
	}
}

func TestClient_Write(t *testing.T) {
	paths := make([]string, 10)
	for k := range paths {
		paths[k] = uuid.New().String()
	}

	tests := []struct {
		name         string
		osCreate     func(name string) (*os.File, error)
		ioCopyN      func(dst io.Writer, src io.Reader, n int64) (written int64, err error)
		ioCopyBuffer func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)
		hasErr       bool
	}{
		{
			"failed os create",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[0], name)
				return nil, &os.PathError{
					Op:   "open",
					Path: "",
					Err:  os.ErrNotExist,
				}
			},
			nil,
			nil,
			true,
		},
		{
			"failed io copyn",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[1], name)
				return &os.File{}, nil
			},
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				return 0, io.EOF
			},
			nil,
			true,
		},
		{
			"failed io copy buffer",
			nil,
			nil,
			func(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
				return 0, io.EOF
			},
			true,
		},
		{
			"success with size",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[3], name)
				return &os.File{}, nil
			},
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				assert.Equal(t, int64(1234), n)
				return 0, nil
			},
			nil,
			false,
		},
		{
			"success with stdout",
			nil,
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				assert.Equal(t, int64(1234), n)
				return 0, nil
			},
			nil,
			false,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Storage{
				osCreate:     v.osCreate,
				ioCopyN:      v.ioCopyN,
				ioCopyBuffer: v.ioCopyBuffer,
				osMkdirAll: func(path string, perm os.FileMode) error {
					return nil
				},
			}

			var pair []*types.Pair
			if v.ioCopyN != nil {
				pair = append(pair, pairs.WithSize(1234))
			}

			var err error
			if v.osCreate == nil {
				err = client.Write("-", nil, pair...)
			} else {
				err = client.Write(paths[k], nil, pair...)
			}
			assert.Equal(t, v.hasErr, err != nil)
		})
	}
}
