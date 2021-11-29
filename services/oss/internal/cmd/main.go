package main

import (
	def "go.beyondstorage.io/v5/definitions"
)

func main() {
	def.GenerateService(Metadata, "generated.go")
}
