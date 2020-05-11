package main

import (
	"log"
	"sort"
	"strings"
)

// Data is the biggest cqontainer for all definitions.
type Data struct {
	Pairs    map[string]*Pair
	Infos    []*Info
	Services []*Service

	// Store all specs for encoding
	pairSpec    *PairSpec
	infoSpec    *InfoSpec
	serviceSpec []*ServiceSpec
}

// Handle will do post actions before generate
func (d *Data) Handle() {
	for _, v := range d.Pairs {
		v.Handle()
	}
	for _, v := range d.Services {
		v.Handle()
	}
}

// Sort will sort the data
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
	Infos   []*Info
}

// Handle will do post actions before generate
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

// Sort will sort the data
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

	Global               bool
	GeneratedDescription string // Description that generated from description
}

// Handle will do post actions before generate
func (p *Pair) Handle() {
	p.GeneratedDescription = strings.ReplaceAll(p.Description, "\n", "\n//")
}

// Info is the metadata definition.
type Info struct {
	Scope       string `hcl:",label"`
	Category    string `hcl:",label"`
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	DisplayName string `hcl:"display_name,optional"`
	ZeroValue   string `hcl:"zero_value,optional"`

	Global bool
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

// InfoSpec is the data parsed from HCL.
type InfoSpec struct {
	Infos []*Info `hcl:"info,block"`
}

// PairSpec is the data parsed from HCL.
type PairSpec struct {
	Pairs []*Pair `hcl:"pair,block"`
}

// NamespaceSpec is the data parsed from HCL.
type NamespaceSpec struct {
	Op []*Op `hcl:"op,block"`
}

// ServiceSpec is the data parsed from HCL.
type ServiceSpec struct {
	Name    string         `hcl:"name"`
	Service *NamespaceSpec `hcl:"service,block"`
	Storage *NamespaceSpec `hcl:"storage,block"`
	Pairs   *PairSpec      `hcl:"pairs,block"`
	Infos   *InfoSpec      `hcl:"infos,block"`
}

// FormatPairs will format pairs for pair spec
func (d *Data) FormatPairs(p *PairSpec, global bool) map[string]*Pair {
	if p == nil {
		return nil
	}

	m := make(map[string]*Pair)
	for _, v := range p.Pairs {
		v := v

		v.Global = global
		m[v.Name] = v
	}
	return m
}

// FormatInfos will format metas for meta spec
func (d *Data) FormatInfos(m *InfoSpec, global bool) []*Info {
	if m == nil {
		return nil
	}

	is := make([]*Info, 0, len(m.Infos))
	for _, v := range m.Infos {
		v := v

		v.Global = global
		is = append(is, v)
	}

	return is
}

// FormatService will format services from service spec
func (d *Data) FormatService(s *ServiceSpec) *Service {
	srv := &Service{
		Name:    s.Name,
		Storage: s.Storage.Op,
		Pairs:   mergePairs(d.Pairs, d.FormatPairs(s.Pairs, false)),
		Infos:   mergeInfos(d.Infos, d.FormatInfos(s.Infos, false)),
	}
	if s.Service != nil {
		srv.Service = s.Service.Op
	}
	return srv
}

// FormatData will format the whole data.
func FormatData(p *PairSpec, m *InfoSpec, s []*ServiceSpec) *Data {
	data := &Data{
		pairSpec:    p,
		infoSpec:    m,
		serviceSpec: s,
	}
	data.Pairs = data.FormatPairs(p, true)
	data.Infos = data.FormatInfos(m, true)

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

func mergeInfos(a, b []*Info) []*Info {
	fn := func(ms ...[]*Info) []*Info {
		ans := make([]*Info, 0)
		for _, m := range ms {
			for _, v := range m {
				v := v
				ans = append(ans, v)
			}
		}
		return ans
	}

	return fn(a, b)
}
