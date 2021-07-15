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

var typeMap = map[string]string{
	"context":             "context.Context",
	"http_client_options": "*httpclient.Options",

	"any":               "interface{}",
	"byte_array":        "[]byte",
	"string_array":      "[]string",
	"string_string_map": "map[string]string",
	"time":              "time.Time",
	"BlockIterator":     "*BlockIterator",
	"IoCallback":        "func([]byte)",
	"Object":            "*Object",
	"ObjectIterator":    "*ObjectIterator",
	"Pairs":             "...Pair",
	"Part":              "*Part",
	"Parts":             "[]*Part",
	"PartIterator":      "*PartIterator",
	"Reader":            "io.Reader",
	"StoragerIterator":  "*StoragerIterator",
	"StorageMeta":       "*StorageMeta",
	"Writer":            "io.Writer",
}

// TODO: We can remove this convert after all service migrated.
func parseType(v string) string {
	s, ok := typeMap[v]
	if !ok {
		return v
	}
	log.Warnf("type %s is not supported anymore, please updated to %s.", v, s)
	return s
}
