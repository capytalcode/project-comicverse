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


type BaseElement struct {
	next       Element
	prev       Element
	parent     Element
	fisrtChild Element
	lastChild  Element
	childCount uint
}

func (e *BaseElement) NextSibling() Element {
	return e.next
}

func (e *BaseElement) SetNextSibling(v Element) {
	e.next = v
}

func (e *BaseElement) PreviousSibling() Element {
	return e.prev
}

func (e *BaseElement) SetPreviousSibling(v Element) {
	e.prev = v
}

func (e *BaseElement) Parent() Element {
	return e.parent
}

func (e *BaseElement) SetParent(v Element) {
	e.parent = v
}

func (e *BaseElement) HasChildren() bool {
	return e.fisrtChild != nil
}

func (e *BaseElement) ChildCount() uint {
	return e.childCount
}

func (e *BaseElement) FirstChild() Element {
	return e.fisrtChild
}

func (e *BaseElement) LastChild() Element {
	return e.lastChild
}

func (e *BaseElement) AppendChild(self, v Element) {
	ensureIsolated(v)

	if e.fisrtChild == nil {
		e.fisrtChild = v
		v.SetNextSibling(nil)
		v.SetPreviousSibling(nil)
	} else {
		l := e.lastChild
		l.SetNextSibling(v)
		v.SetPreviousSibling(l)
	}

	v.SetParent(self)
	e.lastChild = v
	e.childCount++
}

func (e *BaseElement) RemoveChild(self, v Element) {
	if v.Parent() != self {
		return
	}

	if e.childCount <= 0 {
		e.childCount--
	}

	prev := v.PreviousSibling()
	next := v.NextSibling()

	if prev != nil {
		prev.SetNextSibling(next)
	} else {
		e.fisrtChild = next
	}

	if next != nil {
		next.SetNextSibling(prev)
	} else {
		e.lastChild = prev
	}

	v.SetParent(nil)
	v.SetNextSibling(nil)
	v.SetPreviousSibling(nil)
}

func (e *BaseElement) RemoveChildren(_ Element) {
	for c := e.fisrtChild; c != nil; {
		c.SetParent(nil)
		c.SetPreviousSibling(nil)
		next := c.NextSibling()
		c.SetNextSibling(nil)
		c = next
	}
	e.fisrtChild = nil
	e.lastChild = nil
	e.childCount = 0
}

func (e *BaseElement) ReplaceChild(self, v1, insertee Element) {
	e.InsertBefore(self, v1, insertee)
	e.RemoveChild(self, v1)
}

func (e *BaseElement) InsertAfter(self, v1, insertee Element) {
	e.InsertBefore(self, v1.NextSibling(), insertee)
}

func (e *BaseElement) InsertBefore(self, v1, insertee Element) {
	e.childCount++
	if v1 == nil {
		e.AppendChild(self, insertee)
		return
	}

	ensureIsolated(insertee)

	if v1.Parent() == self {
		c := v1
		prev := c.PreviousSibling()
		if prev != nil {
			prev.SetNextSibling(insertee)
			insertee.SetPreviousSibling(prev)
		} else {
			e.fisrtChild = insertee
			insertee.SetPreviousSibling(nil)
		}
		insertee.SetNextSibling(c)
		c.SetPreviousSibling(insertee)
		insertee.SetParent(self)
	}
}

func ensureIsolated(e Element) {
	if p := e.Parent(); p != nil {
		p.RemoveChild(p, e)
	}
}

type ElementKind string

