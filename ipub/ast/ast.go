package ast

import (
	"fmt"
)

type Node interface {
	Kind() NodeKind

	NextSibling() Node
	SetNextSibling(Node)

	PreviousSibling() Node
	SetPreviousSibling(Node)

	Parent() Node
	SetParent(Node)

	HasChildren() bool
	ChildCount() uint

	FirstChild() Node
	LastChild() Node

	AppendChild(self, v Node)
	RemoveChild(self, v Node)

	RemoveChildren(self Node)
	ReplaceChild(self, v1, insertee Node)

	InsertBefore(self, v1, insertee Node)
	InsertAfter(self, v1, insertee Node)
}

type BaseNode struct {
	next       Node
	prev       Node
	parent     Node
	fisrtChild Node
	lastChild  Node
	childCount uint
}

func (e *BaseNode) NextSibling() Node {
	return e.next
}

func (e *BaseNode) SetNextSibling(v Node) {
	e.next = v
}

func (e *BaseNode) PreviousSibling() Node {
	return e.prev
}

func (e *BaseNode) SetPreviousSibling(v Node) {
	e.prev = v
}

func (e *BaseNode) Parent() Node {
	return e.parent
}

func (e *BaseNode) SetParent(v Node) {
	e.parent = v
}

func (e *BaseNode) HasChildren() bool {
	return e.fisrtChild != nil
}

func (e *BaseNode) ChildCount() uint {
	return e.childCount
}

func (e *BaseNode) FirstChild() Node {
	return e.fisrtChild
}

func (e *BaseNode) LastChild() Node {
	return e.lastChild
}

func (e *BaseNode) AppendChild(self, v Node) {
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

func (e *BaseNode) RemoveChild(self, v Node) {
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

func (e *BaseNode) RemoveChildren(_ Node) {
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

func (e *BaseNode) ReplaceChild(self, v1, insertee Node) {
	e.InsertBefore(self, v1, insertee)
	e.RemoveChild(self, v1)
}

func (e *BaseNode) InsertAfter(self, v1, insertee Node) {
	e.InsertBefore(self, v1.NextSibling(), insertee)
}

func (e *BaseNode) InsertBefore(self, v1, insertee Node) {
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

func ensureIsolated(e Node) {
	if p := e.Parent(); p != nil {
		p.RemoveChild(p, e)
	}
}

type NodeKind string

func NewNodeKind(kind string, e Node) NodeKind {
	k := NodeKind(kind)
	if _, ok := elementKindList[k]; ok {
		panic(fmt.Sprintf("Node kind %q is already registered", k))
	}
	elementKindList[k] = e
	return k
}

var elementKindList = make(map[NodeKind]Node)
