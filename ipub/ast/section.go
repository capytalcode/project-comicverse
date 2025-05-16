package ast

import (
	"encoding/xml"
)

type Section struct {
	Body Body `xml:"body"`
}

type Body struct {
	BaseElement
}

var KindBody = NewElementKind("body", &Body{})

func (n *Body) Kind() ElementKind {
	return KindBody
}

