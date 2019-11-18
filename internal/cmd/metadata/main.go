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
	metadataT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("metadata.tmpl"))))
)

//go:generate go-bindata -nometadata -ignore ".*.go" .
func main() {
	metadataPath := "metadata.json"
	content, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	var metadata map[string]string
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

	metadataFile, err := os.Create("metadata.go")
	if err != nil {
		log.Fatal(err)
	}
	defer metadataFile.Close()

	err = metadataT.Execute(metadataFile, struct {
		Data map[string]string
	}{
		metadata,
	})
	if err != nil {
		log.Fatal(err)
	}
}
