package tests

import (
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_GCS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_GCS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestDirer(t *testing.T) {
	if os.Getenv("STORAGE_GCS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_GCS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}
