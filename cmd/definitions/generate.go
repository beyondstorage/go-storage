package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

//go:generate go-bindata -nometadata -ignore "\\.go$" -prefix "../../" ./tmpl ../../definitions

var (
	infoT      = newTmpl("cmd/definitions/tmpl/info")
	pairT      = newTmpl("cmd/definitions/tmpl/pair")
	serviceT   = newTmpl("cmd/definitions/tmpl/service")
	operationT = newTmpl("cmd/definitions/tmpl/operation")
	functionT  = newTmpl("cmd/definitions/tmpl/function")
)

func generateGlobal(data *Data) {
	// Metas generate
	generateT(infoT, "types/info.generated.go", data)

	// Pair generate
	generateT(pairT, "pairs/generated.go", data)

	// Operation generate
	generateT(operationT, "types/operation.generated.go", data)
}

func generateService(data *Data) {
	generateT(serviceT, "generated.go", data.Service)
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
