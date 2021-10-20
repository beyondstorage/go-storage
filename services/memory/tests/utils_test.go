package tests

import (
	"testing"

	"go.beyondstorage.io/v5/types"

	"go.beyondstorage.io/services/memory"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for memory")

	store, err := memory.NewStorager()
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
