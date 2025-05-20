package ast

import (
	"encoding/xml"
	"io"
)

type Content struct {
	BaseElement
}

var KindContent = NewElementKind("content", &Content{})

func (e Content) Name() ElementName {
	return ElementName{Local: "section"}
}

func (e Content) Kind() ElementKind {
	return KindContent
}

func (n *Content) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return n.UnmarshalChildren(n, d, start)
}

