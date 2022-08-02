package tests

import (
	"testing"

	"github.com/beyondstorage/go-storage/v5/tests"
)

func TestStorage(t *testing.T) {
	tests.TestStorager(t, setupTest(t))
}

func TestAppend(t *testing.T) {
	tests.TestAppender(t, setupTest(t))
}

func TestDir(t *testing.T) {
	tests.TestDirer(t, setupTest(t))
}

func TestCopy(t *testing.T) {
	tests.TestCopier(t, setupTest(t))
}

func TestMove(t *testing.T) {
	tests.TestMover(t, setupTest(t))
}
