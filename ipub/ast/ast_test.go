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

func TestUnmarshal(t *testing.T) {
	assert := tinyssert.New(tinyssert.WithTest(t), tinyssert.WithPanic())

	s := []byte(`
	<html>
		<body data-ipub-element="body">
			<section data-ipub-element="content">
				<img data-ipub-element="image" src="https://hello.com/world.png"/>
			</section>
		</body>
	</html>
	`)

	var data ast.Section

	err := xml.Unmarshal(s, &data)
	if err != nil && err != io.EOF {
		t.Error(err.Error())
		t.FailNow()
	}

	body := data.Body
	assert.Equal(ast.KindBody, body.Kind())

	t.Logf("%#v", body)

	content := body.FirstChild()
	assert.Equal(ast.KindContent, content.Kind())

	t.Logf("%#v", content)

	img := content.FirstChild().(*ast.Image)
	assert.Equal(ast.KindImage, img.Kind())
	assert.Equal("https://hello.com/world.png", img.Source())

	t.Logf("%#v", img)
}
