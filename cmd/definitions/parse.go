//go:build tools
// +build tools

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"

	"github.com/beyondstorage/go-storage/v4/cmd/definitions/specs"
)

func parse() (data *Data) {
	injectPairs()

	data = NewData()
	return data
}

func injectPairs() {
	ps := []specs.Pair{
		{
			Name: "context",
			Type: "context.Context",
		},
		{
			Name: "http_client_options",
			Type: "*httpclient.Options",
		},
	}

	for _, v := range ps {
		specs.ParsedPairs = append(specs.ParsedPairs, v)
	}
}

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
