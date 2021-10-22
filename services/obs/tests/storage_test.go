package tests

import (
	"os"
	"testing"

	"go.beyondstorage.io/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_OBS_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_OBS_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}
