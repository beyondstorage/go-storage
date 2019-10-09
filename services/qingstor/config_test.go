package qingstor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_New(t *testing.T) {
	cfg := &Config{
		AccessKeyID:     "test_access_key",
		SecretAccessKey: "secret_access_key",
		Host:            "example.com",
		Port:            80,
		Protocol:        "http",
		Zone:            "test",
		BucketName:      "test",
	}

	s, err := cfg.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.(*Client))
	c := s.(*Client)
	assert.Equal(t, "test_access_key", c.config.AccessKeyID)
	assert.Equal(t, "test", *c.bucket.Properties.Zone)
}
