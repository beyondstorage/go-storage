package tests

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func TestStorageHTTPSignerRead(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Read via QuerySignHTTPRead", func() {
			size := rand.Int63n(4 * 1024 * 1024)
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			_, err = store.Write(path, bytes.NewReader(content), size)
			if err != nil {
				t.Error(err)
			}
			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			req, err := store.QuerySignHTTPRead(path, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			resp, err := client.Do(req)
			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})

			defer resp.Body.Close()

			buf, err := io.ReadAll(resp.Body)
			Convey("The content should be match", func() {
				So(err, ShouldBeNil)
				So(buf, ShouldNotBeNil)

				So(resp.ContentLength, ShouldEqual, size)
				So(sha256.Sum256(buf), ShouldResemble, sha256.Sum256(content))
			})
		})
	})
}

func TestStorageHTTPSignerWrite(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Write via QuerySignHTTPWrite", func() {
			size := rand.Int63n(4 * 1024 * 1024)
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			req, err := store.QuerySignHTTPWrite(path, size, time.Duration(time.Hour))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			req.Body = io.NopCloser(bytes.NewReader(content))

			client := http.Client{}
			_, err = client.Do(req)
			Convey("The request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			defer func() {
				err := store.Delete(path)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("Read should get object data without error", func() {
				var buf bytes.Buffer
				n, err := store.Read(path, &buf)

				Convey("The content should be match", func() {
					So(err, ShouldBeNil)
					So(buf, ShouldNotBeNil)

					So(n, ShouldEqual, size)
					So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
				})
			})
		})
	})
}

func TestStorageHTTPSignerDelete(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Delete via QuerySignHTTPDelete", func() {
			size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
			r := io.LimitReader(randbytes.NewRand(), size)

			path := uuid.New().String()
			_, err := store.Write(path, r, size)
			if err != nil {
				t.Error(err)
			}

			req, err := store.QuerySignHTTPDelete(path, time.Duration(time.Hour))

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

			Convey("Stat should get nil Object and ObjectNotFound error", func() {
				o, err := store.Stat(path)

				So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
				So(o, ShouldBeNil)
			})
		})

		Convey("When Delete with multipart id via QuerySignHTTPDelete", func() {
			path := uuid.New().String()
			o, err := store.CreateMultipart(path)
			if err != nil {
				t.Error(err)
			}

			req, err := store.QuerySignHTTPDelete(path, time.Duration(time.Hour), pairs.WithMultipartID(o.MustGetMultipartID()))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)

				So(req, ShouldNotBeNil)
				So(req.URL, ShouldNotBeNil)
			})

			client := http.Client{}
			_, err = client.Do(req)

			Convey("The first request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			_, err = client.Do(req)

			Convey("The second request returned error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
