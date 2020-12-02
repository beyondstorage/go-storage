// +build tools

package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Xuanwo/templateutils"
	specs "github.com/aos-dev/specs/go"
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
	pairSpec       *specs.Pairs
	infoSpec       *specs.Infos
	operationsSpec *specs.Operations
	serviceSpec    *specs.Service
}

// Service is the service definition.
type Service struct {
	Name       string
	Namespaces []*Namespace
	Pairs      map[string]*Pair
	Infos      []*Info
}

// Sort will sort the service
func (s *Service) Sort() {
	for _, v := range s.Namespaces {
		v.Sort()
	}
}

// Namespace contains all info about a namespace
type Namespace struct {
	Name  string
	New   *Function
	Funcs []*Function
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
	Name    string
	Type    string
	Parser  string
	Default string

	// Runtime generated
	Global      bool
	Description string
}

// Format will formatGlobal current pair
func (p *Pair) Format(s *specs.Pair, global bool) {
	p.Name = s.Name
	p.Type = s.Type
	p.Parser = s.Parser
	p.Default = s.Default

	p.Global = global

	p.Description = formatDescription(templateutils.ToPascal(p.Name), s.Description)
}

// FullName will print full name for current pair
func (p *Pair) FullName() string {
	if p.Global {
		return fmt.Sprintf("\"%s\"", p.Name)
	}
	return "pair" + templateutils.ToPascal(p.Name)
}

// Info is the metadata definition.
type Info struct {
	Scope     string
	Category  string
	Name      string
	Type      string
	ZeroValue string
	Export    bool
	Comment   string

	Global bool

	displayName string
}

// Format will formatGlobal info spec into Info
func (i *Info) Format(s *specs.Info, global bool) {
	i.Scope = s.Scope
	i.Category = s.Category
	i.Name = s.Name
	i.Type = s.Type
	i.displayName = s.DisplayName
	i.ZeroValue = s.ZeroValue
	i.Export = s.Export
	i.Comment = s.Comment

	i.Global = global
}

func (i *Info) TypeName() string {
	if i.Export {
		return templateutils.ToPascal(i.Name)
	} else {
		return templateutils.ToCamel(i.Name)
	}
}
func (i *Info) DisplayName() string {
	if i.displayName != "" {
		return i.displayName
	}
	return templateutils.ToPascal(i.Name)
}

// Interface represents an interface
type Interface struct {
	Name        string
	Description string
	Internal    bool
	Embed       []*Interface
	Ops         map[string]*Operation
}

