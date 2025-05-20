package ast_test

import (
	_ "embed"
	"encoding/xml"
	"io"
	"testing"

	"forge.capytal.company/capytalcode/project-comicverse/ipub/ast"
	"forge.capytal.company/loreddev/x/tinyssert"
)

//go:embed test.xml
var test []byte

func TestMarshal(t *testing.T) {
	b := &ast.Body{}
	c := &ast.Content{}
	i := &ast.Image{}
	i.SetSource("https://hello.com/world.png")
	c.AppendChild(c, i)
	b.AppendChild(b, c)

	s := ast.Section{
		Body: b,
	}
	by, err := xml.Marshal(s)

	if err != nil && err != io.EOF {
		t.Error(err.Error())
		t.FailNow()
	}

	// t.Logf("%#v", s.Body)
	//
	// t.Logf("%#v", f)

	t.Logf("%#v", string(by))
}

