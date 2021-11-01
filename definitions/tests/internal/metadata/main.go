package main

import def "go.beyondstorage.io/v5/definitions"

//go:generate go run .
func main() {
	def.GenerateService(metadata, "generated_test.go")
}

var metadata = def.Metadata{
	Name: "tests",
	Pairs: []def.Pair{
		pairDisableUriCleaning,
		pairStorageClass,
		pairStringPair,
	},
	Infos: []def.Info{
		infoObjectMetaStorageClass,
		infoStorageMetaQueriesPerSecond,
	},
	Namespaces: map[string]def.Namespace{
		def.NamespaceService: {
			Required: []def.Pair{
				def.PairCredential,
			},
			Optional: []def.Pair{
				def.PairEndpoint,
			},
			Features: []def.Feature{},
			Functions: map[string][]def.Pair{
				"create": {
					def.PairLocation,
				},
				"delete": {
					def.PairLocation,
				},
				"get": {
					def.PairLocation,
				},
				"list": {
					def.PairLocation,
				},
			},
		},
		def.NamespaceStorage: {
			Required: []def.Pair{
				def.PairName,
			},
			Optional: []def.Pair{
				def.PairLocation,
				def.PairWorkDir,
				pairDisableUriCleaning,
			},
			Features: []def.Feature{
				def.FeatureVirtualDir,
				def.FeatureLoosePair,
				def.FeatureWriteEmptyObject,
			},
			Functions: map[string][]def.Pair{
				"delete": {
					def.PairMultipartID,
					def.PairObjectMode,
				},
				"list": {
					def.PairListMode,
				},
				"read": {
					def.PairIoCallback,
					def.PairSize,
					def.PairOffset,
				},
				"create": {
					def.PairObjectMode,
				},
				"stat": {
					def.PairObjectMode,
				},
				"write": {
					def.PairContentMD5,
					def.PairContentType,
					def.PairIoCallback,
					pairStorageClass,
				},
			},
		},
	},
}

var pairDisableUriCleaning = def.Pair{
	Name: "disable_uri_cleaning",
	Type: def.Type{Name: "bool"},
}
var pairStorageClass = def.Pair{
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Defaultable: true,
}
var pairStringPair = def.Pair{
	Name:        "string_pair",
	Type:        def.Type{Name: "string"},
	Description: "tests connection string",
}

var infoObjectMetaStorageClass = def.Info{
	Namespace:   def.NamespaceObject,
	Category:    def.CategoryMeta,
	Name:        "storage_class",
	Type:        def.Type{Name: "string"},
	Description: "is the storage class for this object",
}
var infoStorageMetaQueriesPerSecond = def.Info{
	Namespace:   def.NamespaceStorage,
	Category:    def.CategoryMeta,
	Name:        "queries_per_second",
	Type:        def.Type{Name: "int64"},
	Description: "tests storage system metadata",
}

func init() {
	dp := make([]def.Pair, 0)
	for _, v := range metadata.Pairs {
		if !v.Defaultable {
			continue
		}
		dp = append(dp, def.Pair{
			Name:        "default_" + v.Name,
			Type:        v.Type,
			Description: "default value for " + v.Name,
		})
	}
	metadata.Pairs = append(metadata.Pairs, dp...)
}
