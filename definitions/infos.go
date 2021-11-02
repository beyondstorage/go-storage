package definitions

import (
	"github.com/Xuanwo/templateutils"
	"sort"
)

type Info struct {
	Namespace   string
	Category    string
	Name        string
	Type        Type
	Description string

	// Only infos that declared inside definitions can set global or export as true.
	export bool
	global bool
}

func (i Info) Global() bool {
	return i.global
}

func (i Info) NameForStructField() string {
	if i.export {
		return templateutils.ToPascal(i.Name)
	} else {
		return templateutils.ToCamel(i.Name)
	}
}

func (i Info) NameForFunctionName() string {
	return templateutils.ToPascal(i.Name)
}

func SortInfos(is []Info) []Info {
	sort.Slice(is, func(i, j int) bool {
		x, y := is[i], is[j]

		if x.Namespace != y.Namespace {
			return x.Namespace < y.Namespace
		}
		if x.Category != y.Category {
			return x.Category < y.Category
		}
		return x.Name < y.Name
	})
	return is
}

var InfosObjectMetaArray = []Info{
	{
		Name:        "append_offset",
		Type:        Type{Name: "int64"},
		Description: "AppendOffset is the offset of the append object.",
	},
	{
		Name: "content_disposition",
		Type: Type{Name: "string"},
	},
	{
		Name: "content_length",
		Type: Type{Name: "int64"},
	},
	{
		Name: "content_md5",
		Type: Type{Name: "string"},
	},
	{
		Name: "content_type",
		Type: Type{Name: "string"},
	},
	{
		Name: "etag",
		Type: Type{Name: "string"},
	},
	{
		Name:   "id",
		Type:   Type{Name: "string"},
		export: true,
		Description: `ID is the unique key in storage.

ID SHOULD be an absolute path compatible with the target operating system-defined file paths, or a unique identifier specified by service.`,
	},
	{
		Name: "last_modified",
		Type: Type{Package: "time", Name: "Time"},
	},
	{
		Name:        "link_target",
		Type:        Type{Name: "string"},
		Description: "LinkTarget is the symlink target for link object.",
	},
	{
		Name:   "mode",
		Type:   Type{Package: "types", Name: "ObjectMode"},
		export: true,
	},
	{
		Name:        "multipart_id",
		Type:        Type{Name: "string"},
		Description: "MultipartID is the part id of part object.",
	},
	{
		Name:   "path",
		Type:   Type{Name: "string"},
		export: true,
		Description: `Path is either the absolute path or the relative path towards storage's WorkDir depends on user's input.

Path SHOULD be Unix style.`,
	},
	{
		Name:        "system_metadata",
		Type:        Type{Name: "interface{}"},
		Description: "SystemMetadata stores system defined metadata.",
	},
	{
		Name:        "user_metadata",
		Type:        Type{Name: "map[string]string"},
		Description: "UserMetadata stores user defined metadata.",
	},
}

var InfosStorageMetaArray = []Info{
	{
		Name:        "append_number_maximum",
		Type:        Type{Name: "int"},
		Description: "Max append numbers in append operation.",
	},
	{
		Name:        "append_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Max append size in append operation.",
	},
	{
		Name:        "append_total_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Max append total size in append operation.",
	},
	{
		Name:        "copy_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Maximum size for copy operation.",
	},
	{
		Name:        "fetch_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Maximum size for fetch operation.",
	},
	{
		Name: "location",
		Type: Type{Name: "string"},
	},
	{
		Name:        "move_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Maximum size for move operation.",
	},
	{
		Name:        "multipart_number_maximum",
		Type:        Type{Name: "int"},
		Description: "Maximum part number in multipart operation.",
	},
	{
		Name:        "multipart_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Maximum part size in multipart operation.",
	},
	{
		Name:        "multipart_size_minimum",
		Type:        Type{Name: "int64"},
		Description: "Minimum part size in multipart operation.",
	},
	{
		Name:   "name",
		Type:   Type{Name: "string"},
		export: true,
	},
	{
		Name:        "system_metadata",
		Type:        Type{Name: "interface{}"},
		Description: "SystemMetadata stores system defined metadata.",
	},
	{
		Name:   "work_dir",
		Type:   Type{Name: "string"},
		export: true,
	},
	{
		Name:        "write_size_maximum",
		Type:        Type{Name: "int64"},
		Description: "Maximum size for write operation.",
	},
}

func init() {
	for k := range InfosObjectMetaArray {
		InfosObjectMetaArray[k].Namespace = NamespaceObject
		InfosObjectMetaArray[k].Category = CategoryMeta
		InfosObjectMetaArray[k].global = true
	}
	for k := range InfosStorageMetaArray {
		InfosObjectMetaArray[k].Namespace = NamespaceStorage
		InfosObjectMetaArray[k].Category = CategoryMeta
		InfosObjectMetaArray[k].global = true
	}
}
