package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	storj "go.beyondstorage.io/services/storj"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for STORJ")

	store, err := storj.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_STORJ_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_STORJ_NAME")),
		ps.WithLocation(os.Getenv("STORAGE_STORJ_LOCATION")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
