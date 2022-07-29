package tests

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/pkg/randbytes"
	"github.com/beyondstorage/go-storage/v5/services"
	"github.com/beyondstorage/go-storage/v5/types"
)

func TestMover(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Move a file", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			src := uuid.New().String()

			_, err := store.Write(src, bytes.NewReader(content), size)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(src)
				if err != nil {
					t.Error(err)
				}
			}()

			dst := uuid.New().String()
			err = store.Move(src, dst)

			defer func() {
				err = store.Delete(dst)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get src object not exist", func() {
				_, err := store.Stat(src)

				Convey("The error should be ErrObjectNotExist", func() {
					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				})
			})

			Convey("Read should get dst object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(dst, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					So(buf, ShouldNotBeNil)
					So(n, ShouldEqual, size)
					So(md5.Sum(buf.Bytes()), ShouldResemble, md5.Sum(content))
				})
			})
		})

		Convey("When Move to an existing file", func() {
			srcSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			content, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), srcSize))
			src := uuid.New().String()

			_, err := store.Write(src, bytes.NewReader(content), srcSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(src)
				if err != nil {
					t.Error(err)
				}
			}()

			dstSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), dstSize)
			dst := uuid.New().String()

			_, err = store.Write(dst, r, dstSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(dst)
				if err != nil {
					t.Error(err)
				}
			}()

			err = store.Move(src, dst)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get src object not exist", func() {
				_, err := store.Stat(src)

				Convey("The error should be ErrObjectNotExist", func() {
					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				})
			})

			Convey("Read should get dst object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(dst, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					So(buf, ShouldNotBeNil)
					So(n, ShouldEqual, srcSize)
					So(md5.Sum(buf.Bytes()), ShouldResemble, md5.Sum(content))
				})
			})
		})
	})
}

func TestMoverWithDir(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Move to an existing dir", func() {

			srcSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), srcSize)
			src := uuid.New().String()

			_, err := store.Write(src, r, srcSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(src)
				if err != nil {
					t.Error(err)
				}
			}()

			dst := uuid.New().String()
			_, err = store.CreateDir(dst)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(dst, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			err = store.Move(src, dst)
			Convey("The error should be ErrObjectModeInvalid", func() {
				So(errors.Is(err, services.ErrObjectModeInvalid), ShouldBeTrue)
			})
		})
	})
}

func TestMoverWithVirtualDir(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Move to an existing dir", func() {

			srcSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), srcSize)
			src := uuid.New().String()

			_, err := store.Write(src, r, srcSize)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(src)
				if err != nil {
					t.Error(err)
				}
			}()

			dst := uuid.New().String()
			_, err = store.CreateDir(dst)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = store.Delete(dst, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			err = store.Move(src, dst)

			defer func() {
				err = store.Delete(dst)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Stat should get dst object without error", func() {
				o, err := store.Stat(dst)

				So(err, ShouldBeNil)
				So(o, ShouldNotBeNil)

				Convey("The Object Mode should be read", func() {
					So(o.Mode.IsRead(), ShouldBeTrue)
				})

				Convey("The path and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Path, ShouldEqual, dst)

					osize, ok := o.GetContentLength()
					So(ok, ShouldBeTrue)
					So(osize, ShouldEqual, srcSize)
				})
			})
		})
	})
}
