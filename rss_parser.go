package grss

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/ianaindex"
	"io"
)

func RssParse(r io.Reader) (*RssFeed, error) {
	var f = &RssFeed{}
	d := xml.NewDecoder(r)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		f.Charset = charset
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		return enc.NewDecoder().Reader(input), nil
	}

	err := d.Decode(f)
	if err != nil {
		return nil, err
	}

	f.PatchXmlns()
	return f, nil
}
