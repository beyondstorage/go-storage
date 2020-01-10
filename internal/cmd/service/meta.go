package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	meta.Data["storage"] = receiver{
		Pairs: meta.Storage,
		Funcs: parseFunc("storage"),
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
