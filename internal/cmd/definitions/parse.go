package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Xuanwo/templateutils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func parse() (data *Data) {
	data = &Data{}

	// Parse pairs
	files, err := ioutil.ReadDir("pairs")
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	for _, v := range files {
		if !strings.HasSuffix(v.Name(), ".hcl") {
			continue
		}

		filePath := "pairs/" + v.Name()
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}

		ps := struct {
			Pair []*Pair `hcl:"pair,block"`
		}{
			Pair: make([]*Pair, 0),
		}
		err = parseHCL(content, filePath, &ps)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}
		for _, v := range ps.Pair {
			v.Description = strings.ReplaceAll(v.Description, "\n", "\n//")
		}

		data.Pairs = append(data.Pairs, ps.Pair...)
	}
	sort.Slice(data.Pairs, func(i, j int) bool {
		return data.Pairs[i].Name < data.Pairs[j].Name
	})

	data.TypeMap = make(map[string]string)
	for _, v := range data.Pairs {
		data.TypeMap[v.Name] = v.Type
	}

	// Parse metadata
	files, err = ioutil.ReadDir("metadata")
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	for _, v := range files {
		if !strings.HasSuffix(v.Name(), ".hcl") {
			continue
		}

		filePath := "metadata/" + v.Name()
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}

		ms := struct {
			ObjectMeta       []*Metadata `hcl:"object_meta,block"`
			StorageMeta      []*Metadata `hcl:"storage_meta,block"`
			StorageStatistic []*Metadata `hcl:"storage_statistic,block"`
		}{
			ObjectMeta:       make([]*Metadata, 0),
			StorageMeta:      make([]*Metadata, 0),
			StorageStatistic: make([]*Metadata, 0),
		}
		err = parseHCL(content, filePath, &ms)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}

		data.ObjectMeta = append(data.ObjectMeta, ms.ObjectMeta...)
		data.StorageMeta = append(data.StorageMeta, ms.StorageMeta...)
		data.StorageStatistic = append(data.StorageStatistic, ms.StorageStatistic...)
	}
	sort.Slice(data.ObjectMeta, func(i, j int) bool {
		return data.ObjectMeta[i].Name < data.ObjectMeta[j].Name
	})
	sort.Slice(data.StorageMeta, func(i, j int) bool {
		return data.StorageMeta[i].Name < data.StorageMeta[j].Name
	})
	sort.Slice(data.StorageStatistic, func(i, j int) bool {
		return data.StorageStatistic[i].Name < data.StorageStatistic[j].Name
	})

	// Parse service
	files, err = ioutil.ReadDir("services")
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	for _, v := range files {
		if !strings.HasSuffix(v.Name(), ".hcl") {
			continue
		}

		filePath := "services/" + v.Name()
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}

		srv := &Service{}
		err = parseHCL(content, filePath, srv)
		if err != nil {
			log.Fatalf("parse: %v", err)
		}

		servicerFuncs, storagerFuncs := parseFunc(srv.Name, "servicer"), parseFunc(srv.Name, "storager")

		// Register funcs into service
		for _, v := range srv.Service.Op {
			if fn, ok := servicerFuncs[v.Op]; ok {
				fn.hasPair = true
				v.Func = fn
			}
		}
		for _, v := range srv.Storage.Op {
			if fn, ok := storagerFuncs[v.Op]; ok {
				fn.hasPair = true
				v.Func = fn
			}
		}

		// Add missing pairs into service
		for k, v := range servicerFuncs {
			if v.hasPair {
				continue
			}
			srv.Service.Op = append(srv.Service.Op, &Op{
				Op:   k,
				Func: v,
			})
		}
		for k, v := range storagerFuncs {
			if v.hasPair {
				continue
			}
			srv.Storage.Op = append(srv.Storage.Op, &Op{
				Op:   k,
				Func: v,
			})
		}

		// Inject paris
		injectReadCallbackFunc(srv.Storage)
		injectContext(srv.Service)
		injectContext(srv.Storage)

		// Create type map
		srv.TypeMap = data.TypeMap

		for _, v := range srv.Service.Op {
			sort.Strings(v.Optional)
			sort.Strings(v.Required)
		}
		sort.Slice(srv.Service.Op, func(i, j int) bool {
			return srv.Service.Op[i].Op < srv.Service.Op[j].Op
		})
		for _, v := range srv.Storage.Op {
			sort.Strings(v.Optional)
			sort.Strings(v.Required)
		}
		sort.Slice(srv.Storage.Op, func(i, j int) bool {
			return srv.Storage.Op[i].Op < srv.Storage.Op[j].Op
		})

		data.Service = append(data.Service, srv)
	}

	return
}

func injectReadCallbackFunc(ops *Ops) {
	for _, op := range ops.Op {
		fn := op.Func
		if fn == nil {
			continue
		}
		if strings.Contains(fn.Params, "io.Reader") ||
			strings.Contains(fn.Returns, "io.ReadCloser") {
			op.Optional = append(op.Optional, "read_callback_func")
		}
	}
}

func injectContext(ops *Ops) {
	for _, op := range ops.Op {
		op.Optional = append(op.Optional, "context")
	}
}

func parseHCL(src []byte, filename string, in interface{}) (err error) {
	var diags hcl.Diagnostics
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse hcl: %w", err)
		}
	}()

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return diags
	}

	diags = gohcl.DecodeBody(file.Body, nil, in)
	if diags.HasErrors() {
		return diags
	}

	return nil
}

func parseFunc(service, name string) map[string]*Func {
	data := make(map[string]*Func)
	filename := fmt.Sprintf("../services/%s/%s.go", service, name)

	content, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return data
	}
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
		// Ignore non-exported funcs.
		if !fndecl.Name.IsExported() {
			continue
		}
		// Ignore some not needed functions.
		if fndecl.Name.Name == "String" {
			continue
		}

		opName := templateutils.ToSnack(fndecl.Name.Name)

		data[opName] = &Func{
			Parent:   fndecl.Name.Name,
			Receiver: getReceiver(fndecl),
			Returns:  getReturns(fndecl),
			Caller:   getCaller(fndecl),
			Params:   getParams(fndecl),
		}

		if fndecl.Recv != nil {
			data[opName].HasContext = true
		}
	}
	return data
}

func getReceiver(fn *ast.FuncDecl) string {
	if fn.Recv == nil {
		return ""
	}
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
