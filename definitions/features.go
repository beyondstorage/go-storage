package definitions

import (
	"sort"
)

const (
	FeatureTypeOperation = "operation"
	FeatureTypeSystem    = "system"
	FeatureTypeVirtual   = "virtual"
)

type Feature struct {
	Name        string
	Type        string
	Description string
}

func SortFeatures(fe []Feature) []Feature {
	sort.Slice(fe, func(i, j int) bool {
		if fe[i].Type != fe[j].Type {
			return fe[i].Type < fe[j].Type
		}
		return fe[i].Name < fe[j].Name
	})
	return fe
}

var FeaturesService = append([]Feature{
	{
		Name: "loose_pair",
		Type: FeatureTypeVirtual,
		Description: `loose_pair feature is designed for users who don't want strict pair checks.

If this feature is enabled, the service will not return an error for not support pairs.

This feature was introduced in GSP-109.`,
	},
}, buildOperationFeatures(OperationsService)...)

var FeaturesStorage = append([]Feature{
	{
		Name: "loose_pair",
		Type: FeatureTypeVirtual,
		Description: `loose_pair feature is designed for users who don't want strict pair checks.

If this feature is enabled, the service will not return an error for not support pairs.

This feature was introduced in GSP-109.`,
	},
	{
		Name: "virtual_dir",
		Type: FeatureTypeVirtual,
		Description: `virtual_dir feature is designed for a service that doesn't have native dir support but wants to provide simulated operations.

- If this feature is disabled (the default behavior), the service will behave like it doesn't have any dir support.
- If this feature is enabled, the service will support simulated dir behavior in create_dir, create, list, delete, and so on.

This feature was introduced in GSP-109.`,
	},
	{
		Name: "virtual_link",
		Type: FeatureTypeVirtual,
		Description: `virtual_link feature is designed for a service that doesn't have native support for link.
- If this feature is enabled, the service will run compatible mode: create link via native methods, but allow read link from old-style link object.
- If this feature is not enabled, the service will run in native as other service.
This feature was introduced in GSP-86.`,
	},
	{
		Name: "virtual_object_metadata",
		Type: FeatureTypeVirtual,
		Description: `virtual_object_metadata feature is designed for a service that doesn't have native object metadata support but wants to provide simulated operations.

- If this feature is disabled (the default behavior), the service will behave like it doesn't have any object metadata support.
- If this feature is enabled, the service will support simulated object metadata behavior in the list, stat, and so on.

This feature was introduced in GSP-109.`,
	},
	{
		Name: "write_empty_object",
		Type: FeatureTypeSystem,
		Description: `write_empty_object feature is designed for a service that support write empty object.

This behavior was defined in GSP-751 and classified as an operation-related feature in GSP-837.`,
	},
}, buildOperationFeatures(OperationsStorage)...)

func buildOperationFeatures(ops []Operation) []Feature {
	fs := make([]Feature, 0)
	for _, op := range ops {
		fs = append(fs, Feature{
			Name: op.Name,
			Type: FeatureTypeOperation,
		})
	}
	return fs
}
