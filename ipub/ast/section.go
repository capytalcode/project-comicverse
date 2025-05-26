package ast

import (
	"encoding/xml"
)

type Section struct {
	XMLName xml.Name `xml:"html"`
	Body    *Body    `xml:"body"`
}

type Body struct {
	BaseNode
}

var KindBody = NewNodeKind("body", &Body{})

func (e Body) Kind() NodeKind {
	return KindBody
}
