package definitions

type Type struct {
	Expr    string // Raw expr that before type name, e.g. `[]`, `...`, `[]*`
	Package string
	Name    string
}

func (t Type) FullName() string {
	if t.Package == "" {
		return t.Expr + t.Name
	}
	return t.Expr + t.Package + "." + t.Name
}
