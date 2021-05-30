// +build tools

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Xuanwo/templateutils"
	specs "github.com/beyondstorage/specs/go"
)

func parse() (data *Data) {
	injectPairs()

	data = FormatData(specs.ParsedPairs, specs.ParsedInfos, specs.ParsedOperations)
	return data
}

func injectPairs() {
	ps := []specs.Pair{
		{
			Name: "context",
			Type: "context",
		},
		{
			Name: "http_client_options",
			Type: "http_client_options",
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
	// Golang basic types
	"error":  "error",
	"string": "string",
	"int":    "int",
	"int64":  "int64",
	"bool":   "bool",

	// Golang self-defined types.
	"context":             "context.Context",
	"http_client_options": "*httpclient.Options",

	// Compose types
	"any":                 "interface{}",
	"byte_array":          "[]byte",
	"string_array":        "[]string",
	"string_string_map":   "map[string]string",
	"time":                "time.Time",
	"BlockIterator":       "*BlockIterator",
	"DefaultServicePairs": "DefaultServicePairs",
	"DefaultStoragePairs": "DefaultStoragePairs",
	"ServiceFeatures":     "ServiceFeatures",
	"StorageFeatures":     "StorageFeatures",
	"ListMode":            "ListMode",
	"Interceptor":         "Interceptor",
	"IoCallback":          "func([]byte)",
	"Object":              "*Object",
	"ObjectMode":          "ObjectMode",
	"ObjectIterator":      "*ObjectIterator",
	"Pairs":               "...Pair",
	"Part":                "*Part",
	"Parts":               "[]*Part",
	"PartIterator":        "*PartIterator",
	"Reader":              "io.Reader",
	"Storager":            "Storager",
	"StoragerIterator":    "*StoragerIterator",
	"StorageMeta":         "*StorageMeta",
	"Writer":              "io.Writer",
}

func parseType(v string) string {
	s, ok := typeMap[v]
	if ok {
		return s
	}
	log.Fatalf("type %s is not supported.", v)
	return ""
}
