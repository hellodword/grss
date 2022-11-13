package grss

import (
	"fmt"
	"github.com/nbio/xml"
	"golang.org/x/text/encoding/ianaindex"
	"io"
)

type XmlGeneric struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr,omitempty"`
	XmlText
}

type XmlText struct {
	Text     string `xml:",chardata"`
	Cdata    string `xml:",cdata"`
	InnerXml string `xml:",innerxml"`
}

func (a *XmlText) String() string {
	if a.InnerXml != "" {
		return a.InnerXml
	} else if a.Cdata != "" {
		return a.Cdata
	} else {
		return a.Text
	}
}

func newXmlDecoder(r io.Reader) *xml.Decoder {
	d := xml.NewDecoder(r)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		return enc.NewDecoder().Reader(input), nil
	}
	return d
}

func addAttrs(pre [][3]string, src []xml.Attr) (attrs []xml.Attr) {
	for i := range pre {
		var b bool
		for j := range src {
			a := src[j]
			if a.Value == pre[i][2] && a.Name.Space == pre[i][0] && a.Name.Local == pre[i][1] {
				b = true
				break
			}
		}
		if !b {
			attrs = append(attrs, xml.Attr{
				Name: xml.Name{
					Space: pre[i][0],
					Local: pre[i][1],
				},
				Value: pre[i][2],
			})
		}
	}

	return attrs
}
