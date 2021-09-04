//go:build tools
// +build tools

package main

import (
	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

func generateObject(data *Data, path string) {
	f := gg.Group()
	f.LineComment("Code generated by go generate cmd/definitions; DO NOT EDIT.")
	f.Package("types")
	f.Imports().
		Path("fmt").
		Path("time").
		Path("sync")

	f.LineComment("Field index in object bit")
	cons := f.Const()
	for k, v := range data.ObjectMeta {
		pname := templateutils.ToPascal(v.Name)
		cons.TypedField(
			"objectIndex"+pname, "uint64", gg.S("1<<%d", k))
	}

	f.LineComment(`
Object is the smallest unit in go-storage.

NOTES:
  - Object's fields SHOULD not be changed outside services.
  - Object CANNOT be copied
  - Object is concurrent safe.
  - Only ID, Path, Mode are required during list operations, other fields
    could be fetched via lazy stat logic: https://beyondstorage.io/docs/go-storage/internal/object-lazy-stat
`)
	ob := f.Struct("Object")
	for _, v := range data.ObjectMeta {
		if v.Description != "" {
			ob.LineComment(v.Description)
		}
		ob.Field(v.TypeName(), v.Type())
	}

	ob.LineComment("client is the client in which Object is alive.")
	ob.Field("client", "Storager")

	ob.LineComment("bit used as a bitmap for object value, 0 means not set, 1 means set")
	ob.Field("bit", "uint64")
	ob.Field("done", "uint32")
	ob.Field("m", "sync.Mutex")

	for _, v := range data.ObjectMeta {
		pname := templateutils.ToPascal(v.Name)

		if v.Export {
			f.Function("Get"+v.DisplayName()).
				Receiver("o", "*Object").
				Result("", v.Type()).
				Body(gg.Return(gg.S("o.%s", v.TypeName())))
			f.Function("Set"+v.DisplayName()).
				Receiver("o", "*Object").
				Parameter("v", v.Type()).
				Result("", "*Object").
				Body(
					gg.S("o.%s = v", v.TypeName()),
					gg.Return("o"),
				)
			continue
		}
		f.Function("Get"+v.DisplayName()).
			Receiver("o", "*Object").
			NamedLineComment(`will get %s from Object.

%s
`, v.DisplayName(), v.Description).
			Result("", v.Type()).
			Result("", "bool").
			Body(
				gg.S("o.stat()"),
				gg.Line(),
				gg.If(gg.S("o.bit & objectIndex%s != 0", pname)).
					Body(
						gg.Return("o."+v.TypeName(), gg.Lit(true)),
					),
				gg.Return(templateutils.ZeroValue(v.Type()), gg.Lit(false)),
			)
		f.Function("MustGet"+v.DisplayName()).
			Receiver("o", "*Object").
			NamedLineComment(`will get %s from Object.

%s
`, v.DisplayName(), v.Description).
			Result("", v.Type()).
			Body(
				gg.S("o.stat()"),
				gg.Line(),
				gg.If(gg.S("o.bit & objectIndex%s == 0", pname)).
					Body(
						gg.S(`panic(fmt.Sprintf("object %s is not set"))`, v.Name),
					),
				gg.Return("o."+v.TypeName()),
			)
		f.Function("Set"+v.DisplayName()).
			Receiver("o", "*Object").
			NamedLineComment(`will set %s into Object.

%s
`, v.DisplayName(), v.Description).
			Parameter("v", v.Type()).
			Result("", "*Object").
			Body(
				gg.S("o.%s = v", v.TypeName()),
				gg.S("o.bit |= objectIndex%s", pname),
				gg.Return("o"),
			)
	}
	fn := f.Function("clone").
		Receiver("o", "*Object").
		Parameter("xo", "*Object")
	for _, v := range data.ObjectMeta {
		fn.Body(gg.S("o.%s = xo.%s", v.TypeName(), v.TypeName()))
	}
	fn.Body(gg.S("o.bit = xo.bit"))

	err := f.WriteFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}
