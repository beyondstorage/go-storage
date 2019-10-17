package posixfs

import (
	"errors"
	"io"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.NotNil(t, c)
}

func TestClient_Metadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	{
		client := Client{}

		m, err := client.Metadata()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(m))
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
				modTime: time.Now(),
			},
			&types.Object{
				Name: "regular file",
				Type: types.ObjectTypeFile,
				Metadata: types.Metadata{
					types.Size: int64(1234),
				},
			},
		},
		{
			"dir",
			nil,
			fileInfo{
				name:    "dir",
				size:    0,
				mode:    os.ModeDir | 0777,
				modTime: time.Now(),
			},
			&types.Object{
				Name:     "dir",
				Type:     types.ObjectTypeDir,
				Metadata: make(types.Metadata),
			},
		},
		{
			"stream",
			nil,
			fileInfo{
				name:    "stream",
				size:    0,
				mode:    os.ModeDevice | 0777,
				modTime: time.Now(),
			},
			&types.Object{
				Name:     "stream",
				Type:     types.ObjectTypeStream,
				Metadata: make(types.Metadata),
			},
		},
		{
			"invalid",
			nil,
			fileInfo{
				name:    "invalid",
				size:    0,
				mode:    os.ModeIrregular | 0777,
				modTime: time.Now(),
			},
			&types.Object{
				Name:     "invalid",
				Type:     types.ObjectTypeInvalid,
				Metadata: make(types.Metadata),
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
			client := Client{
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
		name      string
		err       error
		recursive bool
	}{
		{"delete file", nil, false},
		{"delete dir", nil, true},
		{"delete nonempty dir", &os.PathError{
			Op:   "remove",
			Path: "delete nonempty dir",
			Err:  syscall.ENOTEMPTY,
		}, false},
	}

	for _, v := range tests {
		v := v

		t.Run(v.name, func(t *testing.T) {

			client := Client{
				osRemove: func(name string) error {
					assert.Equal(t, v.name, name)
					assert.False(t, v.recursive)
					return v.err
				},
				osRemoveAll: func(name string) error {
					assert.Equal(t, v.name, name)
					assert.True(t, v.recursive)
					return v.err
				},
			}
			pairs := make([]*types.Pair, 0)
			if v.recursive {
				pairs = append(pairs, types.WithRecursive(true))
			}
			err := client.Delete(v.name, pairs...)
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestClient_Copy(t *testing.T) {
	t.Run("Failed at open source file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Client{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: syscall.ENONET,
				}
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at open dst file", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Client{
			osOpen: func(name string) (file *os.File, e error) {
				assert.Equal(t, srcName, name)
				return nil, nil
			},
			osCreate: func(name string) (file *os.File, e error) {
				assert.Equal(t, dstName, name)
				return nil, &os.PathError{
					Op:  "open",
					Err: syscall.EEXIST,
				}
			},
		}

		err := client.Copy(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("Failed at io.CopyBuffer", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()
		client := Client{
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
		client := Client{
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
		}

		err := client.Copy(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestClient_Move(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Client{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return &os.LinkError{
					Op:  "rename",
					Old: oldpath,
					New: newpath,
					Err: syscall.EISDIR,
				}
			},
		}

		err := client.Move(srcName, dstName)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		srcName := uuid.New().String()
		dstName := uuid.New().String()

		client := Client{
			osRename: func(oldpath, newpath string) error {
				assert.Equal(t, srcName, oldpath)
				assert.Equal(t, dstName, newpath)
				return nil
			},
		}

		err := client.Move(srcName, dstName)
		assert.NoError(t, err)
	})
}

func TestClient_Reach(t *testing.T) {
	client := Client{}

	assert.Panics(t, func() {
		_, _ = client.Reach(uuid.New().String())
	})
}

func TestClient_CreateDir(t *testing.T) {
	paths := make([]string, 10)
	for k := range paths {
		paths[k] = uuid.New().String()
	}
	tests := []struct {
		name string
		err  error
	}{
		{
			"error",
			&os.PathError{Op: "mkdir", Path: paths[0], Err: syscall.ENOTDIR},
		},
		{
			"success",
			nil,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Client{
				osMkdirAll: func(path string, perm os.FileMode) error {
					assert.Equal(t, paths[k], path)
					assert.Equal(t, os.FileMode(0755), perm)
					return v.err
				},
			}

			err := client.CreateDir(paths[k])
			assert.Equal(t, v.err == nil, err == nil)
		})
	}
}

func TestClient_ListDir(t *testing.T) {
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
					Name: "test_file",
					Type: types.ObjectTypeFile,
					Metadata: types.Metadata{
						types.Size:      int64(1234),
						types.UpdatedAt: time.Unix(1, 0),
					},
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
					Name: "test_dir",
					Type: types.ObjectTypeDir,
					Metadata: types.Metadata{
						types.Size:      int64(0),
						types.UpdatedAt: time.Unix(1, 0),
					},
				},
			},
			nil,
		},
		{
			"error",
			nil,
			nil,
			&os.PathError{Op: "readdir", Path: "", Err: syscall.ENOTDIR},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			path := uuid.New().String()

			client := Client{
				ioutilReadDir: func(dirname string) (infos []os.FileInfo, e error) {
					assert.Equal(t, path, dirname)
					return v.fi, v.err
				},
			}

			x := client.ListDir(path)
			for _, expectItem := range v.items {
				item, err := x.Next()
				if v.err != nil {
					assert.Error(t, err)
					assert.True(t, errors.Is(err, v.err))
				}
				assert.NotNil(t, item)
				assert.EqualValues(t, expectItem, item)
			}
			if len(v.items) == 0 {
				item, err := x.Next()
				if v.err != nil {
					assert.Error(t, err)
				} else {
					assert.True(t, errors.Is(err, iterator.ErrDone))
				}
				assert.Nil(t, item)
			}
		})
	}
}

func TestClient_Read(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		isNil bool
		err   error
	}{
		{
			"success",
			"test_success",
			false,
			nil,
		},
		{
			"error",
			"test_error",
			true,
			&os.PathError{Op: "readdir", Path: "", Err: syscall.ENOTDIR},
		},
		{
			"stdin",
			"-",
			false,
			nil,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			fakeFile := &os.File{}

			client := Client{
				osOpen: func(name string) (file *os.File, e error) {
					assert.Equal(t, v.path, name)
					if v.isNil {
						return nil, v.err
					}
					return fakeFile, v.err
				},
			}

			o, err := client.Read(v.path)
			assert.Equal(t, v.err == nil, err == nil)
			assert.Equal(t, v.isNil, o == nil)
		})
	}
}

func TestClient_WriteFile(t *testing.T) {
	paths := make([]string, 10)
	for k := range paths {
		paths[k] = uuid.New().String()
	}

	tests := []struct {
		name     string
		osCreate func(name string) (*os.File, error)
		ioCopyN  func(dst io.Writer, src io.Reader, n int64) (written int64, err error)
		hasErr   bool
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
			true,
		},
		{
			"success",
			func(name string) (file *os.File, e error) {
				assert.Equal(t, paths[2], name)
				return &os.File{}, nil
			},
			func(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
				return 0, nil
			},
			false,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := Client{
				osCreate: v.osCreate,
				ioCopyN:  v.ioCopyN,
			}

			err := client.WriteFile(paths[k], 1234, nil)
			assert.Equal(t, v.hasErr, err != nil)
		})
	}
}
