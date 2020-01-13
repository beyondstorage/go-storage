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

	TypeMap map[string]string   `json:"-"`
	Data    map[string]receiver `json:"-"`
}

type receiver struct {
	Pairs map[string]map[string]bool
	Funcs map[string]*contextFunc
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

	// Handle TypeMap
	pairsPath := "../../types/pairs/pairs.json"
	content, err = ioutil.ReadFile(pairsPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}
	err = json.Unmarshal(content, &meta.TypeMap)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	// Handle Data
	meta.Data = make(map[string]receiver)
	meta.Data["service"] = receiver{
		Pairs: meta.Service,
		Funcs: parseFunc("service"),
	}
	for k := range meta.Service {
		// If func not implemented, remove.
		if _, ok := meta.Data["service"].Funcs[templateutils.ToPascal(k)]; !ok {
			delete(meta.Service, k)
		}
		// If no paris, remove.
		if len(meta.Service[k]) == 0 {
			delete(meta.Service, k)
		}
	}
	meta.Data["storage"] = receiver{
		Pairs: meta.Storage,
		Funcs: parseFunc("storage"),
	}
	for k := range meta.Storage {
		// If func not implemented, remove.
		if _, ok := meta.Data["storage"].Funcs[templateutils.ToPascal(k)]; !ok {
			delete(meta.Storage, k)
		}
		// If no paris, remove.
		if len(meta.Storage[k]) == 0 {
			delete(meta.Storage, k)
		}
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

	return meta
}
