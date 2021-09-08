//go:build tools
// +build tools

package main

import (
	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

func generateFunc(ns *Namespace, path string) {
	f := gg.NewGroup()

	nsNameP := templateutils.ToPascal(ns.Name)

	for _, fn := range ns.Funcs {
		if fn.Implemented {
			continue
		}

		fnNameC := templateutils.ToCamel(fn.Name)
		fnNameP := templateutils.ToPascal(fn.Name)

		if fn.Local {
			gfn := f.NewFunction(fnNameC).
				WithReceiver("s", "*"+nsNameP)
			for _, v := range fn.Params {
				// We need to remove pair from generated functions.
				if v.Type() == "...Pair" {
					continue
				}
				gfn.AddParameter(v.Name, v.Type())
			}
			gfn.AddParameter("opt", "pair"+nsNameP+fnNameP)
			for _, v := range fn.Results {
				gfn.AddResult(v.Name, v.Type())
			}
			gfn.AddBody(gg.S(`panic("not implemented")`))
			continue
		}

		gfn := f.NewFunction(fnNameC).
			WithReceiver("s", "*"+nsNameP)
		gfn.AddParameter("ctx", "context.Context")
		for _, v := range fn.Params {
			// We need to remove pair from generated functions.
			if v.Type() == "...Pair" {
				continue
			}
			gfn.AddParameter(v.Name, v.Type())
		}
		gfn.AddParameter("opt", "pair"+nsNameP+fnNameP)
		for _, v := range fn.Results {
			gfn.AddResult(v.Name, v.Type())
		}
		gfn.AddBody(`panic("not implemented")`)
		continue
	}

	err := f.AppendFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}
