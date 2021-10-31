package tests

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

// Deprecated: Moved to TestStorager.
func TestAppender(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

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
	})
}

// Deprecated: Moved to TestStorager.
func TestCopier(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Copy a file", func() {
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
			err = store.Copy(src, dst)

			defer func() {
				err = store.Delete(dst)
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
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

		Convey("When Copy to an existing file", func() {
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

			err = store.Copy(src, dst)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
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

// Deprecated: Moved to TestStorager.
func TestCopierWithDir(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Copy to an existing dir", func() {
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

			err = store.Copy(src, dst)
			Convey("The error should be ErrObjectModeInvalid", func() {
				So(errors.Is(err, services.ErrObjectModeInvalid), ShouldBeTrue)
			})
		})
	})
}

// Deprecated: Moved to TestStorager.
func TestCopierWithVirtualDir(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Copy to an existing dir", func() {
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

			err = store.Copy(src, dst)

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

// Deprecated: Moved to TestStorager.
func TestDirer(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When CreateDir", func() {
			path := uuid.New().String()
			_, err := store.CreateDir(path)

			defer func() {
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The first returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			o, err := store.CreateDir(path)
			Convey("The second returned error also should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Path should equal to the input path", func() {
				So(o.Path, ShouldEqual, path)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Create with ModeDir", func() {
			path := uuid.New().String()
			o := store.Create(path, pairs.WithObjectMode(types.ModeDir))

			defer func() {
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			Convey("The Object Path should equal to the input path", func() {
				So(o.Path, ShouldEqual, path)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Stat with ModeDir", func() {
			path := uuid.New().String()
			_, err := store.CreateDir(path)
			if err != nil {
				t.Error(err)
			}

			defer func() {
				err := store.Delete(path, pairs.WithObjectMode(types.ModeDir))
				if err != nil {
					t.Error(err)
				}
			}()

			o, err := store.Stat(path, pairs.WithObjectMode(types.ModeDir))

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Object Path should equal to the input path", func() {
				So(o.Path, ShouldEqual, path)
			})

			Convey("The Object Mode should be dir", func() {
				// Dir object's mode must be Dir.
				So(o.Mode.IsDir(), ShouldBeTrue)
			})
		})

		Convey("When Delete with ModeDir", func() {
			path := uuid.New().String()
			_, err := store.CreateDir(path)
			if err != nil {
				t.Error(err)
			}

			err = store.Delete(path, pairs.WithObjectMode(types.ModeDir))
			Convey("The first returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			err = store.Delete(path, pairs.WithObjectMode(types.ModeDir))
			Convey("The second returned error also should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

// Deprecated: Moved to TestStorager.
func TestLinker(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

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
	})
}

// Deprecated: Moved to TestStorager.
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

// Deprecated: Moved to TestStorager.
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

// Deprecated: Moved to TestStorager.
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

// Deprecated: Moved to TestStorager.
func TestMultipartHTTPSigner(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When CreateMultipart via QuerySignHTTPCreateMultipart", func() {
			path := uuid.New().String()
			req, err := store.QuerySignHTTPCreateMultipart(path, time.Duration(time.Hour))

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

			size := rand.Int63n(4 * 1024 * 1024)
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			req, err := store.QuerySignHTTPWriteMultipart(o, size, 0, time.Duration(time.Hour))

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

			req, err := store.QuerySignHTTPListMultipart(o, time.Duration(time.Hour))

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

			req, err := store.QuerySignHTTPCompleteMultipart(o, []*types.Part{part}, time.Duration(time.Hour))

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

// Deprecated: Moved to TestStorager.
func TestMultiparter(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

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
	})
}

// Deprecated: Moved to TestStorager.
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

			req, err := store.QuerySignHTTPRead(path, time.Hour)

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

// Deprecated: Moved to TestStorager.
func TestStorageHTTPSignerWrite(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {

		Convey("When Write via QuerySignHTTPWrite", func() {
			size := rand.Int63n(4 * 1024 * 1024)
			content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
			if err != nil {
				t.Error(err)
			}

			path := uuid.New().String()
			req, err := store.QuerySignHTTPWrite(path, size, time.Hour)

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

// Deprecated: Moved to TestStorager.
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

			req, err := store.QuerySignHTTPDelete(path, time.Hour)

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

			req, err := store.QuerySignHTTPDelete(path, time.Hour, pairs.WithMultipartID(o.MustGetMultipartID()))

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
