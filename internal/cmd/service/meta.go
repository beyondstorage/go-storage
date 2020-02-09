package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/Xuanwo/templateutils"
)

const (
	metaPath  = "meta.json"
	pairsPath = "../../types/pairs/pairs.json"
)

type metadata struct {
	Name    string                     `json:"name"`
	Service map[string]map[string]bool `json:"service,omitempty"`
	Storage map[string]map[string]bool `json:"storage"`

	TypeMap map[string]string         `json:"-"`
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

	// Inject ReadCallbackFunc in all operations who have a Reader.
	if _, ok := meta.Storage["read"]; !ok {
		meta.Storage["read"] = make(map[string]bool)
	}
	meta.Storage["read"]["read_callback_func"] = false
	if _, ok := meta.Storage["write"]; !ok {
		meta.Storage["write"] = make(map[string]bool)
	}
	meta.Storage["write"]["read_callback_func"] = false
	// TODO: WriteSegment should also be supported.

	// Handle Data
	meta.Data = make(map[string]map[string]*fn)
	meta.Data["service"] = mergeFn(meta.Service, parseFunc("servicer"), parseFunc("utils"))
	meta.Data["storage"] = mergeFn(meta.Storage, parseFunc("storager"))

	return meta
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
