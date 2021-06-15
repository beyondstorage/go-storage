// +build tools

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Xuanwo/templateutils"
	specs "github.com/beyondstorage/specs/go"
	log "github.com/sirupsen/logrus"
)

// Data is the biggest container for all definitions.
type Data struct {
	Pairs      map[string]*Pair
	Infos      []*Info
	InfosMap   map[string][]*Info
	ObjectMeta []*Info
	Service    *Service

	Interfaces    []*Interface
	interfacesMap map[string]*Interface

	// Store all specs for encoding
	pairSpec       specs.Pairs
	infoSpec       specs.Infos
	operationsSpec specs.Operations
	serviceSpec    specs.Service
}

// Service is the service definition.
type Service struct {
	Name       string
	Namespaces []*Namespace
	pairs      map[string]*Pair
	Infos      []*Info
}

// Sort will sort the service
func (s *Service) Sort() {
	// Make sure namespaces sorted by name.
	sort.Slice(s.Namespaces, func(i, j int) bool {
		n := s.Namespaces

		return n[i].Name < n[j].Name
	})

	for _, v := range s.Namespaces {
		v.Sort()
	}
}

// Pairs returns a sorted pair.
func (s *Service) Pairs() []*Pair {
	keys := make([]string, 0, len(s.pairs))

	for k := range s.pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ps := make([]*Pair, 0, len(s.pairs))
	for _, v := range keys {
		ps = append(ps, s.pairs[v])
	}
	return ps
}

// Namespace contains all info about a namespace
type Namespace struct {
	Name       string
	New        *Function
	Funcs      []*Function
	Interfaces []*Interface
}

// Sort will sort the namespace
func (n *Namespace) Sort() {
	sort.Slice(n.Funcs, func(i, j int) bool {
		x := n.Funcs
		return x[i].Name < x[j].Name
	})

	n.New.Sort()
	for _, v := range n.Funcs {
		v.Sort()
	}
}

// Pair is the pair definition.
type Pair struct {
	Name string

	ptype string

	// Runtime generated
	Global      bool
	Description string

	// This is a service pair having the same name and type as a global pair
	Conflict bool
}

func (p *Pair) Type() string {
	return parseType(p.ptype)
}

// Format will formatGlobal current pair
func (p *Pair) Format(s specs.Pair, global bool) {
	p.Name = s.Name
	p.ptype = s.Type
	p.Global = global

	p.Description = formatDescription(templateutils.ToPascal(p.Name), s.Description)
}

// Info is the metadata definition.
type Info struct {
	Scope       string
	Category    string
	Name        string
	Export      bool
	Description string

	itype string

	Global bool
}

// Format will formatGlobal info spec into Info
func (i *Info) Format(s specs.Info, global bool) {
	i.Scope = s.Scope
	i.Category = s.Category
	i.Name = s.Name
	i.itype = s.Type
	i.Export = s.Export
	i.Description = formatDescription(templateutils.ToPascal(s.Name), s.Description)

	i.Global = global
}

func (i *Info) Type() string {
	return parseType(i.itype)
}

func (i *Info) TypeName() string {
	if i.Export {
		return templateutils.ToPascal(i.Name)
	} else {
		return templateutils.ToCamel(i.Name)
	}
}
func (i *Info) DisplayName() string {
	return templateutils.ToPascal(i.Name)
}

// Interface represents an interface
type Interface struct {
	Name        string
	Description string
	Ops         map[string]*Operation
}

// NewInterface will create a new interface from spec.
func NewInterface(in specs.Interface, fields map[string]*Field) *Interface {
	inter := &Interface{
		Name:        in.Name,
		Description: formatDescription(templateutils.ToPascal(in.Name), in.Description),
		Ops:         make(map[string]*Operation),
	}
	for _, v := range in.Ops {
		// Update op maps
		inter.Ops[v.Name] = NewOperation(v, fields)
	}

	return inter
}

