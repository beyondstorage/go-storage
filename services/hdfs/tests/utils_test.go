package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"

	hdfs "go.beyondstorage.io/services/hdfs"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for HDFS")

	store, err := hdfs.NewStorager(
		pairs.WithEndpoint(os.Getenv("STORAGE_HDFS_ENDPOINT")),
		pairs.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	return store
}
