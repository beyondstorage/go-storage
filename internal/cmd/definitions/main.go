package main

import (
	def "github.com/beyondstorage/go-storage/v5/definitions"
)

//go:generate go run .
func main() {
	def.GenerateIterator("../../../types/iterator.generated.go")
	def.GenerateInfo("../../../types/info.generated.go")
	def.GeneratePair("../../../pairs/generated.go")
	def.GenerateOperation("../../../types/operation.generated.go")
	def.GenerateObject("../../../types/object.generated.go")
	def.GenerateNamespace("../../../definitions/namespace.generated.go")
}
