package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	kodo "go.beyondstorage.io/services/kodo/v3"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for kodo")

	store, err := kodo.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_KODO_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_KODO_NAME")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEndpoint(os.Getenv("STORAGE_KODO_ENDPOINT")),
		kodo.WithStorageFeatures(kodo.StorageFeatures{
			VirtualDir: true,
		}),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
