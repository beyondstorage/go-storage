// +build tools

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"text/template"

	"github.com/Xuanwo/templateutils"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	log "github.com/sirupsen/logrus"
)

var (
	infoT      = newTmpl("cmd/definitions/tmpl/info")
	pairT      = newTmpl("cmd/definitions/tmpl/pair")
	serviceT   = newTmpl("cmd/definitions/tmpl/service")
	operationT = newTmpl("cmd/definitions/tmpl/operation")
	functionT  = newTmpl("cmd/definitions/tmpl/function")
	objectT    = newTmpl("cmd/definitions/tmpl/object")
)

func generateGlobal(data *Data) {
	// Metas generate
	generateT(infoT, "types/info.generated.go", data)

	// Pair generate
	generateT(pairT, "pairs/generated.go", data)

	// Operation generate
	generateT(operationT, "types/operation.generated.go", data)

	// Object generate
	generateT(objectT, "types/object.generated.go", data)
}

func generateService(data *Data) {
	generateT(serviceT, "generated.go", data.Service)
	for _, v := range data.Service.Namespaces {
		appendT(functionT, v.Name+".go", v)
		formatService(v.Name + ".go")
	}
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

func formatService(filename string) {
	fset := token.NewFileSet()

	f, err := decorator.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}

	// Sort all methods via name.
	// Only allow functions in service.go and storage.go
	sort.SliceStable(f.Decls, func(i, j int) bool {
		fi, ok := f.Decls[i].(*dst.FuncDecl)
		if !ok {
			return false
		}
		fj, ok := f.Decls[j].(*dst.FuncDecl)
		if !ok {
			return false
		}
		return fi.Name.Name < fj.Name.Name
	})

	// Add empty line for every methods.
	dst.Inspect(f, func(node dst.Node) bool {
		if node == nil {
			return false
		}
		fn, ok := node.(*dst.FuncDecl)
		if !ok {
			return true
		}
		fn.Decorations().Before = dst.EmptyLine
		fn.Decorations().After = dst.EmptyLine
		return true
	})

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("create file %v: %v", filename, err)
	}

	err = decorator.Fprint(file, f)
	if err != nil {
		log.Fatalf("format file %v: %v", filename, err)
	}
}
