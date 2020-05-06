package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

//go:generate go-bindata -nometadata -ignore "\\.go$" -prefix tmpl ./tmpl

var (
	metadataT = template.Must(
		template.New("metadata").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("metadata.tmpl"))))
	pairT = template.Must(
		template.New("pair").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("pair.tmpl"))))
	serviceT = template.Must(
		template.New("service").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("service.tmpl"))))
	openT = template.Must(
		template.New("open").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("open.tmpl"))))
)

func generate(data *Data) {
	var err error

	// Metadata generate
	file, err := os.Create("../types/metadata/generated.go")
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	err = metadataT.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}

	// Pair generate
	file, err = os.Create("../types/pairs/generated.go")
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	err = pairT.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range data.Service {
		file, err = os.Create(fmt.Sprintf("../services/%s/generated.go", v.Name))
		if err != nil {
			log.Fatalf("generate: %v", err)
		}
		err = serviceT.Execute(file, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err = os.Create("../coreutils/generated.go")
	if err != nil {
		log.Fatalf("generate: %v", err)
	}
	err = openT.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}
}
