package tests

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/types"
)

func TestMultipartHTTPSigner(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		signer, ok := store.(types.MultipartHTTPSigner)
		So(ok, ShouldBeTrue)

		Convey("When CreateMultipart via QuerySignHTTPCreateMultipart", func() {
			path := uuid.New().String()
			req, err := signer.QuerySignHTTPCreateMultipart(path, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			_, err = client.Do(req)

			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("List with ModePart should get the object without error", func() {
				it, err := store.List(path, pairs.WithListMode(types.ListModePart))

				So(err, ShouldBeNil)

				o, err := it.Next()
				So(err, ShouldBeNil)
				So(o, ShouldNotBeNil)
				So(o.Path, ShouldEqual, path)
			})

			defer func() {
				it, err := store.List(path, pairs.WithListMode(types.ListModePart))
				if err != nil {
					t.Error(err)
				}

				o, err := it.Next()
				if err != nil {
					t.Error(err)
				}

				err = store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				if err != nil {
					t.Error(err)
				}
			}()
		})

		Convey("When WriteMultipart via QuerySignHTTPWriteMultipart", func() {
			path := uuid.New().String()
			o, err := store.(types.Multiparter).CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path, pairs.WithMultipartID(o.MustGetMultipartID()))
				if err != nil {
					t.Error(err)
				}
			}()

			size := rand.Int63n(4 * 1024 * 1024)
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			req, err := signer.QuerySignHTTPWriteMultipart(o, size, 0, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			req.Body = io.NopCloser(bytes.NewReader(content))

			client := http.Client{}
			resp, err := client.Do(req)

			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})

			Convey("The size should be match", func() {
				So(resp.Request.ContentLength, ShouldEqual, size)
			})
		})

		Convey("When ListMultiPart via QuerySignHTTPListMultiPart", func() {
			mu, ok := store.(types.Multiparter)
			So(ok, ShouldBeTrue)

			path := uuid.New().String()
			o, err := mu.CreateMultipart(path)
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

			_, _, err = mu.WriteMultipart(o, r, size, partNumber)
			if err != nil {
				t.Error(err)
			}

			req, err := signer.QuerySignHTTPListMultipart(o, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			_, err = client.Do(req)

			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When CompletePart via QuerySignHTTPCompletePart", func() {
			mu, ok := store.(types.Multiparter)
			So(ok, ShouldBeTrue)

			path := uuid.New().String()
			o, err := mu.CreateMultipart(path)
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

			_, part, err := mu.WriteMultipart(o, r, size, partNumber)
			if err != nil {
				t.Error(err)
			}

			req, err := signer.QuerySignHTTPCompleteMultipart(o, []*types.Part{part}, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			_, err = client.Do(req)

			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The object should be readable after complete", func() {
				ro, err := store.Stat(path)

				So(err, ShouldBeNil)
				So(ro.Mode.IsRead(), ShouldBeTrue)
				So(ro.Mode.IsPart(), ShouldBeFalse)
			})
		})
	})
}
