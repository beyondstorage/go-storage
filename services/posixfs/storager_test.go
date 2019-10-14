package posixfs

import (
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
			fn := func(name string) (os.FileInfo, error) {
				assert.Equal(t, v.name, name)
				return v.file, v.err
			}
			monkey.Patch(os.Stat, fn)

			client := Client{}
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
