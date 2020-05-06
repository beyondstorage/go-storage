package main

type Data struct {
	Pairs   []*Pair
	TypeMap map[string]string

	ObjectMeta       []*Metadata
	StorageMeta      []*Metadata
	StorageStatistic []*Metadata

	Service []*Service
}

type Pair struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	Description string `hcl:"description,optional"`
	Parser      string `hcl:"parser,optional"`
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

type Ops struct {
	Op []*Op `hcl:"op,block"`
}

type Op struct {
	Op       string   `hcl:",label"`
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`
	Func     *Func    `hcl:"-"`
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
