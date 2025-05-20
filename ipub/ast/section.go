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

func (e Body) Name() ElementName {
	return ElementName{Local: "body"}
}

func (e Body) Kind() ElementKind {
	return KindBody
}

func (n *Body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return n.UnmarshalChildren(n, d, start)
}
