package tests

import (
	"bytes"
	"crypto/sha256"
	"io"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func testAppender(t *testing.T, store types.Storager, f types.StorageFeatures) {
	if f.CreateAppend && f.Delete {
		Convey("When CreateAppend", func() {
			path := uuid.NewString()
			o, err := store.CreateAppend(path)

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Mode should be appendable", func() {
				// Append object's mode must be appendable.
				So(o.Mode.IsAppend(), ShouldBeTrue)
			})
		})

		Convey("When Delete", func() {
			path := uuid.NewString()
			_, err := store.CreateAppend(path)
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path)
			Convey("The first returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			err = store.Delete(path)
			Convey("The second returned error also should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		if f.WriteAppend {
			Convey("When WriteAppend", func() {
				path := uuid.NewString()
				o, err := store.CreateAppend(path)
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
				content, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
				r := bytes.NewReader(content)

				n, err := store.WriteAppend(o, r, size)

				Convey("WriteAppend error should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("WriteAppend size should be equal to n", func() {
					So(n, ShouldEqual, size)
				})
			})

			if f.CommitAppend {
				Convey("When CommitAppend", func() {
					path := uuid.NewString()
					o, err := store.CreateAppend(path)
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
					content, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))

					_, err = store.WriteAppend(o, bytes.NewReader(content), size)
					if err != nil {
						t.Error(err)
					}

					_, err = store.WriteAppend(o, bytes.NewReader(content), size)
					if err != nil {
						t.Error(err)
					}

					err = store.CommitAppend(o)

					Convey("CommitAppend error should be nil", func() {
						So(err, ShouldBeNil)
					})

					var buf bytes.Buffer
					_, err = store.Read(path, &buf, pairs.WithSize(size*2))

					Convey("Read error should be nil", func() {
						So(err, ShouldBeNil)
					})
					Convey("The content should be match", func() {
						So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(bytes.Repeat(content, 2)))
					})
				})

				Convey("When CreateAppend with an existing object", func() {
					path := uuid.NewString()
					o, err := store.CreateAppend(path)

					defer func() {
						err := store.Delete(path)
						if err != nil {
							t.Error(err)
						}
					}()

					Convey("The first returned error should be nil", func() {
						So(err, ShouldBeNil)
					})

					size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
					r := io.LimitReader(randbytes.NewRand(), size)

					_, err = store.WriteAppend(o, r, size)
					if err != nil {
						t.Fatal(err)
					}

					err = store.CommitAppend(o)
					if err != nil {
						t.Fatal(err)
					}

					o, err = store.CreateAppend(path)

					Convey("The second returned error also should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The Object Mode should be appendable", func() {
						// Append object's mode must be appendable.
						So(o.Mode.IsAppend(), ShouldBeTrue)
					})

					Convey("The object append offset should be 0", func() {
						So(o.MustGetAppendOffset(), ShouldBeZeroValue)
					})
				})
			}
		}
	}
}
