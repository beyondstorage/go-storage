//go:build tools
// +build tools

package specs

const (
	featurePath     = "definitions/features.toml"
	fieldPath       = "definitions/fields.toml"
	infoObjectMeta  = "definitions/info_object_meta.toml"
	infoStorageMeta = "definitions/info_storage_meta.toml"
	operationPath   = "definitions/operations.toml"
	pairPath        = "definitions/pairs.toml"
)

var (
	ParsedFeatures   Features
	ParsedPairs      Pairs
	ParsedInfos      Infos
	ParsedOperations Operations
)

func ParseService(filePath string) (Service, error) {
	srv := parseService(filePath)

	// Make sure services has been sorted, so that we can format the specs correctly.
	srv.Sort()

	return srv, nil
}

func init() {
	ParsedFeatures = parseFeatures()
	ParsedFeatures.Sort()

	ParsedPairs = parsePairs()
	ParsedPairs.Sort()

	ParsedInfos = parseInfos()
	ParsedInfos.Sort()

	ParsedOperations = parseOperations()
	ParsedOperations.Sort()
}
