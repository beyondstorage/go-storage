//go:build tools
// +build tools

package main

import (
	"fmt"
	"os"

	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

const (
	featurePath     = "features.toml"
	fieldPath       = "fields.toml"
	infoObjectMeta  = "info_object_meta.toml"
	infoStorageMeta = "info_storage_meta.toml"
	operationPath   = "operations.toml"
	pairPath        = "pairs.toml"
)

func parseFunc(name string) map[string]*templateutils.Method {
	data := make(map[string]*templateutils.Method)
	filename := fmt.Sprintf("%s.go", name)

	content, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return data
	}
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	source, err := templateutils.ParseContent(filename, content)
	if err != nil {
		log.Fatalf("parse content: %v", err)
	}

	for _, v := range source.Methods {
		v := v
		data[v.Name] = v
	}
	return data
}
