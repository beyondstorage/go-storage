package tests

import (
	"io"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func TestLinker(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		f := store.Features()
		if f.CreateLink && f.Write && f.Delete && f.Stat {
			workDir := store.Metadata().WorkDir

			Convey("When create a link object", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				r := io.LimitReader(randbytes.NewRand(), size)
				target := uuid.New().String()

				_, err := store.Write(target, r, size)
				if err != nil {
					t.Fatal(err)
				}

				defer func() {
					err = store.Delete(target)
					if err != nil {
						t.Error(err)
					}
				}()

				path := uuid.New().String()
				o, err := store.CreateLink(path, target)

				defer func() {
					err = store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object mode should be link", func() {
					// Link object's mode must be link.
					So(o.Mode.IsLink(), ShouldBeTrue)
				})

				Convey("The linkTarget of the object must be the same as the target", func() {
					// The linkTarget must be the same as the target.
					linkTarget, ok := o.GetLinkTarget()

					So(ok, ShouldBeTrue)
					So(linkTarget, ShouldEqual, filepath.Join(workDir, target))
				})

				Convey("Stat should get path object without error", func() {
					obj, err := store.Stat(path)

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The object mode should be link", func() {
						// Link object's mode must be link.
						So(obj.Mode.IsLink(), ShouldBeTrue)
					})

					Convey("The linkTarget of the object must be the same as the target", func() {
						// The linkTarget must be the same as the target.
						linkTarget, ok := obj.GetLinkTarget()

						So(ok, ShouldBeTrue)
						So(linkTarget, ShouldEqual, filepath.Join(workDir, target))
					})
				})
			})

			Convey("When create a link object from a not existing target", func() {
				target := uuid.New().String()

				path := uuid.New().String()
				o, err := store.CreateLink(path, target)

				defer func() {
					err = store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object mode should be link", func() {
					// Link object's mode must be link.
					So(o.Mode.IsLink(), ShouldBeTrue)
				})

				Convey("The linkTarget of the object must be the same as the target", func() {
					linkTarget, ok := o.GetLinkTarget()

					So(ok, ShouldBeTrue)
					So(linkTarget, ShouldEqual, filepath.Join(workDir, target))
				})

				Convey("Stat should get path object without error", func() {
					obj, err := store.Stat(path)

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The object mode should be link", func() {
						// Link object's mode must be link.
						So(obj.Mode.IsLink(), ShouldBeTrue)
					})

					Convey("The linkTarget of the object must be the same as the target", func() {
						// The linkTarget must be the same as the target.
						linkTarget, ok := obj.GetLinkTarget()

						So(ok, ShouldBeTrue)
						So(linkTarget, ShouldEqual, filepath.Join(workDir, target))
					})
				})
			})

			Convey("When CreateLink to an existing path", func() {
				firstSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				firstR := io.LimitReader(randbytes.NewRand(), firstSize)
				firstTarget := uuid.New().String()

				_, err := store.Write(firstTarget, firstR, firstSize)
				if err != nil {
					t.Fatal(err)
				}

				defer func() {
					err = store.Delete(firstTarget)
					if err != nil {
						t.Error(err)
					}
				}()

				path := uuid.New().String()
				o, err := store.CreateLink(path, firstTarget)

				defer func() {
					err = store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The first returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				secondSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				secondR := io.LimitReader(randbytes.NewRand(), secondSize)
				secondTarget := uuid.New().String()

				_, err = store.Write(secondTarget, secondR, secondSize)
				if err != nil {
					t.Fatal(err)
				}

				defer func() {
					err = store.Delete(secondTarget)
					if err != nil {
						t.Error(err)
					}
				}()

				o, err = store.CreateLink(path, secondTarget)

				Convey("The second returned error should also be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The object mode should be link", func() {
					// Link object's mode must be link.
					So(o.Mode.IsLink(), ShouldBeTrue)
				})

				Convey("The linkTarget of the object must be the same as the secondTarget", func() {
					// The linkTarget must be the same as the secondTarget.
					linkTarget, ok := o.GetLinkTarget()

					So(ok, ShouldBeTrue)
					So(linkTarget, ShouldEqual, filepath.Join(workDir, secondTarget))
				})
			})
		}
	})
}
