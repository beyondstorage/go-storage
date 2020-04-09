package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Xuanwo/templateutils"
)

const (
	metaPath  = "meta.json"
	pairsPath = "../../types/pairs/pairs.json"
)

type pairs map[string]struct {
	Type        string
	Description string
}

type metadata struct {
	Name    string                     `json:"name"`
	Service map[string]map[string]bool `json:"service,omitempty"`
	Storage map[string]map[string]bool `json:"storage"`

	TypeMap pairs                     `json:"-"`
	Data    map[string]map[string]*fn `json:"-"`
}

type fn struct {
	Pairs map[string]bool
	Funcs *contextFunc
}

func parseMeta() metadata {
	if _, err := os.Stat(metaPath); err != nil {
		log.Fatalf("stat meta failed: %v", err)
	}

	content, err := ioutil.ReadFile(metaPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	var meta metadata

	err = json.Unmarshal(content, &meta)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	// Format input meta.json
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(metaPath, data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	// Handle TypeMap
	content, err = ioutil.ReadFile(pairsPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}
	err = json.Unmarshal(content, &meta.TypeMap)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	// Get all funcs
	servicerFuncs, storagerFuncs, utilsFuncs := parseFunc("servicer"), parseFunc("storager"), parseFunc("utils")

	injectReadCallbackFunc(meta.Storage, storagerFuncs)
	// Handle Data
	meta.Data = make(map[string]map[string]*fn)
	meta.Data["service"] = mergeFn(meta.Service, servicerFuncs, utilsFuncs)
	meta.Data["storage"] = mergeFn(meta.Storage, storagerFuncs)

	return meta
}

// Inject ReadCallbackFunc in all operations who have a Reader.
func injectReadCallbackFunc(mp map[string]map[string]bool, fns map[string]*contextFunc) {
	funcs := make([]string, 0)
	for funcName, fn := range fns {
		funcName = templateutils.ToSnack(funcName)
		if strings.Contains(fn.Params, "io.Reader") ||
			strings.Contains(fn.Returns, "io.ReadCloser") {
			funcs = append(funcs, funcName)
		}
	}

	for _, funcName := range funcs {
		if _, ok := mp[funcName]; !ok {
			mp[funcName] = make(map[string]bool)
		}
		mp[funcName]["read_callback_func"] = false
	}
}

func mergeFn(mp map[string]map[string]bool, mfn ...map[string]*contextFunc) map[string]*fn {
	m := make(map[string]*fn)
	for _, mp := range mfn {
		for k, v := range mp {
			v := v
			k = templateutils.ToKebab(k)

			m[k] = &fn{
				Funcs: v,
			}
		}
	}
	for k, v := range mp {
		v := v
		k = templateutils.ToKebab(k)

		if _, ok := m[k]; !ok {
			m[k] = &fn{}
		}
		m[k].Pairs = v
	}
	return m
}