// NewInterface will create a new interface from spec.
func NewInterface(in *specs.Interface, fields map[string]*Field) *Interface {
	inter := &Interface{
		Name:        in.Name,
		Description: formatDescription(templateutils.ToPascal(in.Name), in.Description),
		Internal:    in.Internal,
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
	if i.Internal {
		return templateutils.ToCamel(i.Name)
	}
	return templateutils.ToPascal(i.Name)
}

// Operation represents an operation.
type Operation struct {
	Name        string
	Description string
	Pairs       []string
	Params      Fields
	Results     Fields
}

// NewOperation will create an new operation from operation spec.
func NewOperation(v *specs.Operation, fields map[string]*Field) *Operation {
	op := &Operation{
		Name:        v.Name,
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
	// Inject error
	op.Results = append(op.Results, fields["err"])

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
		if strings.HasPrefix(v.Type, packageName) {
			s = append(s, v.Name+" "+strings.TrimPrefix(v.Type, packageName+"."))
			continue
		}
		s = append(s, v.String())
	}
	return strings.Join(s, ",")
}

// Function represents a function.
type Function struct {
	*Operation

	Required  []*Pair
	Optional  []*Pair
	Generated []*Pair // This op's generated pairs, they will be treated as optional.

	implemented bool
}

// NewFunction will createn a new function.
func NewFunction(o *Operation) *Function {
	return &Function{Operation: o}
}

// Format will formatGlobal a function with Op.
func (f *Function) Format(s *specs.Op, p map[string]*Pair) {
	for _, v := range s.Required {
		f.Required = append(f.Required, p[v])
	}
	for _, v := range s.Optional {
		f.Optional = append(f.Optional, p[v])
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
	sort.Slice(f.Generated, func(i, j int) bool {
		x := f.Generated
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
		if v.Type == "io.Reader" || v.Type == "io.ReadCloser" {
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
		if v.Name == "seg" {
			x = append(x, "seg.Path()", "seg.ID()")
			continue
		}
		if v.Type != "string" {
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
	Name string
	Type string
}

// String will print field in string formatGlobal.
func (f *Field) String() string {
	if f.Name == "" {
		return f.Type
	}
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}

// Caller will print the caller formatGlobal of field.
func (f *Field) Caller() string {
	if strings.HasPrefix(f.Type, "...") {
		return f.Name + "..."
	}
	return f.Name
}

// Format will create a new field.
func (f *Field) Format(s *specs.Field) {
	f.Type = s.Type
	f.Name = s.Name
}

// FormatPairs will formatGlobal pairs for pair spec
func (d *Data) FormatPairs(p *specs.Pairs, global bool) map[string]*Pair {
	if p == nil {
		return nil
	}

	m := make(map[string]*Pair)
	for _, v := range p.Pairs {
		p := &Pair{}
		p.Format(v, global)

		m[p.Name] = p
	}
	return m
}

// FormatInfos will formatGlobal metas for meta spec
func (d *Data) FormatInfos(m *specs.Infos, global bool) []*Info {
	if m == nil {
		return nil
	}

	is := make([]*Info, 0, len(m.Infos))
	for _, v := range m.Infos {
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
func (d *Data) FormatOperations(o *specs.Operations) (ins []*Interface, inm map[string]*Interface) {
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

	// Handle embed interface
	for _, in := range o.Interfaces {
		for _, v := range in.Embed {
			inm[in.Name].Embed = append(inm[in.Name].Embed, inm[v])
		}
	}
	return
}

// FormatNamespace will formatGlobal a namespace.
func (d *Data) FormatNamespace(srv *Service, n *specs.Namespace) *Namespace {
	ns := &Namespace{Name: n.Name}

	nsInterface := n.Name + "r"

	// Handle New function
	if n.New == nil {
		n.New = &specs.New{}
	}
	ns.New = NewFunction(&Operation{Name: "new"})
	ns.New.Format(&specs.Op{
		Required: n.New.Required,
		Optional: n.New.Optional,
	}, srv.Pairs)

	// Handle other interfaces.
	fns := make(map[string]*Function)

	// Add namespace itself into implements.
	implements := append(n.Implement[:], nsInterface)
	for _, interfaceName := range implements {
		inter := d.interfacesMap[interfaceName]

		for k, v := range inter.Ops {
			v := v

			f := NewFunction(v)
			ns.Funcs = append(ns.Funcs, f)
			fns[k] = f
		}
	}

	for _, v := range n.Op {
		fns[v.Name].Format(v, srv.Pairs)
	}

	implemented := parseFunc(srv.Name, n.Name)
	for _, fn := range fns {
		x := templateutils.ToCamel(fn.Name)
		if _, ok := implemented[x]; ok {
			fn.implemented = true
		}
	}

	// Inject generated pair.
	d.InjectNamespace(srv, ns)
	return ns
}

// InjectNamespace will inject a namespace to insert generated pairs.
func (d *Data) InjectNamespace(srv *Service, n *Namespace) {
	// Inject read_callback_func
	for _, v := range n.Funcs {
		existPairs := map[string]bool{}
		for _, p := range v.Required {
			existPairs[p.Name] = true
		}
		for _, p := range v.Optional {
			existPairs[p.Name] = true
		}

		for _, ps := range v.Pairs {
			if existPairs[ps] {
				continue
			}
			v.Generated = append(v.Generated, srv.Pairs[ps])
		}
	}
}

// FormatService will formatGlobal services from service spec
func (d *Data) FormatService(s *specs.Service) *Service {
	d.serviceSpec = s

	srv := &Service{
		Name:  s.Name,
		Pairs: mergePairs(d.Pairs, d.FormatPairs(s.Pairs, false)),
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

	if d.serviceSpec != nil {
		d.serviceSpec.Sort()
	}
	if d.Service != nil {
		d.Service.Sort()
	}
}

// FormatData will formatGlobal the whole data.
func FormatData(p *specs.Pairs, m *specs.Infos, o *specs.Operations) *Data {
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

func formatDescription(name, desc string) string {
	desc = strings.Trim(desc, "\n")
	if name == "" {
		return strings.ReplaceAll(desc, "\n", "\n// ")
	}
	return fmt.Sprintf("// %s %s", name, strings.ReplaceAll(desc, "\n", "\n// "))
}
