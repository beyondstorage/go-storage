package main

import (
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

const (
	generatedPath = "generated.go"
)

var (
	headerT = template.Must(
		template.New("header").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("header.tmpl"))))
	contextT = template.Must(
		template.New("context").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("context.tmpl"))))
	metaT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("meta.tmpl"))))
	pairsT = template.Must(
		template.New("pairs").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("pairs.tmpl"))))
)

//go:generate go-bindata -nometadata -ignore "\\.go$" -prefix tmpl ./tmpl
func main() {
	f, err := os.Create(generatedPath)
	if err != nil {
		log.Fatal(err)
	}

	meta := parseMeta()

	// Generate header
	err = headerT.Execute(f, meta)
	if err != nil {
		log.Fatal(err)
	}

	// Generate meta
	err = metaT.Execute(f, meta)
	if err != nil {
		log.Fatal(err)
	}

	// Generate pairs
	err = pairsT.Execute(f, meta)
	if err != nil {
		log.Fatal(err)
	}

	// Generate context
	err = contextT.Execute(f, meta)
	if err != nil {
		log.Fatal(err)
	}
}