// DisplayName will output interface's display name.
func (i *Interface) DisplayName() string {
	return templateutils.ToPascal(i.Name)
}

// Operation represents an operation.
type Operation struct {
	Name        string
	Description string
	Pairs       []string
	Params      Fields
	Results     Fields
	ObjectMode  string
	Local       bool
}

// NewOperation will create an new operation from operation spec.
func NewOperation(v specs.Operation, fields map[string]*Field) *Operation {
	op := &Operation{
		Name:        v.Name,
		Local:       v.Local,
		ObjectMode:  v.ObjectMode,
		Description: formatDescription("", v.Description),
	}
	for _, f := range v.Params {
		op.Params = append(op.Params, fields[f])
	}
	// Inject pairs
	op.Params = append(op.Params, fields["pairs"])

	for _, f := range v.Results {
		op.Results = append(op.Results, fields[f])
	}
	// As long as the function is not local, an error may be occur
	if !op.Local {
		// Inject error for non-local functions.
		op.Results = append(op.Results, fields["err"])
	}

	// Add pairs
	op.Pairs = v.Pairs

	return op
}

// FormatParams print params.
func (o *Operation) FormatParams() string {
	s := make([]string, 0)
	for _, v := range o.Params {
		s = append(s, v.String())
	}
	return strings.Join(s, ",")
}

// FormatResults will print results.
func (o *Operation) FormatResults() string {
	s := make([]string, 0)
	for _, v := range o.Results {
		s = append(s, v.String())
	}
	return strings.Join(s, ",")
}

// FormatResultsWithPackageName will print results with package name
//
// If type is starts with this package name, we will ignore it.
func (o *Operation) FormatResultsWithPackageName(packageName string) string {
	s := make([]string, 0)
	for _, v := range o.Results {
		if strings.HasPrefix(v.ftype, packageName) {
			s = append(s, v.Name+" "+strings.TrimPrefix(v.ftype, packageName+"."))
			continue
		}
		s = append(s, v.String())
	}
	return strings.Join(s, ",")
}

// Function represents a function.
type Function struct {
	*Operation

	Simulated bool // This op is simulated, user can decide whether use the virtual function or not.

	Required []*Pair // TODO: other functions could not have required pairs.
	Optional []*Pair
	Virtual  []*Pair // This op's virtual pairs, user can decide whether use the virtual pairs.

	Implemented bool // flag for whether this function has been implemented or not.
}

// NewFunction will createn a new function.
func NewFunction(o *Operation) *Function {
	return &Function{Operation: o}
}

// Format will formatGlobal a function with Op.
func (f *Function) Format(s specs.Op, p map[string]*Pair) {
	f.Simulated = s.Simulated

	for _, v := range s.Required {
		pair, ok := p[v]
		if !ok {
			log.Fatalf("pair %s is not exist", v)
		}
		f.Required = append(f.Required, pair)
	}
	for _, v := range s.Optional {
		pair, ok := p[v]
		if !ok {
			log.Fatalf("pair %s is not exist", v)
		}
		f.Optional = append(f.Optional, pair)
	}
	for _, v := range s.Virtual {
		pair, ok := p[v]
		if !ok {
			log.Fatalf("pair %s is not exist", v)
		}
		f.Virtual = append(f.Virtual, pair)
	}
}

// Sort will sort this function.
func (f *Function) Sort() {
	sort.Slice(f.Required, func(i, j int) bool {
		x := f.Required
		return x[i].Name < x[j].Name
	})
	sort.Slice(f.Optional, func(i, j int) bool {
		x := f.Optional
		return x[i].Name < x[j].Name
	})
}

// Fields is a slice for field.
type Fields []*Field

// String implements the stringer interface.
func (f Fields) String() string {
	x := make([]string, 0)
	for _, v := range f {
		x = append(x, v.String())
	}
	return strings.Join(x, ",")
}

