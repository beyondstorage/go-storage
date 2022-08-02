package main

import (
	def "github.com/beyondstorage/go-storage/v5/definitions"
)

func main() {
	def.GenerateService(Metadata, "generated.go")
}
