package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"
	ipfs "go.beyondstorage.io/services/ipfs"

	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for IPFS")

	store, err := ipfs.NewStorager(
		pairs.WithEndpoint(os.Getenv("STORAGE_IPFS_ENDPOINT")),
		ipfs.WithGateway(os.Getenv("STORAGE_IPFS_GATEWAY")),
		pairs.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
