package tests

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/coreutils"
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
	})
}
