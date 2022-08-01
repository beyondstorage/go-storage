package tests

import (
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_OSS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OSS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestAppend(t *testing.T) {
	if os.Getenv("STORAGE_OSS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OSS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestAppender(t, setupTest(t))
}

func TestMultiparter(t *testing.T) {
	if os.Getenv("STORAGE_OSS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OSS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestMultiparter(t, setupTest(t))
}

func TestDirer(t *testing.T) {
	if os.Getenv("STORAGE_OSS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OSS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}

func TestLinker(t *testing.T) {
	if os.Getenv("STORAGE_OSS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OSS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestLinker(t, setupTest(t))
}
