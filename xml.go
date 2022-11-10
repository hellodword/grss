package grss

import "encoding/xml"

type XmlGeneric struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr,omitempty"`
	Content    string     `xml:",innerxml"`
}
