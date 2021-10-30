package definitions

import "sort"

type Pair struct {
	Name        string
	Type        Type
	Defaultable bool
	Description string

	// Only infos that declared inside definitions can set global as true.
	global bool
}

func (p Pair) Global() bool {
	return p.global
}

var PairMap = make(map[string]Pair)

var PairArray = []Pair{
	PairContentDisposition,
	PairContentMD5,
	PairContentType,
}

func init() {
	for _, v := range PairArray {
		PairMap[v.Name] = v
	}
}

func SortPairs(ps []Pair) []Pair {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Name < ps[j].Name
	})
	return ps
}

var PairContentDisposition = Pair{
	Name: "content_disposition",
	Type: Type{
		Name: "string",
	},
	global: true,
}

var PairContentMD5 = Pair{
	Name: "content_md5",
	Type: Type{
		Name: "string",
	},
	global: true,
}

var PairContentType = Pair{
	Name: "content_type",
	Type: Type{
		Name: "string",
	},
	Defaultable: true,
	global:      true,
}

var PairContinuationToken = Pair{
	Name:        "continuation_token",
	Type:        Type{Name: "string"},
	Description: "specify the continuation token for lis",
	global:      true,
}

var PairCredential = Pair{
	Name: "credential",
	Type: Type{
		Name: "string",
	},
	Description: "specify how to provide credential for service or storage",
	global:      true,
}

var PairEndpoint = Pair{
	Name: "endpoint",
	Type: Type{
		Name: "string",
	},
	Description: "specify how to provide endpoint for service or storage",
	global:      true,
}

var PairListMode = Pair{
	Name: "list_mode",
	Type: Type{
		Package: "types",
		Name:    "ListMode",
	},
	global: true,
}

var PairLocation = Pair{
	Name:        "location",
	Type:        Type{Name: "string"},
	Description: "specify the location for service or storage",
	global:      true,
}

var PairName = Pair{
	Name:        "name",
	Type:        Type{Name: "string"},
	Description: "specify the storage name",
	global:      true,
}
var PairObjectMode = Pair{
	Name: "object_mode",
	Type: Type{
		Package: "types",
		Name:    "ObjectMode",
	},
	global: true,
}
var PairOffset = Pair{
	Name: "offset",
	Type: Type{
		Name: "int64",
	},
	global: true,
}
var PairSize = Pair{
	Name: "size",
	Type: Type{
		Name: "int64",
	},
	global: true,
}
var PairMultipartId = Pair{
	Name: "multipart_id",
	Type: Type{
		Name: "string",
	},
	global: true,
}
var PairIoCallback = Pair{
	Name: "io_callback",
	Type: Type{
		Name: "func([]byte)",
	},
	global: true,
}
