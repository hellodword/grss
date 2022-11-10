package grss

import (
	"encoding/xml"
	"fmt"
	"github.com/hellodword/grss/pkg/etree"
	"golang.org/x/text/encoding/ianaindex"
	"io"
)

func AtomParse(r io.Reader) (*AtomFeed, error) {
	var f = &AtomFeed{}
	d := xml.NewDecoder(r)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
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

	f.PatchNS()
	return f, nil
}

func (f *AtomFeed) Decode(r io.Reader) error {

	doc := etree.NewDocument()
	doc.ReadSettings.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		return enc.NewDecoder().Reader(input), nil
	}

	_, err := doc.ReadFrom(r)
	if err != nil {
		return err
	}

	f.Inner = doc

	return nil
}

func (f *AtomFeed) Encode() (string, error) {
	return f.Inner.WriteToString()
}
