package main

import (
	"sort"
)

type Data struct {
	Pairs   []*Pair `hcl:"pair,block"`
	TypeMap map[string]string

	ObjectMeta       []*Metadata `hcl:"object_meta,block"`
	StorageMeta      []*Metadata `hcl:"storage_meta,block"`
	StorageStatistic []*Metadata `hcl:"storage_statistic,block"`

	Service []*Service
}

func (o *Data) Sort() {
	sort.Slice(o.Pairs, func(i, j int) bool {
		return o.Pairs[i].Name < o.Pairs[j].Name
	})

	sort.Slice(o.ObjectMeta, func(i, j int) bool {
		return o.ObjectMeta[i].Name < o.ObjectMeta[j].Name
	})
	sort.Slice(o.StorageMeta, func(i, j int) bool {
		return o.StorageMeta[i].Name < o.StorageMeta[j].Name
	})
	sort.Slice(o.StorageStatistic, func(i, j int) bool {
		return o.StorageStatistic[i].Name < o.StorageStatistic[j].Name
	})

	for _, v := range o.Service {
		v.Sort()
	}
	sort.Slice(o.Service, func(i, j int) bool {
		return o.Service[i].Name < o.Service[j].Name
	})
}

func (o *Data) ExportPairs() interface{} {
	return struct {
		Pairs []*Pair `hcl:"pair,block"`
	}{o.Pairs}
}

func (o *Data) ExportMetadata() interface{} {
	return struct {
		ObjectMeta       []*Metadata `hcl:"object_meta,block"`
		StorageMeta      []*Metadata `hcl:"storage_meta,block"`
		StorageStatistic []*Metadata `hcl:"storage_statistic,block"`
	}{
		o.ObjectMeta,
		o.StorageMeta,
		o.StorageStatistic,
	}
}

type Pair struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	Description string `hcl:"description,optional"`
	Parser      string `hcl:"parser,optional"`

	GeneratedDescription string
}

type Metadata struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	DisplayName string `hcl:"display_name,optional"`
	ZeroValue   string `hcl:"zero_value,optional"`
}

type Service struct {
	Name    string `hcl:"name"`
	Service *Ops   `hcl:"service,block"`
	Storage *Ops   `hcl:"storage,block"`

	TypeMap map[string]string
}

func (o *Service) Sort() {
	if o.Service != nil {
		o.Service.Sort()
	}
	if o.Storage != nil {
		o.Storage.Sort()
	}
}

type Ops struct {
	Op []*Op `hcl:"op,block"`
}

func (o *Ops) Sort() {
	for _, v := range o.Op {
		v.Sort()
	}
	sort.Slice(o.Op, func(i, j int) bool {
		return o.Op[i].Op < o.Op[j].Op
	})
}

type Op struct {
	Op       string   `hcl:",label"`
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`

	Generated []string // This op's generated pairs, they will be treated as optional.
	Func      *Func    // Function related to this op
}

func (o *Op) Sort() {
	sort.Strings(o.Optional)
	sort.Strings(o.Required)
}

type Func struct {
	Parent     string // Old method name: "AbortSegment"
	Receiver   string // Receiver's name: "s *Storage"
	Params     string // Method's Params: "ctx context.Context, id string, pairs ...*types.Pair"
	Returns    string // Method's returns: "err error"
	Caller     string // How to call Parent method: "id, pairs..."
	HasContext bool

	hasPair bool
}
