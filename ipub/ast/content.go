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

