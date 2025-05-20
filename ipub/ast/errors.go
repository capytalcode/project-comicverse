package ast

import (
	"encoding/xml"
	"fmt"
)

type ErrInvalidAttrName struct {
	Actual   xml.Name
	Expected xml.Name
}

var _ error = ErrInvalidAttrName{}

func (err ErrInvalidAttrName) Error() string {
	return fmt.Sprintf("attribute %q has invalid name, expected %q", fmtXMLName(err.Expected), fmtXMLName(err.Actual))
}

type ErrInvalidAttrValue struct {
	Attr    xml.Attr
	Message string
}

var _ error = ErrInvalidAttrValue{}

func (err ErrInvalidAttrValue) Error() string {
	return fmt.Sprintf("attribute %q's value %q is invalid: %s", fmtXMLName(err.Attr.Name), err.Attr.Value, err.Message)
}

func fmtXMLName(n xml.Name) string {
	s := n.Local
	if n.Space != "" {
		s = fmt.Sprintf("%s:%s", n.Space, n.Local)
	}
	return s
}
