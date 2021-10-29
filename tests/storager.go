package tests

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/randbytes"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func TestStorage(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		So(store, ShouldNotBeNil)

		Convey("Test basic operations for storage service", func() {
			TestStorager(t, store)
		})

		Convey("Test Append related operations", func() {
			TestAppender(t, store)
		})

		Convey("Test Copy operation", func() {
			TestCopier(t, store)
		})

		Convey("Test CreateDir operation", func() {
			TestDirer(t, store)
		})

		Convey("Test CreateLink operation", func() {
			TestLinker(t, store)
		})

		Convey("Test Move operation", func() {
			TestMover(t, store)
		})

		Convey("Test multipart related operations which support authentication", func() {
			TestMultipartHTTPSigner(t, store)
		})

		Convey("Test Multipart related operations", func() {
			TestMultiparter(t, store)
		})

		Convey("Test basic operations which support authentication", func() {
			TestStorageHTTPSignerDelete(t, store)
			TestStorageHTTPSignerRead(t, store)
			TestStorageHTTPSignerWrite(t, store)
		})
	})
}

func TestStorager(t *testing.T, store types.Storager) {
	Convey("Given a basic Storager", t, func() {
		So(store, ShouldNotBeNil)

		f := store.Features()
		if f.Create && f.Delete && f.Metadata && f.List && f.Read && f.Stat && f.Write {
			Convey("When String called", func() {
				s := store.String()

				Convey("The string should not be empty", func() {
					So(s, ShouldNotBeEmpty)
				})
			})

			Convey("When Metadata called", func() {
				m := store.Metadata()

				Convey("The metadata should not be empty", func() {
					So(m, ShouldNotBeEmpty)
				})
			})

			workDir := store.Metadata().WorkDir

			Convey("When Read a file", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

				var buf bytes.Buffer

				n, err := store.Read(path, &buf)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The content should be match", func() {
					So(buf, ShouldNotBeNil)

					So(n, ShouldEqual, size)
					So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
				})
			})

			Convey("When Read a file with offset or size", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

				Convey("When Read with offset", func() {
					offset := rand.Int63n(size)

					var buf bytes.Buffer
					n, err := store.Read(path, &buf, ps.WithOffset(offset))

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The content should be match", func() {
						So(buf, ShouldNotBeNil)

						So(n, ShouldEqual, size-offset)
						So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content[offset:]))
					})
				})

				Convey("When Read with size", func() {
					len := rand.Int63n(size)

					var buf bytes.Buffer
					n, err := store.Read(path, &buf, ps.WithSize(len))

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The content should be match", func() {
						So(buf, ShouldNotBeNil)

						So(n, ShouldEqual, len)
						So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content[:len]))
					})
				})

				Convey("When Read with offset and size", func() {
					offset := rand.Int63n(size)
					len := rand.Int63n(size - offset)

					var buf bytes.Buffer
					n, err := store.Read(path, &buf, ps.WithOffset(offset), ps.WithSize(len))

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The content should be match", func() {
						So(buf, ShouldNotBeNil)

						So(n, ShouldEqual, len)
						So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content[offset:offset+len]))
					})
				})
			})

			Convey("When Write a file", func() {
				firstSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				r := io.LimitReader(randbytes.NewRand(), firstSize)
				path := uuid.New().String()

				_, err := store.Write(path, r, firstSize)

				defer func() {
					err := store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The first returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				secondSize := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				content, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), secondSize))

				_, err = store.Write(path, bytes.NewReader(content), secondSize)

				Convey("The second returned error also should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("Stat should get Object without error", func() {
					o, err := store.Stat(path)

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The name and size should be match", func() {
						So(o, ShouldNotBeNil)
						So(o.Path, ShouldEqual, path)

						osize, ok := o.GetContentLength()
						So(ok, ShouldBeTrue)
						So(osize, ShouldEqual, secondSize)
					})
				})

				Convey("Read should get Object data without error", func() {
					var buf bytes.Buffer
					n, err := store.Read(path, &buf)

					Convey("The error should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The content should be match", func() {
						So(buf, ShouldNotBeNil)

						So(n, ShouldEqual, secondSize)
						So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
					})
				})
			})

			Convey("When Write and Read a file with IoCallback", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
				if err != nil {
					t.Error(err)
				}

				path := uuid.New().String()

				curWrite := int64(0)
				writeFn := func(bs []byte) {
					curWrite += int64(len(bs))
				}
				_, err = store.Write(path, bytes.NewReader(content), size, ps.WithIoCallback(writeFn))
				defer func() {
					err := store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				Convey("The error returned by Write should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The write size should be match", func() {
					So(curWrite, ShouldEqual, size)
				})

				curRead := int64(0)
				readFn := func(bs []byte) {
					curRead += int64(len(bs))
				}
				var buf bytes.Buffer
				n, err := store.Read(path, &buf, ps.WithIoCallback(readFn))

				Convey("The error returned be Read should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The read size should be match", func() {
					So(curRead, ShouldEqual, n)
				})

				Convey("The content should be match", func() {
					So(buf, ShouldNotBeNil)

					So(n, ShouldEqual, size)
					So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
				})
			})

			if f.WriteEmptyObject {
				Convey("When write a file with a nil io.Reader and 0 size", func() {
					path := uuid.New().String()
					var size int64 = 0

					_, err := store.Write(path, nil, size)

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
							So(o.Path, ShouldEqual, path)

							osize, ok := o.GetContentLength()
							So(ok, ShouldBeTrue)
							So(osize, ShouldEqual, size)
						})
					})
				})

				Convey("When write a file with a nil io.Reader and valid size", func() {
					size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
					path := uuid.New().String()

					_, err := store.Write(path, nil, size)

					Convey("The error should not be nil", func() {
						So(err, ShouldNotBeNil)
					})

					Convey("Stat should get nil Object and ObjectNotFound error", func() {
						o, err := store.Stat(path)

						So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
						So(o, ShouldBeNil)
					})
				})

				Convey("When write a file with a valid io.Reader and 0 size", func() {
					var size int64 = 0
					n := rand.Int63n(4 * 1024 * 1024)
					r := io.LimitReader(randbytes.NewRand(), n)
					path := uuid.New().String()

					_, err := store.Write(path, r, size)

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
							So(o.Path, ShouldEqual, path)

							osize, ok := o.GetContentLength()
							So(ok, ShouldBeTrue)
							So(osize, ShouldEqual, size)
						})
					})
				})

				Convey("When write a file with a valid io.Reader and length greater than size", func() {
					n := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
					size := rand.Int63n(n)
					r, _ := io.ReadAll(io.LimitReader(randbytes.NewRand(), n))
					path := uuid.New().String()

					_, err := store.Write(path, bytes.NewReader(r), size)

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
							So(o.Path, ShouldEqual, path)

							osize, ok := o.GetContentLength()
							So(ok, ShouldBeTrue)
							So(osize, ShouldEqual, size)
						})
					})

					Convey("Read should get Object without error", func() {
						content, _ := io.ReadAll(io.LimitReader(bytes.NewReader(r), size))
						var buf bytes.Buffer
						n, err := store.Read(path, &buf)

						Convey("The error should be nil", func() {
							So(err, ShouldBeNil)
						})

						Convey("The content should match the size limit of the content", func() {
							So(buf, ShouldNotBeNil)

							So(n, ShouldEqual, size)
							So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
						})
					})
				})
			}

			Convey("When Stat a file", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
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

				o, err := store.Stat(path)

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The Object name and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Path, ShouldEqual, path)

					osize, ok := o.GetContentLength()
					So(ok, ShouldBeTrue)
					So(osize, ShouldEqual, size)
				})
			})

			Convey("When Delete a file", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
				if err != nil {
					t.Error(err)
				}

				path := uuid.New().String()
				_, err = store.Write(path, bytes.NewReader(content), size)
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

				Convey("Stat should get nil Object and ObjectNotFound error", func() {
					o, err := store.Stat(path)

					So(errors.Is(err, services.ErrObjectNotExist), ShouldBeTrue)
					So(o, ShouldBeNil)
				})
			})

			Convey("When List an empty dir", func() {
				it, err := store.List("", ps.WithListMode(types.ListModeDir))

				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("The iterator should not be nil", func() {
					So(it, ShouldNotBeNil)
				})

				o, err := it.Next()

				Convey("The next should be done", func() {
					So(err, ShouldBeError, types.IterateDone)
				})
				Convey("The object should be nil", func() {
					So(o, ShouldBeNil)
				})
			})

			Convey("When List a dir within files", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				r := io.LimitReader(randbytes.NewRand(), size)
				path := uuid.New().String()
				_, err := store.Write(path, r, size)
				if err != nil {
					t.Error(err)
				}
				defer func() {
					err := store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				it, err := store.List("", ps.WithListMode(types.ListModeDir))
				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("The iterator should not be nil", func() {
					So(it, ShouldNotBeNil)
				})

				o, err := it.Next()
				Convey("The name and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Path, ShouldEqual, path)

					osize, ok := o.GetContentLength()
					So(ok, ShouldBeTrue)
					So(osize, ShouldEqual, size)
				})
			})

			Convey("When List without ListMode", func() {
				size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
				r := io.LimitReader(randbytes.NewRand(), size)
				path := uuid.New().String()
				_, err := store.Write(path, r, size)
				if err != nil {
					t.Error(err)
				}
				defer func() {
					err := store.Delete(path)
					if err != nil {
						t.Error(err)
					}
				}()

				it, err := store.List("")
				Convey("The error should be nil", func() {
					So(err, ShouldBeNil)
				})
				Convey("The iterator should not be nil", func() {
					So(it, ShouldNotBeNil)
				})

				o, err := it.Next()
				Convey("The name and size should be match", func() {
					So(o, ShouldNotBeNil)
					So(o.Path, ShouldEqual, path)

					osize, ok := o.GetContentLength()
					So(ok, ShouldBeTrue)
					So(osize, ShouldEqual, size)
				})
			})

			Convey("When testing GSP-749 unify path behavior", func() {
				Convey("When using absolute path", func() {
					size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
					content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
					if err != nil {
						t.Error(err)
					}

					path := uuid.New().String()
					absPath := filepath.Join(workDir, path)
					_, err = store.Write(absPath, bytes.NewReader(content), size)
					if err != nil {
						t.Error(err)
					}
					defer func() {
						err := store.Delete(absPath)
						if err != nil {
							t.Error(err)
						}
					}()

					Convey("Stat should get Object without error", func() {
						o, err := store.Stat(absPath)

						Convey("The error should be nil", func() {
							So(err, ShouldBeNil)
							So(o, ShouldNotBeNil)
							So(o.Path, ShouldEqual, strings.ReplaceAll(absPath, "\\", "/"))
						})
					})

					Convey("Read should get Object content without error", func() {
						var buf bytes.Buffer
						n, err := store.Read(absPath, &buf)

						Convey("The error should be nil", func() {
							So(err, ShouldBeNil)
						})

						Convey("The content should be match", func() {
							So(buf, ShouldNotBeNil)

							So(n, ShouldEqual, size)
							So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
						})
					})
				})

				Convey("When using backslash in path", func() {
					size := rand.Int63n(4 * 1024 * 1024) // Max file size is 4MB
					content, err := io.ReadAll(io.LimitReader(randbytes.NewRand(), size))
					if err != nil {
						t.Error(err)
					}

					path := uuid.New().String() + "\\" + uuid.New().String()
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

					Convey("Stat should get Object without error", func() {
						o, err := store.Stat(path)

						Convey("The error should be nil", func() {
							So(err, ShouldBeNil)
							So(o, ShouldNotBeNil)
							So(o.Path, ShouldEqual, strings.ReplaceAll(path, "\\", "/"))
						})
					})

					Convey("Read should get Object content without error", func() {
						var buf bytes.Buffer
						n, err := store.Read(path, &buf)

						Convey("The error should be nil", func() {
							So(err, ShouldBeNil)
						})

						Convey("The content should be match", func() {
							So(buf, ShouldNotBeNil)

							So(n, ShouldEqual, size)
							So(sha256.Sum256(buf.Bytes()), ShouldResemble, sha256.Sum256(content))
						})
					})
				})
			})
		}
	})
}
