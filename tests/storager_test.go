package tests

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/coreutils"
	"github.com/Xuanwo/storage/pkg/randbytes"
	"github.com/Xuanwo/storage/services"
	"github.com/Xuanwo/storage/types"
	ps "github.com/Xuanwo/storage/types/pairs"
)

func TestStorager(t *testing.T) {
	srv := loadConfig()

	for _, v := range srv {
		pairs, err := ps.Parse(v.Options)
		if err != nil {
			t.Error(err)
		}
		println("Start test for ", v.Type)
		testStorager(t, v.Type, pairs)
		testDirLister(t, v.Type, pairs)
	}
}

func testStorager(t *testing.T, typ string, pair []*types.Pair) {
	Convey("Given a basic Storager", t, func() {
		var store storage.Storager
		var err error

		store, err = coreutils.OpenStorager(typ, pair...)
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

		Convey("When Read a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			err = store.Write(path, bytes.NewReader(content), ps.WithSize(size))
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			r, err := store.Read(path)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The content should be match", func() {
				So(r, ShouldNotBeNil)

				readContent, err := ioutil.ReadAll(r)
				if err != nil {
					t.Error(err)
				}
				So(readContent, ShouldResemble, content)
			})

		})

		Convey("When Write a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)
			path := uuid.New().String()

			err := store.Write(path, r, ps.WithSize(size))

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get Object without error", func() {
				o, err := store.Stat(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The name and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Name, ShouldEqual, path)
					So(o.Size, ShouldEqual, size)
				})
			})

			Convey("Read should get Object data without error", func() {
				r, err := store.Read(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The reader should not be nil", func() {
					So(r, ShouldNotBeNil)
				})
			})

		})

		Convey("When Stat a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			err = store.Write(path, bytes.NewReader(content), ps.WithSize(size))
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			o, err := store.Stat(path)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object name and size should be match", func() {
				So(o, ShouldNotBeNil)
				So(o.Name, ShouldEqual, path)
				So(o.Size, ShouldEqual, size)
			})
		})

		Convey("When Delete a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, err := ioutil.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			err = store.Write(path, bytes.NewReader(content), ps.WithSize(size))
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get nil Object and ObjectNotFound error", func() {
				o, err := store.Stat(path)

				So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				So(o, ShouldBeNil)
			})
		})
	})
}

func testDirLister(t *testing.T, typ string, pair []*types.Pair) {
	Convey("Given a dir lister", t, func() {
		var store storage.Storager
		var lister storage.DirLister
		var err error

		store, err = coreutils.OpenStorager(typ, pair...)
		if err != nil {
			t.Error(err)
		}

		lister, ok := store.(storage.DirLister)
		if !ok {
			t.Skip()
		}

		Convey("When List an empty dir", func() {
			called := false
			fn := types.ObjectFunc(func(_ *types.Object) {
				called = true
			})

			err := lister.ListDir("", ps.WithFileFunc(fn))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The file func should not be called", func() {
				So(called, ShouldBeFalse)
			})

		})

		Convey("When List a dir within files", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

			err = lister.ListDir("", ps.WithFileFunc(fn))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The file func should be called", func() {
				So(called, ShouldBeTrue)
			})

			Convey("The name and size should be match", func() {
				So(obj, ShouldNotBeNil)
				So(obj.Name, ShouldEqual, path)
				So(obj.Size, ShouldEqual, size)
			})
		})

	})
}
