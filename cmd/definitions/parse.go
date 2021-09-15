//go:build tools
// +build tools

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

const (
	featurePath     = "definitions/features.toml"
	fieldPath       = "definitions/fields.toml"
	infoObjectMeta  = "definitions/info_object_meta.toml"
	infoStorageMeta = "definitions/info_storage_meta.toml"
	operationPath   = "definitions/operations.toml"
	pairPath        = "definitions/pairs.toml"
)

func parseFunc(name string) map[string]*templateutils.Method {
	data := make(map[string]*templateutils.Method)
	filename := fmt.Sprintf("%s.go", name)

	content, err := ioutil.ReadFile(filename)
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
