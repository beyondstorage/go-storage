package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	oss "github.com/beyondstorage/go-storage/services/oss/v3"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for oss")

	store, err := oss.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_OSS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_OSS_NAME")),
		ps.WithEndpoint(os.Getenv("STORAGE_OSS_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEnableVirtualDir(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
