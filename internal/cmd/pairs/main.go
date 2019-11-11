package main

import (
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"

	"github.com/Xuanwo/storage/types"
)

var (
	pairT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("pair.tmpl"))))
	metadataT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("metadata.tmpl"))))
)

//go:generate go-bindata -ignore ".*.go" .
func main() {
	pairFile, err := os.Create("pairs.go")
	if err != nil {
		log.Fatal(err)
	}
	defer pairFile.Close()

	err = pairT.Execute(pairFile, struct {
		Data map[string]string
	}{
		types.AvailablePairs,
	})
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
		types.AvailablePairs,
	})
	if err != nil {
		log.Fatal(err)
	}
}
