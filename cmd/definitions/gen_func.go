//go:build tools
// +build tools

package main

import (
	"github.com/Xuanwo/gg"
	"github.com/Xuanwo/templateutils"
	log "github.com/sirupsen/logrus"
)

func generateFunc(ns *Namespace, path string) {
	f := gg.Group()

	nsNameP := templateutils.ToPascal(ns.Name)

	for _, fn := range ns.Funcs {
		if fn.Implemented {
			continue
		}

		fnNameC := templateutils.ToCamel(fn.Name)
		fnNameP := templateutils.ToPascal(fn.Name)

		if fn.Local {
			gfn := f.Function(fnNameC).
				Receiver("s", "*"+nsNameP)
			for _, v := range fn.Params {
				// We need to remove pair from generated functions.
				if v.Type() == "...Pair" {
					continue
				}
				gfn.Parameter(v.Name, v.Type())
			}
			gfn.Parameter("opt", "pair"+nsNameP+fnNameP)
			for _, v := range fn.Results {
				gfn.Result(v.Name, v.Type())
			}
			gfn.Body(gg.String(`panic("not implemented")`))
			continue
		}

		gfn := f.Function(fnNameC).
			Receiver("s", "*"+nsNameP)
		gfn.Parameter("ctx", "context.Context")
		for _, v := range fn.Params {
			// We need to remove pair from generated functions.
			if v.Type() == "...Pair" {
				continue
			}
			gfn.Parameter(v.Name, v.Type())
		}
		gfn.Parameter("opt", "pair"+nsNameP+fnNameP)
		for _, v := range fn.Results {
			gfn.Result(v.Name, v.Type())
		}
		gfn.Body(gg.String(`panic("not implemented")`))
		continue
	}

	err := f.AppendFile(path)
	if err != nil {
		log.Fatalf("generate to %s: %v", path, err)
	}
}
