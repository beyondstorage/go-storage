package tests

import (
	"io"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func TestMultiparter(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		f := store.Features()
		if f.CreateMultipart && f.WriteMultipart && f.ListMultipart && f.CompleteMultipart && f.Create && f.Delete && f.Stat {
			Convey("When CreateMultipart", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)

				Convey("The first returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				defer func(multipartID string) {
					err := store.Delete(path, pairs.WithMultipartID(multipartID))
					if err != nil {
						t.Error(err)
					}
				}(o.MustGetMultipartID())

				o, err = store.CreateMultipart(path)

				Convey("The second returned error also should be nil", func() {
					So(err, ShouldBeNil)
				})

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The Object Mode should be part", func() {
					// Multipart object's mode must be Part.
					So(o.Mode.IsPart(), ShouldBeTrue)
					// Multipart object's mode must not be Read.
					So(o.Mode.IsRead(), ShouldBeFalse)
				})

				Convey("The Object must have multipart id", func() {
					// Multipart object must have multipart id.
					_, ok := o.GetMultipartID()
					So(ok, ShouldBeTrue)
				})
			})

			Convey("When Delete with multipart id", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				err = store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				Convey("The first returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				err = store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				Convey("The second returned error also should be nil", func() {
					So(err, ShouldBeNil)
				})
			})

			Convey("When Stat with multipart id", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				multipartId := o.MustGetMultipartID()

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(multipartId))
					if err != nil {
						t.Error(err)
					}
				}()

				mo, err := store.Stat(path, pairs.WithMultipartID(multipartId))

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
					So(mo, ShouldNotBeNil)
				})

				Convey("The Object Mode should be part", func() {
					// Multipart object's mode must be Part.
					So(mo.Mode.IsPart(), ShouldBeTrue)
					// Multipart object's mode must not be Read.
					So(mo.Mode.IsRead(), ShouldBeFalse)
				})

				Convey("The Object must have multipart id", func() {
					// Multipart object must have multipart id.
					mid, ok := mo.GetMultipartID()
					So(ok, ShouldBeTrue)
					So(mid, ShouldEqual, multipartId)
				})
			})

			Convey("When Create with multipart id", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				multipartId := o.MustGetMultipartID()

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(multipartId))
					if err != nil {
						t.Error(err)
					}
				}()

				mo := store.Create(path, pairs.WithMultipartID(multipartId))

				Convey("The Object Mode should be part", func() {
					// Multipart object's mode must be Part.
					So(mo.Mode.IsPart(), ShouldBeTrue)
					// Multipart object's mode must not be Read.
					So(mo.Mode.IsRead(), ShouldBeFalse)
				})

				Convey("The Object must have multipart id", func() {
					// Multipart object must have multipart id.
					mid, ok := mo.GetMultipartID()
					So(ok, ShouldBeTrue)
					So(mid, ShouldEqual, multipartId)
				})
			})

			Convey("When WriteMultipart", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
					if err != nil {
						t.Error(err)
					}
				}()

				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				r := io.LimitReader(randbytes.NewRand(), size)

				n, part, err := store.WriteMultipart(o, r, size, 0)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The part should not be nil", func() {
					So(part, ShouldNotBeNil)
				})

				Convey("The size should be match", func() {
					So(n, ShouldEqual, size)
				})
			})

			Convey("When ListMultiPart", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
					if err != nil {
						t.Error(err)
					}
				}()

				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				partNumber := rand.Intn(1000)        // Choose a random part number from [0, 1000)
				r := io.LimitReader(randbytes.NewRand(), size)

				_, _, err = store.WriteMultipart(o, r, size, partNumber)
				if err != nil {
					t.Error(err)
				}

				it, err := store.ListMultipart(o)

				Convey("ListMultipart error should be nil", func() {
					So(err, ShouldBeNil)
					So(it, ShouldNotBeNil)
				})

				p, err := it.Next()
				Convey("Next error should be nil", func() {
					So(err, ShouldBeNil)
					So(p, ShouldNotBeNil)
				})
				Convey("The part number and size should be match", func() {
					So(p.Index, ShouldEqual, partNumber)
					So(p.Size, ShouldEqual, size)
				})
			})

			Convey("When List with part type", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				defer func() {
					err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
					if err != nil {
						t.Error(err)
					}
				}()

				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				partNumber := rand.Intn(1000)        // Choose a random part number from [0, 1000)
				r := io.LimitReader(randbytes.NewRand(), size)

				_, _, err = store.WriteMultipart(o, r, size, partNumber)
				if err != nil {
					t.Error(err)
				}

				it, err := store.List("", pairs.WithListMode(types.ListModePart))
				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("The iterator should not be nil", func() {
					So(it, ShouldNotBeNil)
				})

				mo, err := it.Next()
				Convey("Next error should be nil", func() {
					So(err, ShouldBeNil)
					So(mo, ShouldNotBeNil)
				})
				Convey("The path and multipart id should be match", func() {
					So(mo.Path, ShouldEqual, path)
					So(mo.Mode.IsPart(), ShouldBeTrue)

					// Multipart object must have multipart id.
					mid, ok := mo.GetMultipartID()
					So(ok, ShouldBeTrue)
					So(mid, ShouldEqual, o.MustGetMultipartID())
				})
			})

			Convey("When CompletePart", func() {
				path := uuid.New().String()
				o, err := store.CreateMultipart(path)
				if err != nil {
					t.Error(err)
				}

				defer func() {
					err := store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				// Set 0 to `partNumber` here as the part numbers must be continuous for `CompleteMultipartUpload` in `cos` which is different with other storages.
				partNumber := 0
				r := io.LimitReader(randbytes.NewRand(), size)

				_, part, err := store.WriteMultipart(o, r, size, partNumber)
				if err != nil {
					t.Error(err)
				}

				err = store.CompleteMultipart(o, []*types.Part{part})

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object should be readable after complete", func() {
					ro, err := store.Stat(path)

					So(err, ShouldBeNil)
					So(ro.Mode.IsRead(), ShouldBeTrue)
					So(ro.Mode.IsPart(), ShouldBeFalse)
				})
			})
		}
	})
}
