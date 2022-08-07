package tests

import (
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorger(t *testing.T) {
	if os.Getenv("STORAGE_FTP_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_FTP_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, initTest(t))
}
