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
