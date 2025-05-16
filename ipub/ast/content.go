package ast

import (
	"encoding/xml"
	"io"
)

type Content struct {
	BaseElement
}

var KindContent = NewElementKind("content", &Content{})

func (n *Content) Kind() ElementKind {
	return KindContent
}

func (n *Content) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return n.UnmarshalChildren(n, d, start)
}

