package ast

import (
	"encoding/xml"
)

type Section struct {
	XMLName xml.Name `xml:"html"`
	Body    *Body    `xml:"body"`
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

func (e *Body) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return e.MarshalXMLElement(e, enc, start)
}

func (e *Body) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	return e.UnmarshalXMLElement(e, dec, start)
}
