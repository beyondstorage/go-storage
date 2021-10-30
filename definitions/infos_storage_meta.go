package definitions

var InfosStorageMetaArray = []Info{
	InfoStorageMetaAppendNumberMaximum,
	InfoStorageMetaAppendSizeMaximum,
	InfoStorageMetaAppendTotalSizeMaximum,
	InfoStorageMetaCopySizeMaximum,
	InfoStorageMetaFetchSizeMaximum,
	InfoStorageMetaLocation,
	InfoStorageMetaMoveSizeMaximum,
	InfoStorageMetaMultipartNumberMaximum,
	InfoStorageMetaMultipartSizeMaximum,
	InfoStorageMetaMultipartSizeMinimum,
	InfoStorageMetaName,
	InfoStorageMetaSystemMetadata,
	InfoStorageMetaWorkDir,
	InfoStorageMetaWriteSizeMaximum,
}

var InfoStorageMetaAppendNumberMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "append_number_maximum",
	Type:        Type{Name: "int"},
	global:      true,
	Description: "Max append numbers in append operation.",
}

var InfoStorageMetaAppendSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "append_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Max append size in append operation.",
}

var InfoStorageMetaAppendTotalSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "append_total_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Max append total size in append operation.",
}

var InfoStorageMetaCopySizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "copy_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Maximum size for copy operation.",
}

var InfoStorageMetaFetchSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "fetch_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Maximum size for fetch operation.",
}

var InfoStorageMetaLocation = Info{
	Namespace: NamespaceStorage,
	Category:  CategoryMeta,
	Name:      "location",
	global:    true,
	Type:      Type{Name: "string"},
}

var InfoStorageMetaMoveSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "move_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Maximum size for move operation.",
}

var InfoStorageMetaMultipartNumberMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "multipart_number_maximum",
	Type:        Type{Name: "int"},
	global:      true,
	Description: "Maximum part number in multipart operation.",
}

var InfoStorageMetaMultipartSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "multipart_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Maximum part size in multipart operation.",
}

var InfoStorageMetaMultipartSizeMinimum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "multipart_size_minimum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Minimum part size in multipart operation.",
}

var InfoStorageMetaName = Info{
	Namespace: NamespaceStorage,
	Category:  CategoryMeta,
	Name:      "name",
	Type:      Type{Name: "string"},
	Export:    true,
	global:    true,
}

var InfoStorageMetaSystemMetadata = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "system_metadata",
	Type:        Type{Name: "interface{}"},
	global:      true,
	Description: "SystemMetadata stores system defined metadata.",
}

var InfoStorageMetaWorkDir = Info{
	Namespace: NamespaceStorage,
	Category:  CategoryMeta,
	Name:      "work_dir",
	Type:      Type{Name: "string"},
	Export:    true,
	global:    true,
}

var InfoStorageMetaWriteSizeMaximum = Info{
	Namespace:   NamespaceStorage,
	Category:    CategoryMeta,
	Name:        "write_size_maximum",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: "Maximum size for write operation.",
}
