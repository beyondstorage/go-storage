package tests

import (
	"testing"

	fs "go.beyondstorage.io/services/fs/v4"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/types"
)

func setupTest(t *testing.T) types.Storager {
	tmpDir := t.TempDir()
	t.Logf("Setup test at %s", tmpDir)

	store, err := fs.NewStorager(ps.WithWorkDir(tmpDir))
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
