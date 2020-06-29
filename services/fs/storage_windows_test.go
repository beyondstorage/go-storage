package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Xuanwo/storage/types/info"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
)

func TestStorage_ListDirUnderWindows(t *testing.T) {
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
					ID:         filepath.Join(paths[1], "test_file"),
					Name:       fmt.Sprintf("%s/%s", paths[1], "test_file"),
					Type:       types.ObjectTypeFile,
					Size:       1234,
					UpdatedAt:  time.Unix(1, 0),
					ObjectMeta: info.NewObjectMeta(),
				},
			},
			nil,
		},
	}

	for k, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			client := &Storage{
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
