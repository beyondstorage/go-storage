package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	azfile "go.beyondstorage.io/services/azfile"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for azfile")

	store, err := azfile.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_AZFILE_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_AZFILE_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithName(os.Getenv("STORAGE_AZFILE_NAME")),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
