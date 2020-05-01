package tests

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const TestPrefix = "STORAGE_TEST_SERVICE_"

type config struct {
	Type    string
	Options map[string]interface{}
}

func loadConfig() []config {
	srv := make([]config, 0)

	// Read from local config
	content, err := ioutil.ReadFile("storager.yaml")
	if err == nil {
		err = yaml.Unmarshal(content, &srv)
		if err != nil {
			log.Fatal("load storager.yaml failed")
		}
	}

	// Read from env
	env := os.Environ()
	for _, v := range env {
		values := strings.SplitN(v, "=", 2)

		if !strings.HasPrefix(values[0], TestPrefix) {
			continue
		}

		cfg := config{
			// STORAGE_TEST_SERVICE_FS => FS => fs
			Type: strings.ToLower(strings.TrimPrefix(values[0], TestPrefix)),
		}
		// Config via env will be yaml content after base64.
		content, err := base64.StdEncoding.DecodeString(values[0])
		if err != nil {
			log.Fatal("base64 decode config failed")
		}
		err = yaml.Unmarshal(content, cfg.Options)
		if err != nil {
			log.Fatal("yaml unmarshal config failed")
		}

		srv = append(srv, cfg)
	}
	return srv
}
