//go:build tools
// +build tools

package main

import (
	"sort"

	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

func generatePair(data *Data, path string) {
	f := gg.NewGroup()
	f.AddLineComment("Code generated by go generate cmd/definitions; DO NOT EDIT.")
	f.AddPackage("pairs")
	f.NewImport().
		AddPath("context").
		AddPath("time").
		AddLine().
		AddPath("github.com/beyondstorage/go-storage/v4/pkg/httpclient").
		AddDot("github.com/beyondstorage/go-storage/v4/types")

	ps := make([]*Pair, 0, len(data.PairsMap))
	for _, v := range data.PairsMap {
		v := v
		ps = append(ps, v)
	}
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Name < ps[j].Name
	})

	for _, v := range ps {
		pname := templateutils.ToPascal(v.Name)

		f.AddLineComment(`With%s will apply %s value to Options.

%s %s`, pname, v.Name, pname, v.Description)
		xfn := f.NewFunction("With" + pname)
		xfn.AddParameter("v", v.Type)
		xfn.AddResult("p", "Pair")
		xfn.AddBody(
			gg.Return(
				gg.Value("Pair").
					AddField("Key", gg.Lit(v.Name)).
					AddField("Value", "v")))
	}

	err := f.WriteFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}
