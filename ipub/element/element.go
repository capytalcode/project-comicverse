package element

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"forge.capytal.company/capytalcode/project-comicverse/ipub/element/attr"
)

type Element interface {
	Kind() ElementKind
}

type ElementChildren []Element

func (ec *ElementChildren) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	elErr := fmt.Errorf("unable to unsmarshal element %q", attr.FmtXMLName(start.Name))

	i := slices.IndexFunc(start.Attr, func(a xml.Attr) bool {
		return a.Name == elementKindAttrName
	})
	if i == -1 {
		return errors.Join(elErr, fmt.Errorf("element kind not specified"))
	}

	var k ElementKind
	if err := k.UnmarshalXMLAttr(start.Attr[i]); err != nil {
		return err
	}

	ks := elementKindList[k]

	// Get a pointer of a new instance of the underlying implementation so we can
	// change it without manipulating the value inside the elementKindList.
	ep := reflect.New(reflect.TypeOf(ks))
	if ep.Elem().Kind() == reflect.Pointer {
		// If the implementation is a pointer, we need the underlying value so we can
		// manipulate it.
		ep = reflect.New(reflect.TypeOf(ks).Elem())
	}

	if err := d.DecodeElement(ep.Interface(), &start); err != nil && err != io.EOF {
		return errors.Join(elErr, err)
	}

	if ec == nil {
		c := ElementChildren{}
		ec = &c
	}

	s := *ec
	s = append(s, ep.Interface().(Element))
	*ec = s

	return nil
}

type ElementKind string

// NewElementKind registers a new Element implementation to a private list which is
// consumed bu [ElementChildren] to properly find what underlying type is a children
// of another element struct.
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
