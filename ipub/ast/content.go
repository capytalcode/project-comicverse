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

func (e *Content) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return e.MarshalXMLElement(e, enc, start)
}

func (e *Content) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	return e.UnmarshalXMLElement(e, dec, start)
}

