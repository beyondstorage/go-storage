package qingstor

import (
	"fmt"

	"github.com/yunify/qingstor-sdk-go/v3/config"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
)

// Config is the qingstor service config.
type Config struct {
	AccessKeyID     string
	SecretAccessKey string

	Host     string
	Port     int
	Protocol string

	Zone       string
	BucketName string
}

// New implements Configurer interface.
func (c *Config) New() (storage.Storager, error) {
	errorMessage := "create new qingstor client failed: %w"

	cfg, err := config.New(c.AccessKeyID, c.SecretAccessKey)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	cfg.Host = c.Host
	cfg.Port = c.Port
	cfg.Protocol = c.Protocol

	client := &Client{
		config:   c,
		segments: make(map[string]*segment.Segment),
	}
	client.service, err = service.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	if c.BucketName == "" {
		return client, nil
	}

	err = client.setupBucket(c.BucketName, c.Zone)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}
	return client, nil
}
