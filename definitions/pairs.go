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
	PairContinuationToken,
	PairCredential,
	PairEndpoint,
	PairListMode,
	PairLocation,
	PairName,
	PairObjectMode,
	PairOffset,
	PairSize,
	PairMultipartID,
	PairIoCallback,
	PairWorkDir,
}

func init() {
	// Normalize default pairs.
	dps := make([]Pair, 0)
	for _, v := range PairArray {
		if !v.Defaultable {
			continue
		}
		dps = append(dps, Pair{
			Name:        "default_" + v.Name,
			Type:        v.Type,
			Description: "default value for " + v.Name,
			global:      true,
		})
	}
	PairArray = append(PairArray, dps...)

	// Build feature pairs.
	fps := make([]Pair, 0)
	for _, f := range FeaturesArray {
		fps = append(fps, Pair{
			Name:        "enable_" + f.Name,
			Type:        Type{Name: "bool"},
			Description: "Enable feature " + f.Name,
			global:      true,
		})
	}
	PairArray = append(PairArray, fps...)

	// Setup maps.
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
	Name:   "content_disposition",
	Type:   Type{Name: "string"},
	global: true,
}

var PairContentMD5 = Pair{
	Name:   "content_md5",
	Type:   Type{Name: "string"},
	global: true,
}

var PairContentType = Pair{
	Name:        "content_type",
	Type:        Type{Name: "string"},
	Defaultable: true,
	global:      true,
}

var PairContinuationToken = Pair{
	Name:        "continuation_token",
	Type:        Type{Name: "string"},
	Description: "specify the continuation token for list",
	global:      true,
}

var PairCredential = Pair{
	Name:        "credential",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "specify how to provide credential for service or storage",
}

var PairEndpoint = Pair{
	Name:        "endpoint",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "specify how to provide endpoint for service or storage",
}

var PairListMode = Pair{
	Name:   "list_mode",
	Type:   Type{Package: "types", Name: "ListMode"},
	global: true,
}

var PairLocation = Pair{
	Name:        "location",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "specify the location for service or storage",
}

var PairName = Pair{
	Name:        "name",
	Type:        Type{Name: "string"},
	global:      true,
	Description: "specify the storage name",
}
var PairObjectMode = Pair{
	Name:        "object_mode",
	Type:        Type{Package: "types", Name: "ObjectMode"},
	global:      true,
	Description: `ObjectMode hint`,
}
var PairOffset = Pair{
	Name:        "offset",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: `specify offset for this request, storage will seek to this offset before read`,
}
var PairSize = Pair{
	Name:        "size",
	Type:        Type{Name: "int64"},
	global:      true,
	Description: `specify size for this request, storage will only read limited content data`,
}
var PairMultipartID = Pair{
	Name:   "multipart_id",
	Type:   Type{Name: "string"},
	global: true,
}
var PairIoCallback = Pair{
	Name:        "io_callback",
	Type:        Type{Name: "func([]byte)"},
	Defaultable: true,
	global:      true,
	Description: `specify what todo every time we read data from source`,
}
var PairWorkDir = Pair{
	Name:   "work_dir",
	Type:   Type{Name: "string"},
	global: true,
	Description: `specify the work dir for service or storage, every operation will be relative to this dir.
work_dir SHOULD be an absolute path.
work_dir will be default to / if not set.
work_dir SHOULD be Unix style for object storage services.
For fs storage service on windows platform, the behavior is defined separately.`,
}
