package attr

import (
	"encoding/xml"
	"fmt"
)

type Attribute interface {
	xml.MarshalerAttr
	xml.UnmarshalerAttr
	fmt.Stringer
}

type BaseAttribute string

func (a BaseAttribute) MarshalXMLAttr(n xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: n, Value: a.String()}, nil
}

func (a *BaseAttribute) UnmarshalXMLAttr(attr xml.Attr) error {
	*a = BaseAttribute(attr.Value)
	return nil
}

func (a BaseAttribute) String() string {
	return string(a)
}


