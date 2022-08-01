package tests

import (
	"os"
	"testing"

	"github.com/beyondstorage/go-storage/services/gdrive"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"

	"github.com/google/uuid"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for gdrive")

	store, err := gdrive.NewStorager(
		ps.WithName(os.Getenv("STORAGE_GDRIVE_NAME")),
		ps.WithCredential(os.Getenv("STORAGE_GDRIVE_CREDENTIAL")),
		ps.WithWorkDir("/"+uuid.New().String()),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	t.Cleanup(func() {
		err = store.Delete("")
		if err != nil {
			t.Errorf("cleanup: %v", err)
		}
	})
	return store
}
