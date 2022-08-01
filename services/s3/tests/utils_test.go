package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	s3 "github.com/beyondstorage/go-storage/services/s3/v3"
	ps "github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for s3")

	store, err := s3.NewStorager(
		ps.WithCredential(os.Getenv("STORAGE_S3_CREDENTIAL")),
		ps.WithName(os.Getenv("STORAGE_S3_NAME")),
		ps.WithLocation(os.Getenv("STORAGE_S3_LOCATION")),
		ps.WithEndpoint(os.Getenv("STORAGE_S3_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
		ps.WithEnableVirtualDir(),
		ps.WithEnableVirtualLink(),
		s3.WithForcePathStyle(),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
