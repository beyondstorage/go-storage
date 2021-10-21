package tests

import (
	"os"
	"testing"

	"go.beyondstorage.io/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_AZBLOB_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_AZBLOB_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestAppend(t *testing.T) {
	if os.Getenv("STORAGE_AZBLOB_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_AZBLOB_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestAppender(t, setupTest(t))
}

func TestDir(t *testing.T) {
	if os.Getenv("STORAGE_AZBLOB_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_AZBLOB_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}
