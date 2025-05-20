package ast

import (
	"encoding/xml"
	"slices"
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

type Image struct {
	src string

	BaseElement
}

var KindImage = NewElementKind("image", &Image{})

func (e *Image) Name() ElementName {
	return ElementName{Local: "img"}
}

func (e *Image) Kind() ElementKind {
	return KindImage
}

func (e *Image) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if src := e.Source(); src != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "src"},
			Value: src,
		})
	}
	return e.MarshalXMLElement(e, enc, start)
}

func (e *Image) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	i := slices.IndexFunc(start.Attr, func(a xml.Attr) bool {
		return a.Name == xml.Name{Local: "src"}
	})
	if i > -1 {
		e.SetSource(start.Attr[i].Value)
	}
	return e.UnmarshalXMLElement(e, dec, start)
}

func (e Image) Source() string {
	return e.src
}

func (e *Image) SetSource(src string) {
	e.src = src
}