// StringEndswithComma will print string with comma aware.
func (f Fields) StringEndswithComma() string {
	content := f.String()
	if content == "" {
		return ""
	}
	return content + ","
}

// Caller will print caller foramt.
func (f Fields) Caller() string {
	x := make([]string, 0)
	for _, v := range f {
		x = append(x, v.Caller())
	}
	return strings.Join(x, ",")
}

// HasReader will check whether we have reader here.
func (f Fields) HasReader() bool {
	for _, v := range f {
		if v.ftype == "io.Reader" || v.ftype == "io.ReadCloser" {
			return true
		}
	}
	return false
}

// CallerEndswithComma will print caller with comma aware.
func (f Fields) CallerEndswithComma() string {
	content := f.Caller()
	if content == "" {
		return ""
	}
	return content + ","
}

// TrimLast will trim the last fields.
func (f Fields) TrimLast() Fields {
	return f[:len(f)-1]
}

// PathCaller will print caller with path aware.
func (f Fields) PathCaller() string {
	x := make([]string, 0)
	for _, v := range f {
		if v.ftype != "string" {
			break
		}

		x = append(x, v.Caller())
	}

	content := strings.Join(x, ",")
	if content == "" {
		return ""
	}
	return "," + content
}

// Field represent a field.
type Field struct {
	Name  string
	ftype string
}

// String will print field in string formatGlobal.
func (f *Field) String() string {
	if f.Name == "" {
		return f.Type()
	}
	return fmt.Sprintf("%s %s", f.Name, f.Type())
}

func (f *Field) Type() string {
	return parseType(f.ftype)
}

// Caller will print the caller formatGlobal of field.
func (f *Field) Caller() string {
	if strings.HasPrefix(f.Type(), "...") {
		return f.Name + "..."
	}
	return f.Name
}

// Format will create a new field.
func (f *Field) Format(s specs.Field) {
	f.ftype = s.Type
	f.Name = s.Name
}

// FormatPairs will formatGlobal pairs for pair spec
func (d *Data) FormatPairs(p specs.Pairs, global bool) map[string]*Pair {
	m := make(map[string]*Pair)
	for _, v := range p {
		pair := &Pair{}
		pair.Format(v, global)

		m[pair.Name] = pair
	}
	return m
}

// FormatInfos will formatGlobal metas for meta spec
func (d *Data) FormatInfos(m specs.Infos, global bool) []*Info {
	is := make([]*Info, 0, len(m))
	for _, v := range m {
		i := &Info{}
		i.Format(v, global)

		is = append(is, i)
	}

	d.InfosMap = make(map[string][]*Info)
	for _, v := range is {
		v := v

		typeName := fmt.Sprintf("%s-%s", v.Scope, v.Category)
		if typeName == "object-meta" {
			d.ObjectMeta = append(d.ObjectMeta, v)
			continue
		}

		d.InfosMap[typeName] = append(d.InfosMap[typeName], v)
	}

	return is
}

// FormatOperations will formatGlobal operations from operation spec
func (d *Data) FormatOperations(o specs.Operations) (ins []*Interface, inm map[string]*Interface) {
	fileds := make(map[string]*Field)
	for _, v := range o.Fields {
		f := &Field{}
		f.Format(v)

		fileds[v.Name] = f
	}

	// Build all interfaces.
	inm = make(map[string]*Interface)
	for _, in := range o.Interfaces {
		inter := NewInterface(in, fileds)

		ins = append(ins, inter)
		inm[inter.Name] = inter
	}

	return
}

