package attr

import (
	"encoding/xml"
	"fmt"
)

type ErrInvalidName struct {
	Actual   xml.Name
	Expected xml.Name
}

var _ error = ErrInvalidName{}

func (err ErrInvalidName) Error() string {
	return fmt.Sprintf("attribute %q has invalid name, expected %q", FmtXMLName(err.Actual), FmtXMLName(err.Expected))
}

type ErrInvalidValue struct {
	Attr    xml.Attr
	Message string
}

var _ error = ErrInvalidValue{}

func (err ErrInvalidValue) Error() string {
	return fmt.Sprintf("attribute %q's value %q is invalid: %s", FmtXMLName(err.Attr.Name), err.Attr.Value, err.Message)
}

func FmtXMLName(n xml.Name) string {
	s := n.Local
	if n.Space != "" {
		s = fmt.Sprintf("%s:%s", n.Space, n.Local)
	}
	return s
}
