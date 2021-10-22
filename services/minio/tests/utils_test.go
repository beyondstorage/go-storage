package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	minio "go.beyondstorage.io/services/minio"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for minio")

	srv, err := minio.NewServicer(
		ps.WithCredential(os.Getenv("STORAGE_MINIO_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_MINIO_ENDPOINT")),
	)
	if err != nil {
		t.Errorf("new servicer: %v", err)
	}

	bucketName := os.Getenv("STORAGE_MINIO_NAME")

	_, err = srv.Create(bucketName)
	if err != nil {
		t.Errorf("create storager: %v", err)
	}

	store, err := minio.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_MINIO_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_MINIO_ENDPOINT")),
		ps.WithName(bucketName),
		ps.WithWorkDir("/"+uuid.New().String()),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}

	t.Cleanup(func() {
		err = store.Delete("")
		if err != nil {
			t.Errorf("cleanup: %v", err)
		}

		err = srv.Delete(bucketName)
		if err != nil {
			t.Errorf("cleanup: %v", err)
		}
	})
	return store
}
