package main

// InterfaceSpec is the spec for interface.
type InterfaceSpec struct {
	Name        string           `hcl:",label"`
	Description string           `hcl:"description,optional"`
	Internal    bool             `hcl:"internal,optional"`
	Embed       []string         `hcl:"embed,optional"`
	Ops         []*OperationSpec `hcl:"op,block"`
}

// OperationsSpec is the spec for operations.
type OperationsSpec struct {
	Interfaces []*InterfaceSpec `hcl:"interface,block"`
	Fields     []*FieldSpec     `hcl:"field,block"`
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

// InfosSpec is the spec for infos.
type InfosSpec struct {
	Infos []*InfoSpec `hcl:"info,block"`
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

// OpSpec means an operation definition.
type OpSpec struct {
	Name     string   `hcl:",label"`
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`
}

// NewSpec is the spec for new function.
type NewSpec struct {
	Required []string `hcl:"required,optional"`
	Optional []string `hcl:"optional,optional"`
}
