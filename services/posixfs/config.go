package posixfs

import (
	"github.com/Xuanwo/storage"
)

// Config is the posixfs config.
type Config struct {
	Path string
}

// New implements Configurer interface.
func (c *Config) New() (storage.Storager, error) {
	client := &Client{}
	return client, nil
}
