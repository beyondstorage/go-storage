package definitions

type Type struct {
	Expr    string // Raw expr that before type name, e.g. `[]`, `...`, `[]*`
	Package string
	Name    string
}

func (t Type) FullName(pkg ...string) string {
	// The type is a builtin type, we can use directly.
	if t.Package == "" {
		return t.Expr + t.Name
	}
	// The types package name is the same with input one.
	if len(pkg) > 0 && t.Package == pkg[0] {
		return t.Expr + t.Name
	}
	return t.Expr + t.Package + "." + t.Name
}
