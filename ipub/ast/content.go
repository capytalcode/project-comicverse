package ast

type Content struct {
	BaseNode
}

var KindContent = NewNodeKind("content", &Content{})

func (e Content) Kind() NodeKind {
	return KindContent
}

type Image struct {
	src string

	BaseNode
}

var KindImage = NewNodeKind("image", &Image{})

func (e *Image) Kind() NodeKind {
	return KindImage
}

func (e Image) Source() string {
	return e.src
}

func (e *Image) SetSource(src string) {
	e.src = src
}
