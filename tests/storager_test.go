package tests

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/coreutils"
	"github.com/Xuanwo/storage/pkg/randbytes"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

const TestPrefix = "STORAGE_TEST_SERVICE_"

func TestStorager(t *testing.T) {
	services := make(map[string]string)

	content, err := ioutil.ReadFile("storager.yaml")
	if err == nil {
		err = yaml.Unmarshal(content, &services)
		if err != nil {
			t.Error(err)
		}
	}

	env := os.Environ()
	for _, v := range env {
		values := strings.SplitN(v, "=", 2)

		if !strings.HasPrefix(values[0], TestPrefix) {
			continue
		}

		services[values[0]] = values[1]
	}

	for _, v := range services {
		testStorager(t, v)
	}
}

func testStorager(t *testing.T, config string) {
	Convey("Given a basic Storager", t, func() {
		var store storage.Storager
		var err error

		_, store, err = coreutils.Open(config)
		if err != nil {
			t.Error(err)
		}

		Convey("The Storager should not be nil", func() {
			So(store, ShouldNotBeNil)
		})

		Convey("The error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When String called", func() {
			s := store.String()

			Convey("The string should not be empty", func() {
				So(s, ShouldNotBeEmpty)
			})
		})

		Convey("When Metadata called", func() {
			m, err := store.Metadata()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The metadata should not be empty", func() {
				So(m, ShouldNotBeEmpty)
			})
		})

		Convey("When List an empty dir", func() {
			called := false
			fn := types.ObjectFunc(func(_ *types.Object) {
				called = true
			})
			err := store.List("", ps.WithFileFunc(fn))

			Convey("The file func should not be called", func() {
				So(called, ShouldBeFalse)
			})

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When List a dir within files", func() {
			size := int64(4 * 1024 * 1024) // 4MB
			r := io.LimitReader(randbytes.NewRand(), size)
			path := uuid.New().String()
			err := store.Write(path, r, ps.WithSize(size))
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			called := false
			var obj *types.Object
			fn := types.ObjectFunc(func(o *types.Object) {
				called = true
				obj = o

			})
			err = store.List("", ps.WithFileFunc(fn))

			Convey("The file func should be called", func() {
				So(called, ShouldBeTrue)
			})

			Convey("The name and size should be match", func() {
				So(obj, ShouldNotBeNil)
				So(obj.Name, ShouldEqual, path)
				So(obj.Size, ShouldEqual, size)
			})

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
