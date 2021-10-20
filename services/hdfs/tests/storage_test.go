package tests

import (
	"os"
	"testing"

	"go.beyondstorage.io/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_HDFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_HDFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestDirer(t *testing.T) {
	if os.Getenv("STORAGE_HDFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_HDFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}

func TestMover(t *testing.T) {
	if os.Getenv("STORAGE_HDFS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_HDFS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestMover(t, setupTest(t))
}
