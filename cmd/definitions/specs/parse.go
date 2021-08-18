// +build tools

package specs

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"

	"github.com/beyondstorage/go-storage/v4/cmd/definitions/bindata"
)

type tomlFeature struct {
	Description string `toml:"description"`
}

type tomlFeatures map[string]tomlFeature

type tomlField struct {
	Type string
}

type tomlFields map[string]tomlField

type tomlInfo struct {
	Type        string `toml:"type"`
	Export      bool   `toml:"export"`
	Description string `toml:"description"`
}

type tomlInfos map[string]tomlInfo

type tomlPair struct {
	Type        string `toml:"type"`
	Description string `toml:"description,optional"`
}

type tomlPairs map[string]tomlPair

type tomlOperation struct {
	Description string   `toml:"description"`
	Params      []string `toml:"params"`
	Pairs       []string `toml:"pairs"`
	Results     []string `toml:"results"`
	ObjectMode  string   `toml:"object_mode"`
	Local       bool     `toml:"local"`
}

type tomlInterface struct {
	Description string                   `toml:"description"`
	Ops         map[string]tomlOperation `toml:"op"`
}

type tomlInterfaces map[string]tomlInterface

type tomlOp struct {
	Required []string `toml:"required"`
	Optional []string `toml:"optional"`
}

type tomlNamespace struct {
	Features    []string          `toml:"features"`
	Implement   []string          `toml:"implement"`
	Defaultable []string          `toml:"defaultable"`
	New         tomlOp            `toml:"new"`
	Op          map[string]tomlOp `toml:"op"`
}

type tomlService struct {
	Name      string                                    `toml:"name"`
	Namespace map[string]tomlNamespace                  `toml:"namespace"`
	Pairs     map[string]tomlPair                       `toml:"pairs"`
	Infos     map[string]map[string]map[string]tomlInfo `toml:"infos"`
}

func parseTOML(src []byte, in interface{}) (err error) {
	return toml.Unmarshal(src, in)
}

func parseFeatures() Features {
	var tp tomlFeatures

	err := parseTOML(bindata.MustAsset(featurePath), &tp)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	var ps Features

	for k, v := range tp {
		p := Feature{
			Name:        k,
			Description: v.Description,
		}

		ps = append(ps, p)
	}

	return ps
}

func parsePairs() Pairs {
	var tp tomlPairs

	err := parseTOML(bindata.MustAsset(pairPath), &tp)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	var ps Pairs

	for k, v := range tp {
		p := Pair{
			Name:        k,
			Type:        v.Type,
			Description: v.Description,
		}

		ps = append(ps, p)
	}

	return ps
}

func parseInfos() Infos {
	var ti tomlInfos
	var ps Infos

	// Parse object meta
	err := parseTOML(bindata.MustAsset(infoObjectMeta), &ti)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	for k, v := range ti {
		p := Info{
			Scope:       "object",
			Category:    "meta",
			Name:        k,
			Type:        v.Type,
			Export:      v.Export,
			Description: v.Description,
		}

		ps = append(ps, p)
	}

	// Parse storage meta
	err = parseTOML(bindata.MustAsset(infoStorageMeta), &ti)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	for k, v := range ti {
		p := Info{
			Scope:       "storage",
			Category:    "meta",
			Name:        k,
			Type:        v.Type,
			Export:      v.Export,
			Description: v.Description,
		}

		ps = append(ps, p)
	}

	return ps
}

func parseOperations() Operations {
	var (
		ti tomlInterfaces
		tf tomlFields
		ps Operations
	)

	// Parse operations.
	err := parseTOML(bindata.MustAsset(operationPath), &ti)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	// Parse fields
	err = parseTOML(bindata.MustAsset(fieldPath), &tf)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}

	for k, v := range tf {
		ps.Fields = append(ps.Fields, Field{
			Name: k,
			Type: v.Type,
		})
	}

	for k, v := range ti {
		it := Interface{
			Name:        k,
			Description: v.Description,
		}

		for k, v := range v.Ops {
			it.Ops = append(it.Ops, Operation{
				Name:        k,
				Description: v.Description,
				Params:      v.Params,
				Pairs:       v.Pairs,
				Results:     v.Results,
				ObjectMode:  v.ObjectMode,
				Local:       v.Local,
			})
		}

		ps.Interfaces = append(ps.Interfaces, it)
	}

	return ps
}

func parseService(filePath string) Service {
	var (
		ts tomlService
		ps Service
	)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read file %s: %v", filePath, err)
	}

	err = parseTOML(content, &ts)
	if err != nil {
		log.Fatalf("parse toml %s: %v", filePath, err)
	}

	ps.Name = ts.Name

	// Parse pairs
	for k, v := range ts.Pairs {
		ps.Pairs = append(ps.Pairs, Pair{
			Name:        k,
			Type:        v.Type,
			Description: v.Description,
		})
	}

	// Parse infos
	for scope, v := range ts.Infos {
		for category, v := range v {
			for name, v := range v {
				ps.Infos = append(ps.Infos, Info{
					Scope:       scope,
					Category:    category,
					Name:        name,
					Type:        v.Type,
					Export:      v.Export,
					Description: v.Description,
				})
			}
		}
	}

	// Parse namespace.
	for name, v := range ts.Namespace {
		n := Namespace{
			Name:        name,
			Implement:   v.Implement,
			Features:    v.Features,
			Defaultable: v.Defaultable,
			New: New{
				Required: v.New.Required,
				Optional: v.New.Optional,
			},
		}

		for opName, op := range v.Op {
			n.Op = append(n.Op, Op{
				Name:     opName,
				Required: op.Required,
				Optional: op.Optional,
			})
		}

		ps.Namespaces = append(ps.Namespaces, n)
	}

	return ps
}
