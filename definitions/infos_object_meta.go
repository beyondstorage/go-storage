package definitions

var InfosObjectMetaArray = []Info{
	InfoObjectMetaAppendOffset,
	InfoObjectMetaContentDisposition,
	InfoObjectMetaContentLength,
	InfoObjectMetaContentMd5,
	InfoObjectMetaContentType,
	InfoObjectMetaEtag,
	InfoObjectMetaId,
	InfoObjectMetaLastModified,
	InfoObjectMetaLinkTarget,
	InfoObjectMetaMode,
	InfoObjectMetaMultipartID,
	InfoObjectMetaPath,
	InfoObjectMetaSystemMetadata,
	InfoObjectMetaUserMetadata,
}

var InfoObjectMetaAppendOffset = Info{
	Namespace:   NamespaceObject,
	Category:    CategoryMeta,
	Name:        "append_offset",
	Type:        Type{Name: "int64"},
	Description: "AppendOffset is the offset of the append object.",
	global:      true,
}

var InfoObjectMetaContentDisposition = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "content_disposition",
	Type:      Type{Name: "string"},
	global:    true,
}

var InfoObjectMetaContentLength = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "content_length",
	Type:      Type{Name: "int64"},
	global:    true,
}

var InfoObjectMetaContentMd5 = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "content_md5",
	Type:      Type{Name: "string"},
	global:    true,
}

var InfoObjectMetaContentType = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "content_type",
	Type:      Type{Name: "string"},
	global:    true,
}

var InfoObjectMetaEtag = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "etag",
	Type:      Type{Name: "string"},
	global:    true,
}

var InfoObjectMetaId = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "id",
	Type:      Type{Name: "string"},
	Export:    true,
	global:    true,
	Description: `ID is the unique key in storage.

ID SHOULD be an absolute path compatible with the target operating system-defined file paths, or a unique identifier specified by service.`,
}

var InfoObjectMetaLastModified = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "last_modified",
	Type:      Type{Package: "time", Name: "Time"},
	global:    true,
}

var InfoObjectMetaLinkTarget = Info{
	Namespace:   NamespaceObject,
	Category:    CategoryMeta,
	Name:        "link_target",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "LinkTarget is the symlink target for link object.",
}
var InfoObjectMetaMode = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "mode",
	Type:      Type{Package: "types", Name: "ObjectMode"},
	Export:    true,
	global:    true,
}
var InfoObjectMetaMultipartID = Info{
	Namespace:   NamespaceObject,
	Category:    CategoryMeta,
	Name:        "multipart_id",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "MultipartID is the part id of part object.",
}
var InfoObjectMetaPath = Info{
	Namespace: NamespaceObject,
	Category:  CategoryMeta,
	Name:      "path",
	Type:      Type{Name: "string"},
	Export:    true,
	global:    true,
	Description: `Path is either the absolute path or the relative path towards storage's WorkDir depends on user's input.

Path SHOULD be Unix style.`,
}

var InfoObjectMetaSystemMetadata = Info{
	Namespace:   NamespaceObject,
	Category:    CategoryMeta,
	Name:        "system_metadata",
	Type:        Type{Name: "interface{}"},
	global:      true,
	Description: "SystemMetadata stores system defined metadata.",
}

var InfoObjectMetaUserMetadata = Info{
	Namespace:   NamespaceObject,
	Category:    CategoryMeta,
	Name:        "user_metadata",
	Type:        Type{Name: "map[string]string"},
	global:      true,
	Description: "UserMetadata stores user defined metadata.",
}
