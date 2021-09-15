//go:build tools
// +build tools

package main

import (
	"go/parser"
	"go/token"
	"os"
	"sort"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	log "github.com/sirupsen/logrus"
)

func generateGlobal(data *Data) {
	// Iterator generate
	generateIterator("types/iterator.generated.go")

	// Metas generate
	generateInfo(data, "types/info.generated.go")

	// Pair generate
	generatePair(data, "pairs/generated.go")

	// Operation generate
	generateOperation(data, "types/operation.generated.go")

	// Object generate
	generateObject(data, "types/object.generated.go")
}

func generateService(data *Data) {
	generateSrv(data.Service, "generated.go")
	for _, v := range data.Service.SortedNamespaces() {
		generateFunc(v, v.Name+".go")
		formatService(v.Name + ".go")
	}
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
