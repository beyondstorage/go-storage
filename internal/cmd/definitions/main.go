package main

import (
	def "go.beyondstorage.io/v5/definitions"
)

//go:generate go run .
func main() {
	// Iterator generate
	def.GenerateIterator("../../../types/iterator.generated.go")

	// Metas generate
	def.GenerateInfo("../../../types/info.generated.go")

	// Pair generate
	def.GeneratePair("../../../pairs/generated.go")

	// Operation generate
	def.GenerateOperation("../../../types/operation.generated.go")

	// Object generate
	def.GenerateObject("../../../types/object.generated.go")

	// Feature generate
	def.GenerateFeatures("../../../types/feature.generated.go")
}
