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
	infoT      = newTmpl("info")
	pairT      = newTmpl("pair")
	serviceT   = newTmpl("service")
	openT      = newTmpl("open")
	operationT = newTmpl("operation")
	functionT  = newTmpl("function")
)

func generate(data *Data) {
	// Metas generate
	generateT(infoT, "../types/info/generated.go", data)

	// Pair generate
	generateT(pairT, "../types/pairs/generated.go", data)

	// Operation generate
	generateT(operationT, "../generated.go", data)

	// Service generate
	for _, v := range data.Services {
		fp := fmt.Sprintf("../services/%s/generated.go", v.Name)
		generateT(serviceT, fp, v)

		for _, ns := range v.Namespaces {
			sp := fmt.Sprintf("../services/%s/%s.go", v.Name, ns.Name)
			for _, fn := range ns.Funcs {
				if fn.implemented {
					continue
				}
				appendT(functionT, sp, struct {
					Namespace string
					Func      *Function
				}{
					ns.Name,
					fn,
				})
			}
		}
	}

	// Coreutils generate
	generateT(openT, "../coreutils/generated.go", data)
}

func generateT(tmpl *template.Template, filePath string, data interface{}) {
	errorMsg := fmt.Sprintf("generate template %s to %s", tmpl.Name(), filePath) + ": %v"

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
}

func appendT(tmpl *template.Template, filePath string, data interface{}) {
	errorMsg := fmt.Sprintf("append template %s to %s", tmpl.Name(), filePath) + ": %v"

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
}

func newTmpl(name string) *template.Template {
	return template.Must(
		template.New(name).
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset(name + ".tmpl"))))
}
