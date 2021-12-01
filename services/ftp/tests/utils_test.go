package tests

import (
	"os"
	"testing"

	"github.com/google/uuid"

	_ "go.beyondstorage.io/services/ftp"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"go.beyondstorage.io/v5/types"
)

func initTest(t *testing.T) (store types.Storager) {
	t.Log("Setup test for ftp")

	store, err := services.NewStorager("ftp",
		ps.WithCredential(os.Getenv("STORAGE_FTP_CREDENTIAL")),
		ps.WithEndpoint(os.Getenv("STORAGE_FTP_ENDPOINT")),
		ps.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("create storager: %v", err)
	}

	return
}
