package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	bos "github.com/beyondstorage/go-storage/services/bos/v2"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for bos")

	store, err := bos.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_BOS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_BOS_NAME")),
		ps.WithEndpoint(os.Getenv("STORAGE_BOS_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEnableVirtualDir(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	return store
}
