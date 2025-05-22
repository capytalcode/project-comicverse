package element

import "encoding/xml"

type Section struct {
	XMLName xml.Name `xml:"html"`
	Body    Body     `xml:"body"`
}

var KindSection = NewElementKind("section", Section{})

func (Section) Kind() ElementKind {
	return KindSection
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Test    string   `xml:"test,attr"`
}

var KindBody = NewElementKind("body", Body{})

func (Body) Kind() ElementKind {
	return KindBody
}

type Paragraph struct {
	XMLName     xml.Name    `xml:"p"`
	DataElement ElementKind `xml:"data-ipub-element,attr"`
	Test        string      `xml:"test,attr"`

	Text string `xml:",chardata"`
}

var KindParagraph = NewElementKind("paragraph", Paragraph{})

func (Paragraph) Kind() ElementKind {
	return KindParagraph
}
