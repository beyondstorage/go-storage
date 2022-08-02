package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	azfile "github.com/beyondstorage/go-storage/services/azfile"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
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
