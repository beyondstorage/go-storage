package main

import (
	def "go.beyondstorage.io/v5/definitions"
)

var Metadata = def.Metadata{
	Name:  "webdav",
	Pairs: []def.Pair{},
	Infos: []def.Info{},
	Factory: []def.Pair{
		def.PairWorkDir,
	},
	Service: def.Service{},
	Storage: def.Storage{},
}
