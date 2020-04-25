package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Xuanwo/templateutils"
)

var (
	pairsT = template.Must(
		template.New("meta").
			Funcs(templateutils.FuncMap()).
			Parse(string(MustAsset("pair.tmpl"))))
)

// Pairs is the struct for pairs.json
type Pairs map[string]struct {
	Type        string
	Description string
}

//go:generate go-bindata -nometadata -ignore "\\.go$" .
func main() {
	pairsPath := "pairs.json"
	content, err := ioutil.ReadFile(pairsPath)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	var pairs Pairs
	err = json.Unmarshal(content, &pairs)
	if err != nil {
		log.Fatalf("json unmarshal failed: %v", err)
	}

	// Format input meta.json
	data, err := json.MarshalIndent(pairs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(pairsPath, data, 0664)
	if err != nil {
		log.Fatal(err)
	}

	pairsFile, err := os.Create("pairs.go")
	if err != nil {
		log.Fatal(err)
	}
	defer pairsFile.Close()

	err = pairsT.Execute(pairsFile, struct {
		Data Pairs
	}{
		pairs,
	})
	if err != nil {
		log.Fatal(err)
	}
}
