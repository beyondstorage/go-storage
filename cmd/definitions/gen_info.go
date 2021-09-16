//go:build tools
// +build tools

package main

import (
	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

func generateInfo(data *Data, path string) {
	f := gg.NewGroup()
	f.AddLineComment("Code generated by go generate cmd/definitions; DO NOT EDIT.")
	f.AddPackage("types")
	f.NewImport().AddPath("fmt")

	serviceName := "storage-meta"
	infos := data.StorageMeta()

	serviceNameC := templateutils.ToCamel(serviceName)
	serviceNameP := templateutils.ToPascal(serviceName)
	serviceNameK := templateutils.ToKebab(serviceName)

	// Generate field bits
	f.AddLineComment("Field index in %s bit", serviceNameC)
	consts := f.NewConst()
	for k, v := range infos {
		consts.AddField(
			gg.S("%sIndex%s",
				serviceNameC, templateutils.ToPascal(v.Name)),
			gg.S("1<<%d", k),
		)
	}

	// Generate struct
	st := f.NewStruct(serviceNameP)
	for _, v := range infos {
		st.AddField(v.TypeName(), v.Type)
	}
	st.AddLine()
	st.AddLineComment("bit used as a bitmap for object value, 0 means not set, 1 means set")
	st.AddField("bit", "uint64")
	st.AddField("m", "map[string]interface{}")

	// Generate Get/Set functions.
	for _, v := range infos {
		// If the value is export, we don't need to generate MustGetXxxx
		if v.Export {
			f.NewFunction("Get"+v.DisplayName()).
				WithReceiver("m", "*"+serviceNameP).
				AddResult("", v.Type).
				AddBody(gg.Return("m." + v.TypeName()))
			f.NewFunction("Set"+v.DisplayName()).
				WithReceiver("m", "*"+serviceNameP).
				AddParameter("v", v.Type).
				AddResult("", "*"+serviceNameP).
				AddBody(
					gg.S("m.%s = v", v.TypeName()),
					gg.Return("m"),
				)
			continue
		}
		f.NewFunction("Get"+v.DisplayName()).
			WithReceiver("m", "*"+serviceNameP).
			AddResult("", v.Type).
			AddResult("", "bool").
			AddBody(
				gg.If(gg.S("m.bit & %sIndex%s != 0",
					serviceNameC, templateutils.ToPascal(v.Name))).
					AddBody(gg.Return("m."+v.TypeName(), gg.Lit(true))),
				gg.Return(templateutils.ZeroValue(v.Type), gg.Lit(false)),
			)
		f.NewFunction("MustGet"+v.DisplayName()).
			WithReceiver("m", "*"+serviceNameP).
			AddResult("", v.Type).
			AddBody(
				gg.If(gg.S("m.bit & %sIndex%s == 0",
					serviceNameC, templateutils.ToPascal(v.Name))).
					AddBody(gg.S(
						`panic(fmt.Sprintf("%s %s is not set"))`,
						serviceNameK, v.Name)),
				gg.Return("m."+v.TypeName()))
		f.NewFunction("Set"+v.DisplayName()).
			WithReceiver("m", "*"+serviceNameP).
			AddParameter("v", v.Type).
			AddResult("", "*"+serviceNameP).
			AddBody(
				gg.S("m.%s = v", v.TypeName()),
				gg.S("m.bit |= %sIndex%s", serviceNameC, templateutils.ToPascal(v.Name)),
				gg.Return("m"),
			)
	}

	err := f.WriteFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}
