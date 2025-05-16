package ast

import (
	"encoding/xml"
	"fmt"
	"io"
	"slices"
)

type Element interface {
	Kind() ElementKind

	NextSibling() Element
	SetNextSibling(Element)

	PreviousSibling() Element
	SetPreviousSibling(Element)

	Parent() Element
	SetParent(Element)

	HasChildren() bool
	ChildCount() uint

	FirstChild() Element
	LastChild() Element

	AppendChild(self, v Element)
	RemoveChild(self, v Element)

	RemoveChildren(self Element)
	ReplaceChild(self, v1, insertee Element)

	InsertBefore(self, v1, insertee Element)
	InsertAfter(self, v1, insertee Element)

}
