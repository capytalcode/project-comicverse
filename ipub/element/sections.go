package element

import "encoding/xml"

type Section struct {
	XMLName xml.Name `xml:"html"`
	Body    Body     `xml:"body"`
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Test    string   `xml:"test,attr"`
}

type Paragraph struct {
	XMLName     xml.Name    `xml:"p"`
	DataElement ElementKind `xml:"data-ipub-element,attr"`
	Test        string      `xml:"test,attr"`

	Text string `xml:",chardata"`
}

