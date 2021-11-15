package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	cos "go.beyondstorage.io/services/cos/v3"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for oss")

	store, err := cos.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_COS_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_COS_NAME")),
		ps.WithLocation(os.Getenv("STORAGE_COS_LOCATION")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEnableVirtualDir(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
