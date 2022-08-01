package tests

import (
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_MINIO_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_MINIO_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestCopier(t *testing.T) {
	if os.Getenv("STORAGE_MINIO_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_MINIO_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestCopier(t, setupTest(t))
}
