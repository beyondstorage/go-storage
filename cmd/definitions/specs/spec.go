//go:build tools
// +build tools

package specs

import (
	"sort"
)

// Interface is the spec for interface.
type Interface struct {
	Name        string
	Description string
	Ops         []Operation
}

// Operations is the spec for operations.
type Operations struct {
	Interfaces []Interface
	Fields     []Field
}

// Sort will sort the Operations
func (o *Operations) Sort() {
	sort.Slice(o.Fields, func(i, j int) bool {
		x := o.Fields
		return x[i].Name < x[j].Name
	})

	sort.Slice(o.Interfaces, func(i, j int) bool {
		x := o.Interfaces
		return x[i].Name < x[j].Name
	})
}

// Operation is the spec for operation.
type Operation struct {
	Name        string
	Description string
	Params      []string
	Pairs       []string
	Results     []string
	ObjectMode  string
	Local       bool
}

// Field is the spec for field.
type Field struct {
	Name string
	Type string
}

// Service is the data parsed from TOML.
type Service struct {
	Name       string
	Namespaces []Namespace
	Pairs      Pairs
	Infos      Infos
}

// Sort will sort the service spec.
func (s *Service) Sort() {
	s.Pairs.Sort()
	s.Infos.Sort()

	sort.Slice(s.Namespaces, func(i, j int) bool {
		ns := s.Namespaces
		return ns[i].Name < ns[j].Name
	})

	for _, v := range s.Namespaces {
		v.Sort()
	}
}

// Features is all global features that available.
//
// Features will be defined in features.toml.
type Features []Feature

func (p Features) Sort() {
	if p == nil || len(p) == 0 {
		return
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].Description < p[j].Description
	})
}

type Feature struct {
	Name        string
	Description string
}

// Infos is the spec for infos.
type Infos []Info

// Sort will sort the pair spec.
func (p Infos) Sort() {
	if p == nil || len(p) == 0 {
		return
	}

	sort.Slice(p, func(i, j int) bool {
		return compareInfoSpec(p[i], p[j])
	})
}

// Info is the spec for info.
type Info struct {
	Scope       string
	Category    string
	Name        string
	Type        string
	Export      bool
	Description string
}

type Pairs []Pair

// Sort will sort the pair spec.
func (p Pairs) Sort() {
	if p == nil || len(p) == 0 {
		return
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].Name < p[j].Name
	})
}

// Pair is the data parsed from TOML.
type Pair struct {
	Name        string
	Type        string
	Defaultable bool
	Description string
}

// Namespace is the data parsed from TOML.
type Namespace struct {
	Name      string
	Features  []string // The feature names that provided by current namespace.
	Implement []string
	New       New
	Op        []Op
}

// Sort will sort the Namespace
func (n *Namespace) Sort() {
	n.New.Sort()

	sort.Strings(n.Features)
	sort.Strings(n.Implement)

	sort.Slice(n.Op, func(i, j int) bool {
		x := n.Op
		return x[i].Name < x[j].Name
	})

	for _, v := range n.Op {
		v.Sort()
	}
}

// Op means an operation definition.
type Op struct {
	Name string

	Required []string
	Optional []string
}

// Sort will sort the Op
func (o *Op) Sort() {
	sort.Strings(o.Required)
	sort.Strings(o.Optional)
}

// New is the spec for new function.
type New struct {
	Required []string
	Optional []string
}

// Sort will sort the New
func (o *New) Sort() {
	sort.Strings(o.Required)
	sort.Strings(o.Optional)
}

func compareInfoSpec(x, y Info) bool {
	if x.Scope != y.Scope {
		return x.Scope < y.Scope
	}
	if x.Category != y.Category {
		return x.Category < y.Category
	}
	return x.Name < y.Name
}
