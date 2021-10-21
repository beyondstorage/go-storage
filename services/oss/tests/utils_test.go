package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	oss "go.beyondstorage.io/services/oss/v3"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for oss")

	store, err := oss.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_OSS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_OSS_NAME")),
		ps.WithEndpoint(os.Getenv("STORAGE_OSS_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		oss.WithStorageFeatures(oss.StorageFeatures{
			VirtualDir: true,
		}),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
