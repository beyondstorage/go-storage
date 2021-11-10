package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_FromString(t *testing.T) {
	f := &Factory{}

	err := f.FromString("hmac:ak:sk@http:xxx:xxx/bucket/dir/?disable_uri_cleaning")
	assert.NoError(t, err)
	assert.Equal(t, &Factory{
		Credential:         "hmac:ak:sk",
		DisableURICleaning: true,
		Endpoint:           "http:xxx:xxx",
		Name:               "bucket",
		WorkDir:            "/dir/",
	}, f)
}
