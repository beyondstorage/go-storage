package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	switch v := len(os.Args); v {
	case 1:
		actionGlobal()
	case 2:
		actionService()
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

func actionService() {
	data := parse()

	filePath := os.Args[1]
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
