package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

const (
	generatedPath = "generated.go"
	servicesPath  = "../services"
	metaPath      = "meta.json"
)

var (
	openT = template.Must(
		template.New("open").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("open.tmpl"))))
)

type metadata struct {
	Name    string                     `json:"name"`
	Service map[string]map[string]bool `json:"service,omitempty"`
	Storage map[string]map[string]bool `json:"storage"`
}

//go:generate go-bindata -nometadata -ignore "\\.go$" -prefix tmpl ./tmpl
func main() {
	f, err := os.Create(generatedPath)
	if err != nil {
		log.Fatal(err)
	}

	services, err := ioutil.ReadDir(servicesPath)
	if err != nil {
		log.Fatal(err)
	}

	metas := make([]metadata, 0)
	for _, v := range services {
		if !v.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(servicesPath, v.Name(), metaPath))
		if err != nil {
			log.Fatal(err)
		}

		var meta metadata

		err = json.Unmarshal(content, &meta)
		if err != nil {
			log.Fatalf("json unmarshal failed: %v", err)
		}

		metas = append(metas, meta)
	}

	err = openT.Execute(f, metas)
	if err != nil {
		log.Fatal(err)
	}
}
