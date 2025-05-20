package ast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"slices"
)

type Element interface {
	Name() ElementName
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

	xml.Unmarshaler
}

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

func (e *BaseElement) UnmarshalChildren(self Element, d *xml.Decoder, start xml.StartElement) error {
	elErr := fmt.Errorf("unable to unmarshal element kind %q", self.Kind())

	if n := self.Name(); n != (xml.Name{}) {
		if n != start.Name {
			return errors.Join(
				elErr,
				fmt.Errorf("element has different name (%q) than expected (%q)",
					fmtXMLName(start.Name), fmtXMLName(n)),
			)
		}
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}

		switch tt := token.(type) {
		case xml.StartElement:
			elErr = errors.Join(elErr, fmt.Errorf("unable to unmarshal child element %q", fmtXMLName(tt.Name)))

			i := slices.IndexFunc(tt.Attr, func(a xml.Attr) bool {
				return a.Name == elementKindAttrName
			})
			if i == -1 {
				return errors.Join(elErr, fmt.Errorf("element kind not specified"))
			}

			var k ElementKind
			if err := k.UnmarshalXMLAttr(tt.Attr[i]); err != nil {
				return errors.Join(elErr, err)
			}

			c := k.Element()

			err := d.DecodeElement(c, &tt)
			if err != nil && err != io.EOF {
				return err
			}

			e.AppendChild(self, c)
		}
	}
}

func ensureIsolated(e Element) {
	if p := e.Parent(); p != nil {
		p.RemoveChild(p, e)
	}
}

type (
	ElementName = xml.Name
	ElementKind string
)

func NewElementKind(kind string, e Element) ElementKind {
	k := ElementKind(kind)
	if _, ok := elementKindList[k]; ok {
		panic(fmt.Sprintf("Element kind %q is already registered", k))
	}
	elementKindList[k] = e
	return k
}

var (
	_ xml.MarshalerAttr   = ElementKind("")
	_ xml.UnmarshalerAttr = (*ElementKind)(nil)
	_ fmt.Stringer        = ElementKind("")
)

func (k ElementKind) MarshalXMLAttr(n xml.Name) (xml.Attr, error) {
	if n != elementKindAttrName {
		return xml.Attr{}, ErrInvalidAttrName{Actual: n, Expected: elementKindAttrName}
	}

	return xml.Attr{
		Name:  elementKindAttrName,
		Value: k.String(),
	}, nil
}

func (k *ElementKind) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Name != elementKindAttrName {
		return ErrInvalidAttrName{Actual: attr.Name, Expected: elementKindAttrName}
	}

	ak := ElementKind(attr.Value)
	if _, ok := elementKindList[ak]; !ok {
		return ErrInvalidAttrValue{Attr: attr, Message: "element kind not registered"}
	}

	*k = ak

	return nil
}

func (k ElementKind) String() string {
	return string(k)
}

func (k ElementKind) Element() Element {
	return elementKindList[k]
}

var (
	elementKindList     = make(map[ElementKind]Element)
	elementKindAttrName = xml.Name{Local: "data-ipub-element"}
)
