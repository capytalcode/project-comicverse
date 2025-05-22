package element

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"forge.capytal.company/capytalcode/project-comicverse/ipub/element/attr"
)

type Element interface {
	Kind() ElementKind
}

type ElementKind string

func NewElementKind(n string, s Element) ElementKind {
	k := ElementKind(n)

	if _, ok := elementKindList[k]; ok {
		panic(fmt.Sprintf("element kind %q already registered", n))
	}

	elementKindList[k] = s
	return k
}

func (k ElementKind) MarshalXMLAttr(n xml.Name) (xml.Attr, error) {
	if n != elementKindAttrName {
		return xml.Attr{}, attr.ErrInvalidName{Actual: n, Expected: elementKindAttrName}
	}
	return xml.Attr{Name: elementKindAttrName, Value: k.String()}, nil
}

func (k *ElementKind) UnmarshalXMLAttr(a xml.Attr) error {
	ak := ElementKind(a.Value)

	if _, ok := elementKindList[ak]; !ok {
		v := make([]string, 0, len(elementKindList))
		for k := range elementKindList {
			v = append(v, k.String())
		}
		return attr.ErrInvalidValue{
			Attr:    a,
			Message: fmt.Sprintf("must be a registered element (%q)", strings.Join(v, `", "`)),
		}
	}

	*k = ak

	return nil
}

func (k ElementKind) String() string {
	return string(k)
}

var (
	elementKindList     = make(map[ElementKind]Element)
	elementKindAttrName = xml.Name{Local: "data-ipub-element"}
)