// FormatNamespace will formatGlobal a namespace.
func (d *Data) FormatNamespace(srv *Service, n specs.Namespace) *Namespace {
	ns := &Namespace{Name: n.Name}

	nsInterface := n.Name + "r"

	// Handle New function
	ns.New = NewFunction(&Operation{Name: "new"})
	ns.New.Format(specs.Op{
		Required: n.New.Required,
		Optional: n.New.Optional,
	}, srv.pairs)

	// Handle other interfaces.
	fns := make(map[string]*Function)

	// Add namespace itself into implements.
	implements := append(n.Implement[:], nsInterface)
	for _, interfaceName := range implements {
		inter := d.interfacesMap[interfaceName]

		// Add interface into namespace's interface list.
		ns.Interfaces = append(ns.Interfaces, inter)

		// Add all functions under interface into namespace's func list.
		for k, v := range inter.Ops {
			v := v

			f := NewFunction(v)
			ns.Funcs = append(ns.Funcs, f)
			fns[k] = f
		}
	}

	for _, v := range n.Op {
		fns[v.Name].Format(v, srv.pairs)
	}

	implemented := parseFunc(n.Name)
	for _, fn := range fns {
		x := templateutils.ToCamel(fn.Name)
		if _, ok := implemented[x]; ok {
			fn.Implemented = true
		}
	}

	d.ValidateNamespace(srv, ns)
	return ns
}

// ValidateNamespace will inject a namespace to insert generated pairs.
func (d *Data) ValidateNamespace(srv *Service, n *Namespace) {
	for _, v := range n.Funcs {
		// For now, we disallow required pairs for Storage.
		if n.Name == "Storage" && len(v.Required) > 0 {
			log.Fatalf("Operation [%s] cannot specify required pairs.", v.Name)
		}

		existPairs := map[string]bool{}
		for _, p := range v.Optional {
			existPairs[p.Name] = true
		}
		for _, p := range v.Virtual {
			existPairs[p.Name] = true
		}

		for _, ps := range v.Pairs {
			if existPairs[ps] {
				continue
			}
			log.Fatalf("Operation [%s] requires Pair [%s] support, please add virtual implementation for this pair.", v.Name, ps)
		}
	}
}

// FormatService will formatGlobal services from service spec
func (d *Data) FormatService(s specs.Service) *Service {
	d.serviceSpec = s

	srv := &Service{
		Name:  s.Name,
		pairs: mergePairs(d.Pairs, d.FormatPairs(s.Pairs, false)),
		Infos: mergeInfos(d.Infos, d.FormatInfos(s.Infos, false)),
	}

	for _, v := range s.Namespaces {
		ns := d.FormatNamespace(srv, v)

		srv.Namespaces = append(srv.Namespaces, ns)
	}

	return srv
}

// Sort will sort the data.
func (d *Data) Sort() {
	// Sort all specs.
	d.pairSpec.Sort()
	d.infoSpec.Sort()
	d.operationsSpec.Sort()

	d.serviceSpec.Sort()

	if d.Service != nil {
		d.Service.Sort()
	}
}

// FormatData will formatGlobal the whole data.
func FormatData(p specs.Pairs, m specs.Infos, o specs.Operations) *Data {
	data := &Data{
		pairSpec:       p,
		infoSpec:       m,
		operationsSpec: o,
	}
	data.Pairs = data.FormatPairs(p, true)
	data.Infos = data.FormatInfos(m, true)
	data.Interfaces, data.interfacesMap = data.FormatOperations(o)

	return data
}

func mergePairs(global, service map[string]*Pair) map[string]*Pair {
	ans := make(map[string]*Pair)
	for k, v := range global {
		v := v
		ans[k] = v
	}
	for k, v := range service {
		if p, ok := ans[k]; ok {
			v.Conflict = true
			if v.ptype == p.ptype {
				log.Warnf("pair conflict: %s", k)
			} else {
				log.Fatalf("pair (%s, %s) conflicts with global pair (%s, %s)", k, v.ptype, k, p.ptype)
			}
		}
		v := v
		ans[k] = v
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

func formatDescription(name, desc string) string {
	desc = strings.Trim(desc, "\n")
	if name == "" {
		return strings.ReplaceAll(desc, "\n", "\n// ")
	}
	return fmt.Sprintf("// %s %s", name, strings.ReplaceAll(desc, "\n", "\n// "))
}
