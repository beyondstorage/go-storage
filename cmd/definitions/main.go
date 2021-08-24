//go:build tools
// +build tools

package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/beyondstorage/go-storage/v4/cmd/definitions/specs"
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
}

func actionService(filePath string) {
	data := parse()

	srv, err := specs.ParseService(filePath)
	if err != nil {
		log.Fatalf("parse: %v", err)
	}
	data.Service = data.FormatService(srv)

	data.Sort()

	generateService(data)
	log.Printf("%s generate finished", filePath)
}
