package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

var (
	metadataT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("metadata.tmpl"))))
)

//go:generate go-bindata -nometadata -ignore "\\.go$" .
func main() {
	fi, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalf("read dir failed: %v", err)
	}

	for _, v := range fi {
		if !strings.HasSuffix(v.Name(), "json") {
			continue
		}

		metadataName := strings.TrimSuffix(v.Name(), ".json")
		metadataPath := metadataName + ".json"
		content, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			log.Fatalf("read file failed: %v", err)
		}

		var metadata map[string]struct {
			Name      string
			Type      string
			ZeroValue string `json:",omitempty"`
		}
		err = json.Unmarshal(content, &metadata)
		if err != nil {
			log.Fatalf("json unmarshal failed: %v", err)
		}

		// Format input meta.json
		data, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(metadataPath, data, 0664)
		if err != nil {
			log.Fatal(err)
		}

		metadataFile, err := os.Create(metadataName + ".go")
		if err != nil {
			log.Fatal(err)
		}

		err = metadataT.Execute(metadataFile, struct {
			Name string
			Data map[string]struct {
				Name      string
				Type      string
				ZeroValue string `json:",omitempty"`
			}
		}{
			metadataName,
			metadata,
		})
		if err != nil {
			log.Fatal(err)
		}

		metadataFile.Close()
	}
}
