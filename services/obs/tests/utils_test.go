package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	obs "github.com/beyondstorage/go-storage/services/obs/v2"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for obs")

	store, err := obs.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_OBS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_OBS_NAME")),
		ps.WithEndpoint(os.Getenv("STORAGE_OBS_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEnableVirtualDir(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	return store
}
