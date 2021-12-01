package main

import (
	def "go.beyondstorage.io/v5/definitions"
	// "go.beyondstorage.io/v5/types"
)

var Metadata = def.Metadata{
	Name:  "fs",
	Pairs: []def.Pair{},
	Infos: []def.Info{},
	Factory: []def.Pair{
		def.PairWorkDir,
	},
	Service: def.Service{},
	Storage: def.Storage{
		// Features: types.StorageFeatures{
		// 	WriteEmptyObject: true,

		// 	Copy:         true,
		// 	Create:       true,
		// 	CreateAppend: true,
		// 	CreateDir:    true,
		// 	CreateLink:   true,
		// 	Delete:       true,
		// 	Fetch:        true,
		// 	List:         true,
		// 	Metadata:     true,
		// 	Move:         true,
		// 	Read:         true,
		// 	Stat:         true,
		// 	Write:        true,
		// 	WriteAppend:  true,
		// 	CommitAppend: true,
		// },

		// Create: []def.Pair{
		// 	def.PairObjectMode,
		// },
		// Delete: []def.Pair{
		// 	def.PairObjectMode,
		// },
		// List: []def.Pair{
		// 	def.PairListMode,
		// 	def.PairContinuationToken,
		// },
		// Read: []def.Pair{
		// 	def.PairOffset,
		// 	def.PairIoCallback,
		// 	def.PairSize,
		// },
		// Write: []def.Pair{
		// 	def.PairContentMD5,
		// 	def.PairContentType,
		// 	def.PairOffset,
		// 	def.PairIoCallback,
		// },
		// Stat: []def.Pair{
		// 	def.PairObjectMode,
		// },
	},
}
