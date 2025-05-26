package ast

type Package struct {
	BaseNode
}

var KindPackage = NewNodeKind("package", &Package{})

func (e Package) Kind() NodeKind {
	return KindPackage
}
