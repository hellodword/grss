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
	Content    string     `xml:",innerxml"`
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
