package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"
	azblob "go.beyondstorage.io/services/azblob/v3"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for azblob")

	store, err := azblob.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_AZBLOB_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_AZBLOB_NAME")),
		ps.WithEndpoint(os.Getenv("STORAGE_AZBLOB_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
