// +build tools

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Xuanwo/templateutils"
	specs "github.com/aos-dev/specs/go"
)

func parse() (data *Data) {
	data = FormatData(&specs.ParsedPairs, &specs.ParsedInfos, &specs.ParsedOperations)
	return data
}

func parseFunc(name string) map[string]*templateutils.Method {
	data := make(map[string]*templateutils.Method)
	filename := fmt.Sprintf("%s.go", name)

	content, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		return data
	}
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}

	source := &templateutils.Source{}
	err = source.ParseContent(filename, content)
	if err != nil {
		log.Fatalf("parse content: %v", err)
	}

	for _, v := range source.Methods {
		v := v
		data[v.Name] = v
	}
	return data
}
