package main

import (
	"log"
	"sort"
	"strings"
)

// Data is the biggest cqontainer for all definitions.
type Data struct {
	Pairs    map[string]*Pair
	Metas    *Metas
	Services []*Service

	// Store all specs for encoding
	pairSpec    *PairSpec
	metaSpec    *MetaSpec
	serviceSpec []*ServiceSpec
}

func (d *Data) Handle() {
	for _, v := range d.Pairs {
		v.Handle()
	}
	for _, v := range d.Services {
		v.Handle()
	}
}

func (d *Data) Sort() {
	for _, v := range d.Services {
		v.Sort()
	}
}

// Service is the service definition.
type Service struct {
	Name    string
	Service []*Op
	Storage []*Op
	Pairs   map[string]*Pair
	Metas   *Metas
}

func (srv *Service) Handle() {
	servicerFuncs, storagerFuncs := parseFunc(srv.Name, "servicer"), parseFunc(srv.Name, "storager")

	// Register funcs into service
	if srv.Service != nil {
		for _, v := range srv.Service {
			if fn, ok := servicerFuncs[v.Name]; ok {
				fn.hasPair = true
				v.Func = fn
			}
		}

		// Add missing pairs into service
		for k, v := range servicerFuncs {
			v := v
			if v.hasPair {
				continue
			}
			srv.Service = append(srv.Service, &Op{
				Name: k,
				Func: v,
			})
		}
	}

	for _, v := range srv.Storage {
		if fn, ok := storagerFuncs[v.Name]; ok {
			fn.hasPair = true
			v.Func = fn
		}
	}
	for k, v := range storagerFuncs {
		v := v
		if v.hasPair {
			continue
		}
		srv.Storage = append(srv.Storage, &Op{
			Name: k,
			Func: v,
		})
	}

	// Inject pairs
	injectReadCallbackFunc(srv.Storage)
	injectContext(srv.Service)
	injectContext(srv.Storage)
	return
}

func (srv *Service) Sort() {
	sort.Slice(srv.Service, func(i, j int) bool {
		return srv.Service[i].Name < srv.Service[j].Name
	})
	sort.Slice(srv.Storage, func(i, j int) bool {
		return srv.Storage[i].Name < srv.Storage[j].Name
	})

	for _, v := range srv.Service {
		v.Sort()
	}
	for _, v := range srv.Storage {
		v.Sort()
	}
}

// Pair is the pair definition.
type Pair struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	Description string `hcl:"description,optional"`
	Parser      string `hcl:"parser,optional"`

	GeneratedDescription string // Description that generated from description
}

func (p *Pair) Handle() {
	p.GeneratedDescription = strings.ReplaceAll(p.Description, "\n", "\n//")
}

type Metas struct {
	ObjectMeta       map[string]*Metadata
	StorageMeta      map[string]*Metadata
	StorageStatistic map[string]*Metadata
}

// Metas is the metadata definition.
type Metadata struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	DisplayName string `hcl:"display_name,optional"`
	ZeroValue   string `hcl:"zero_value,optional"`
}

// Op means an operation definition.
type Op struct {
	Name     string   `hcl:",label"`
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`

	Generated []string // This op's generated pairs, they will be treated as optional.
	Func      *Func    // Function related to this op
}

// Sort will sort the operation
func (o *Op) Sort() {
	sort.Strings(o.Optional)
	sort.Strings(o.Required)
}

// Func is the function related the op.
type Func struct {
	Parent     string // Old method name: "AbortSegment"
	Receiver   string // Receiver's name: "s *Storage"
	Params     string // Method's Params: "ctx context.Context, id string, pairs ...*types.Pair"
	Returns    string // Method's returns: "err error"
	Caller     string // How to call Parent method: "id, pairs..."
	HasContext bool

	hasPair bool
}

type MetaSpec struct {
	ObjectMeta       []*Metadata `hcl:"object_meta,block"`
	StorageMeta      []*Metadata `hcl:"storage_meta,block"`
	StorageStatistic []*Metadata `hcl:"storage_statistic,block"`
}

type PairSpec struct {
	Pairs []*Pair `hcl:"pair,block"`
}

type NamespaceSpec struct {
	Op []*Op `hcl:"op,block"`
}

type ServiceSpec struct {
	Name    string         `hcl:"name"`
	Service *NamespaceSpec `hcl:"service,block"`
	Storage *NamespaceSpec `hcl:"storage,block"`
	Pairs   *PairSpec      `hcl:"pairs,block"`
	Metas   *MetaSpec      `hcl:"metas,block"`
}

func (d *Data) FormatPairs(p *PairSpec) map[string]*Pair {
	if p == nil {
		return nil
	}

	m := make(map[string]*Pair)
	for _, v := range p.Pairs {
		v := v
		m[v.Name] = v
	}
	return m
}

func (d *Data) FormatMeta(m *MetaSpec) *Metas {
	if m == nil {
		return &Metas{}
	}

	meta := &Metas{}
	meta.ObjectMeta = make(map[string]*Metadata)
	for _, v := range m.ObjectMeta {
		v := v
		meta.ObjectMeta[v.Name] = v
	}

	meta.StorageMeta = make(map[string]*Metadata)
	for _, v := range m.StorageMeta {
		v := v
		meta.StorageMeta[v.Name] = v
	}

	meta.StorageStatistic = make(map[string]*Metadata)
	for _, v := range m.StorageStatistic {
		v := v
		meta.StorageStatistic[v.Name] = v
	}

	return meta
}

func (d *Data) FormatService(s *ServiceSpec) *Service {
	srv := &Service{
		Name:    s.Name,
		Storage: s.Storage.Op,
		Pairs:   mergePairs(d.Pairs, d.FormatPairs(s.Pairs)),
		Metas:   mergeMetas(d.Metas, d.FormatMeta(s.Metas)),
	}
	if s.Service != nil {
		srv.Service = s.Service.Op
	}
	return srv
}

func FormatData(p *PairSpec, m *MetaSpec, s []*ServiceSpec) *Data {
	data := &Data{
		pairSpec:    p,
		metaSpec:    m,
		serviceSpec: s,
	}
	data.Pairs = data.FormatPairs(p)
	data.Metas = data.FormatMeta(m)

	for _, v := range s {
		data.Services = append(data.Services, data.FormatService(v))
	}
	return data
}

func mergePairs(ms ...map[string]*Pair) map[string]*Pair {
	ans := make(map[string]*Pair)
	for _, m := range ms {
		for k, v := range m {
			if _, ok := ans[k]; ok {
				log.Fatalf("pair conflict: %s", k)
			}
			v := v
			ans[k] = v
		}
	}
	return ans
}

func mergeMetas(a, b *Metas) *Metas {
	fn := func(ms ...map[string]*Metadata) map[string]*Metadata {
		ans := make(map[string]*Metadata)
		for _, m := range ms {
			for k, v := range m {
				if _, ok := ans[k]; ok {
					log.Fatalf("metadata conflict: %s", k)
				}
				v := v
				ans[k] = v
			}
		}
		return ans
	}

	m := &Metas{
		ObjectMeta:       fn(a.ObjectMeta, b.ObjectMeta),
		StorageMeta:      fn(a.StorageMeta, b.StorageMeta),
		StorageStatistic: fn(a.StorageStatistic, b.StorageStatistic),
	}
	return m
}
