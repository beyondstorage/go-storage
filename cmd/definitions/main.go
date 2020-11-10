// +build tools

package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	switch v := len(args); v {
	case 1:
		actionGlobal()
	case 2:
		actionService(args[1])
	default:
		log.Fatalf("args length should be 1 or 2, actual %d", v)
	}
}

func actionGlobal() {
	data := parse()
	data.Sort()

	generateGlobal(data)
	formatGlobal(data)
}

func actionService(filePath string) {
	data := parse()

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	srv := &ServiceSpec{}
	err = parseHCL(content, filePath, srv)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	data.Service = data.FormatService(srv)

	data.Sort()

	generateService(data)
	formatService(data)
	log.Printf("%s generate finished", filePath)
}
