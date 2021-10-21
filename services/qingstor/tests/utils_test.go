package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	"go.beyondstorage.io/services/qingstor/v4"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for qingstor")

	store, err := qingstor.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_QINGSTOR_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_QINGSTOR_ENDPOINT")),
		ps.WithName(os.Getenv("STORAGE_QINGSTOR_NAME")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		qingstor.WithStorageFeatures(qingstor.StorageFeatures{
			VirtualDir:  true,
			VirtualLink: true,
		}),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
