//go:build tools
// +build tools

package specs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	assert.NotEmpty(t, ParsedFeatures)
	assert.NotEmpty(t, ParsedPairs)
	assert.NotEmpty(t, ParsedOperations)
	assert.NotEmpty(t, ParsedInfos)
}

func TestParseService(t *testing.T) {
	srv, err := ParseService("testdata/service.toml")
	if err != nil {
		t.Error("parse service", err)
		return
	}

	assert.Equal(t, "tests", srv.Name)
	assert.Equal(t, "storage", srv.Namespaces[1].Name)
	assert.Equal(t, []string{"virtual_dir"}, srv.Namespaces[1].Features)
}
