package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"
	onedrive "go.beyondstorage.io/services/onedrive"

	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for onedrive")

	store, err := onedrive.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_ONEDRIVE_CREDENTIAL")),
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
