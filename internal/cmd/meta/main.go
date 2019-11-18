package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

var (
	metaT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("meta.tmpl"))))
)

type metadata struct {
	Name    string                     `json:"name"`
	Service map[string]map[string]bool `json:"service,omitempty"`
	Storage map[string]map[string]bool `json:"storage"`

	TypeMap map[string]string `json:"-"`
}

//go:generate go-bindata -nometadata -ignore ".*.go" .
func main() {
	_, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalf("read dir failed: %v", err)
	}

	metaPath := "meta.json"
	if _, err := os.Stat(metaPath); err != nil {
		log.Fatalf("stat meta failed: %v", err)
	}

	content, err := ioutil.ReadFile(metaPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	var meta metadata

	err = json.Unmarshal(content, &meta)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	pairsPath := "../../types/pairs/pairs.json"
	content, err = ioutil.ReadFile(pairsPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}
	err = json.Unmarshal(content, &meta.TypeMap)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	// Format input meta.json
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(metaPath, data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	filePath := "meta.go"
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = metaT.Execute(f, meta)
	if err != nil {
		log.Fatal(err)
	}
}
