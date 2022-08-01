package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	gcs "github.com/beyondstorage/go-storage/services/gcs/v3"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for gcs")

	store, err := gcs.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_GCS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_GCS_NAME")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		gcs.WithProjectID(os.Getenv("STORAGE_GCS_PROJECT_ID")),
		ps.WithEnableVirtualDir(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
