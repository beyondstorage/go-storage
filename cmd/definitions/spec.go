package main

import (
	"sort"
)

// InterfaceSpec is the spec for interface.
type InterfaceSpec struct {
	Name        string           `hcl:",label"`
	Description string           `hcl:"description,optional"`
	Internal    bool             `hcl:"internal,optional"`
	Embed       []string         `hcl:"embed,optional"`
	Ops         []*OperationSpec `hcl:"op,block"`
}

// Sort will sort the InterfaceSpec
func (i *InterfaceSpec) Sort() {
	sort.Strings(i.Embed)
}

// OperationsSpec is the spec for operations.
type OperationsSpec struct {
	Interfaces []*InterfaceSpec `hcl:"interface,block"`
	Fields     []*FieldSpec     `hcl:"field,block"`
}

// Sort will sort the OperationsSpec
func (o *OperationsSpec) Sort() {
	sort.Slice(o.Fields, func(i, j int) bool {
		x := o.Fields
		return x[i].Name < x[j].Name
	})

	sort.Slice(o.Interfaces, func(i, j int) bool {
		x := o.Interfaces
		return x[i].Name < x[j].Name
	})

	for _, v := range o.Interfaces {
		v.Sort()
	}
}

// OperationSpec is the spec for operation.
type OperationSpec struct {
	Name        string   `hcl:",label"`
	Description string   `hcl:"description,optional"`
	Params      []string `hcl:"params,optional"`
	Results     []string `hcl:"results,optional"`
}

// FieldSpec is the spec for field.
type FieldSpec struct {
	Name string `hcl:",label"`
	Type string `hcl:"type"`
}

// ServiceSpec is the data parsed from HCL.
type ServiceSpec struct {
	Name       string           `hcl:"name"`
	Namespaces []*NamespaceSpec `hcl:"namespace,block"`
	Pairs      *PairsSpec       `hcl:"pairs,block"`
	Infos      *InfosSpec       `hcl:"infos,block"`
}

// Sort will sort ther service spec.
func (s *ServiceSpec) Sort() {
	s.Pairs.Sort()
	s.Infos.Sort()

	for _, v := range s.Namespaces {
		v.Sort()
	}
}

// InfosSpec is the spec for infos.
type InfosSpec struct {
	Infos []*InfoSpec `hcl:"info,block"`
}

// Sort will sort the pair spec.
func (p *InfosSpec) Sort() {
	if p == nil || len(p.Infos) == 0 {
		return
	}

	sort.Slice(p.Infos, func(i, j int) bool {
		x := p.Infos
		return compareInfoSpec(x[i], x[j])
	})
}

// InfoSpec is the spec for info.
type InfoSpec struct {
	Scope       string `hcl:",label"`
	Category    string `hcl:",label"`
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	DisplayName string `hcl:"display_name,optional"`
	ZeroValue   string `hcl:"zero_value,optional"`
}

// PairsSpec is the data parsed from HCL.
type PairsSpec struct {
	Pairs []*PairSpec `hcl:"pair,block"`
}

// Sort will sort the pair spec.
func (p *PairsSpec) Sort() {
	if p == nil || len(p.Pairs) == 0 {
		return
	}

	sort.Slice(p.Pairs, func(i, j int) bool {
		x := p.Pairs
		return x[i].Name < x[j].Name
	})
}

// PairSpec is the data parsed from HCL.
type PairSpec struct {
	Name        string `hcl:",label"`
	Type        string `hcl:"type"`
	Description string `hcl:"description,optional"`
	Parser      string `hcl:"parser,optional"`
	Default     string `hcl:"default,optional"`
}

// NamespaceSpec is the data parsed from HCL.
type NamespaceSpec struct {
	Name      string    `hcl:",label"`
	Implement []string  `hcl:"implement,optional"`
	New       *NewSpec  `hcl:"new,block"`
	Op        []*OpSpec `hcl:"op,block"`
}

// Sort will sort the NamespaceSpec
func (n *NamespaceSpec) Sort() {
	n.New.Sort()

	sort.Strings(n.Implement)

	sort.Slice(n.Op, func(i, j int) bool {
		x := n.Op
		return x[i].Name < x[j].Name
	})

	for _, v := range n.Op {
		v.Sort()
	}
}

// OpSpec means an operation definition.
type OpSpec struct {
	Name     string   `hcl:",label"`
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`
}

// Sort will sort the OpSpec
func (o *OpSpec) Sort() {
	sort.Strings(o.Required)
	sort.Strings(o.Optional)
}

// NewSpec is the spec for new function.
type NewSpec struct {
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`
}

// Sort will sort the NewSpec
func (o *NewSpec) Sort() {
	sort.Strings(o.Required)
	sort.Strings(o.Optional)
}

func compareInfoSpec(x, y *InfoSpec) bool {
	if x.Scope != y.Scope {
		return x.Scope < y.Scope
	}
	if x.Category != y.Category {
		return x.Category < y.Category
	}
	return x.Name < y.Name
}
