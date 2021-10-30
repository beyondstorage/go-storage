package definitions

import "sort"

type Feature struct {
	Name        string
	Namespaces  []string
	Description string
}

func (f Feature) HasNamespace(ns string) bool {
	for _, n := range f.Namespaces {
		if n == ns {
			return true
		}
	}
	return false
}

var FeaturesMap = make(map[string]Feature)

var FeaturesArray = []Feature{
	FeatureLoosePair,
	FeatureVirtualDir,
	FeatureVirtualObjectMetadata,
	FeatureVirtualLink,
	FeatureWriteEmptyObject,
}

func init() {
	for _, v := range FeaturesArray {
		FeaturesMap[v.Name] = v
	}
}

func SortFeatures(fe []Feature) []Feature {
	sort.Slice(fe, func(i, j int) bool {
		return fe[i].Name < fe[j].Name
	})
	return fe
}

var FeatureLoosePair = Feature{
	Name:       "loose_pair",
	Namespaces: []string{NamespaceService, NamespaceStorage},
	Description: `loose_pair feature is designed for users who don't want strict pair checks.

If this feature is enabled, the service will not return an error for not support pairs.

This feature was introduced in GSP-109.`,
}

var FeatureVirtualDir = Feature{
	Name:       "virtual_dir",
	Namespaces: []string{NamespaceStorage},
	Description: `virtual_dir feature is designed for a service that doesn't have native dir support but wants to provide simulated operations.

- If this feature is disabled (the default behavior), the service will behave like it doesn't have any dir support.
- If this feature is enabled, the service will support simulated dir behavior in create_dir, create, list, delete, and so on.

This feature was introduced in GSP-109.`,
}

var FeatureVirtualObjectMetadata = Feature{
	Name:       "virtual_object_metadata",
	Namespaces: []string{NamespaceStorage},
	Description: `virtual_object_metadata feature is designed for a service that doesn't have native object metadata support but wants to provide simulated operations.

- If this feature is disabled (the default behavior), the service will behave like it doesn't have any object metadata support.
- If this feature is enabled, the service will support simulated object metadata behavior in the list, stat, and so on.

This feature was introduced in GSP-109.`,
}

var FeatureVirtualLink = Feature{
	Name:       "virtual_link",
	Namespaces: []string{NamespaceStorage},
	Description: `virtual_link feature is designed for a service that doesn't have native support for link.

- If this feature is enabled, the service will run compatible mode: create link via native methods, but allow read link from old-style link object.
- If this feature is not enabled, the service will run in native as other service.

This feature was introduced in GSP-86.`,
}

var FeatureWriteEmptyObject = Feature{
	Name:       "write_empty_object",
	Namespaces: []string{NamespaceStorage},
	Description: `write_empty_object feature is designed for a service that support write empty object.

This behavior was defined in GSP-751 and classified as an operation-related feature in GSP-837.`,
}
