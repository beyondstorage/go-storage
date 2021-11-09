package definitions

import (
	"fmt"
)

type Field struct {
	Name string
	Type Type
}

func getField(name string) Field {
	f, ok := fieldMap[name]
	if !ok {
		panic(fmt.Errorf("field %s is not exist", name))
	}
	return f
}

var fieldMap = buildFieldMap()

var fieldArray = []Field{
	{
		Name: "bi",
		Type: Type{Expr: "*", Package: "types", Name: "BlockIterator"}},
	{
		Name: "bid",
		Type: Type{Name: "string"},
	},
	{
		Name: "bids",
		Type: Type{Expr: "[]", Name: "string"},
	},
	{
		Name: "ctx",
		Type: Type{Package: "context", Name: "Context"},
	},
	{
		Name: "dst",
		Type: Type{Name: "string"},
	},
	{
		Name: "err",
		Type: Type{Name: "error"},
	},
	{
		Name: "expire",
		Type: Type{Package: "time", Name: "Duration"},
	},
	{
		Name: "index",
		Type: Type{Name: "int"},
	},
	{
		Name: "meta",
		Type: Type{Expr: "*", Package: "types", Name: "StorageMeta"},
	},
	{
		Name: "n",
		Type: Type{Name: "int64"},
	},
	{
		Name: "name",
		Type: Type{Name: "string"},
	},
	{
		Name: "o",
		Type: Type{Expr: "*", Package: "types", Name: "Object"},
	},
	{
		Name: "offset",
		Type: Type{Name: "int64"},
	},
	{
		Name: "oi",
		Type: Type{Expr: "*", Package: "types", Name: "ObjectIterator"},
	},
	{
		Name: "op",
		Type: Type{Name: "string"},
	},
	{
		Name: "pairs",
		Type: Type{Expr: "...", Package: "types", Name: "Pair"},
	},
	{
		Name: "part",
		Type: Type{Expr: "*", Package: "types", Name: "Part"},
	},
	{
		Name: "parts",
		Type: Type{Expr: "[]*", Package: "types", Name: "Part"},
	},
	{
		Name: "path",
		Type: Type{Name: "string"},
	},
	{
		Name: "pi",
		Type: Type{Expr: "*", Package: "types", Name: "PartIterator"},
	},
	{
		Name: "r",
		Type: Type{Package: "io", Name: "Reader"},
	},
	{
		Name: "req",
		Type: Type{Expr: "*", Package: "http", Name: "Request"},
	},
	{
		Name: "srvf",
		Type: Type{Package: "types", Name: "ServiceFeatures"},
	},
	{
		Name: "size",
		Type: Type{Name: "int64"},
	},
	{
		Name: "src",
		Type: Type{Name: "string"},
	},
	{
		Name: "url",
		Type: Type{Name: "string"},
	},
	{
		Name: "sti",
		Type: Type{Expr: "*", Package: "types", Name: "StoragerIterator"},
	},
	{
		Name: "stof",
		Type: Type{Package: "types", Name: "StorageFeatures"},
	},
	{
		Name: "store",
		Type: Type{Package: "types", Name: "Storager"},
	},
	{
		Name: "w",
		Type: Type{Package: "io", Name: "Writer"},
	},
	{
		Name: "target",
		Type: Type{Name: "string"},
	},
}

func buildFieldMap() map[string]Field {
	m := make(map[string]Field)
	for _, v := range fieldArray {
		v := v
		m[v.Name] = v
	}
	return m
}
