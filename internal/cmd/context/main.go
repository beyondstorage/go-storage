package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

var (
	contextT = template.Must(
		template.New("context").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("context.tmpl"))))
)

type contextFunc struct {
	Parent   string // Old method name: "AbortSegment"
	Receiver string // Receiver's name: "s *Storage"
	Params   string // Method's Params: "ctx context.Context, id string, pairs ...*types.Pair"
	Returns  string // Method's returns: "err error"
	Caller   string // How to call Parent method: "id, pairs..."
}

//go:generate go-bindata -nometadata -ignore ".*.go" .
func main() {
	files := []string{"Servicer", "Storager"}
	data := map[string]map[string]*contextFunc{
		"Servicer": make(map[string]*contextFunc),
		"Storager": make(map[string]*contextFunc),
	}

	for _, name := range files {
		filename := strings.ToLower(name) + ".go"
		_, err := os.Stat(filename)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("read file failed: %v", err)
		}

		f, err := parser.ParseFile(token.NewFileSet(), filename, string(content), 0)
		if err != nil {
			log.Fatalf("decorator parse failed: %v", err)
		}

		for _, fn := range f.Decls {
			fndecl, ok := fn.(*ast.FuncDecl)
			// Ignore Non-FuncDecl node.
			if !ok {
				continue
			}
			// Ignore functions.
			if fndecl.Recv == nil {
				continue
			}
			// Ignore non-exported funcs.
			if !fndecl.Name.IsExported() {
				continue
			}
			// Ignore some not needed functions.
			if fndecl.Name.Name == "String" {
				continue
			}

			data[name][fndecl.Name.Name] = &contextFunc{
				Parent:   fndecl.Name.Name,
				Receiver: getReceiver(fndecl),
				Returns:  getReturns(fndecl),
				Caller:   getCaller(fndecl),
			}

			// Add context
			fndecl.Type.Params.List = append([]*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("ctx"),
					},
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("context"),
						Sel: ast.NewIdent("Context"),
					},
				},
			}, fndecl.Type.Params.List...)
			data[name][fndecl.Name.Name].Params = getParams(fndecl)
		}
	}

	filePath := "context.go"
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = contextT.Execute(f, struct {
		Name string
		Data map[string]map[string]*contextFunc
	}{
		Name: os.Getenv("GOPACKAGE"),
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getReceiver(fn *ast.FuncDecl) string {
	return fmt.Sprintf("s %s", formatExpr(fn.Recv.List[0].Type))
}

func getParams(fn *ast.FuncDecl) string {
	parms := []string{}
	for _, v := range fn.Type.Params.List {
		parms = append(parms, formatField(v))
	}
	ans := fmt.Sprintf("%s", strings.Join(parms, ","))
	return ans
}

func getReturns(fn *ast.FuncDecl) string {
	results := []string{}
	for _, v := range fn.Type.Results.List {
		results = append(results, formatField(v))
	}
	ans := fmt.Sprintf("%s", strings.Join(results, ","))
	return ans
}

func getCaller(fn *ast.FuncDecl) string {
	parms := []string{}
	for _, v := range fn.Type.Params.List {
		if _, ok := v.Type.(*ast.Ellipsis); ok {
			parms = append(parms, v.Names[0].Name+"...")
			continue
		}
		for _, name := range v.Names {
			parms = append(parms, name.Name)
		}
	}
	ans := fmt.Sprintf("%s", strings.Join(parms, ","))
	return ans
}

func formatField(f *ast.Field) string {
	s := []string{}
	for _, name := range f.Names {
		s = append(s, name.Name)
	}
	return strings.Join(s, ",") + " " + formatExpr(f.Type)
}

func formatExpr(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", formatExpr(v.X), v.Sel.Name)
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return "*" + formatExpr(v.X)
	case *ast.Ellipsis:
		return "..." + formatExpr(v.Elt)
	default:
		println(fmt.Sprintf("not handled type %+#v", v))
		return ""
	}
}
