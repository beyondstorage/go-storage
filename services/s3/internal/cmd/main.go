package main

import (
	"go.beyondstorage.io/services/s3/v3/internal/meta"
	def "go.beyondstorage.io/v5/definitions"
)

func main() {
	def.GenerateService(meta.Metadata, "generated.go")
}
