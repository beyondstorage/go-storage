package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorage(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorager(t, setupTest(t))
}

func TestMultiparter(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestMultiparter(t, setupTest(t))
}

func TestDirer(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestDirer(t, setupTest(t))
}

func TestLinker(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestLinker(t, setupTest(t))
}

func TestHTTPSigner(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	tests.TestStorageHTTPSignerWrite(t, setupTest(t))
	tests.TestStorageHTTPSignerRead(t, setupTest(t))
	// presign operations don't support DeleteObject & CreateMultipartUpload
	// DeleteObject is used in TestStorageHTTPSignerDelete
	// CreateMultipartUpload is used in TestMultipartHTTPSigner
	// tests.TestStorageHTTPSignerDelete(t, setupTest(t)) is not supported
	// tests.TestMultipartHTTPSigner(t, setupTest(t)) is not supported
}

// https://github.com/beyondstorage/go-storage/issues/741
func TestIssue741(t *testing.T) {
	if os.Getenv("STORAGE_S3_INTEGRATION_TEST") != "on" {
		t.Skipf("STORAGE_S3_INTEGRATION_TEST is not 'on', skipped")
	}
	store := setupTest(t)

	content := []byte("Hello, World!")
	r := bytes.NewReader(content)

	_, err := store.Write("IMG@@@¥&_0960.jpg", r, int64(len(content)))
	if err != nil {
		t.Errorf("write: %v", err)
		return
	}
	return
}
