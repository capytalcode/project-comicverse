package element_test

import (
	"encoding/xml"
	"testing"

	"forge.capytal.company/capytalcode/project-comicverse/ipub/element"
)

func Test(t *testing.T) {
	d := element.Section{
		Body: element.Body{
			Test: "helloworld",
			Children: []element.Element{
				&element.Paragraph{
					DataElement: element.ParagraphKind,
					Text:        "hello world",
					Test:        "testvalue",
				},
				&element.Paragraph{
					DataElement: element.ParagraphKind,

					Text: "hello world 2",
				},
			},
		},
	}

	b, err := xml.Marshal(d)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	t.Logf("%#v", string(b))

	var ud element.Section
	err = xml.Unmarshal(b, &ud)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("%#v", ud)
	t.Logf("%#v", ud.Body.Children[0])
	t.Logf("%#v", ud.Body.Children[1])
}
