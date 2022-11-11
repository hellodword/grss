package grss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hellodword/grss/common/header"
	"github.com/hellodword/grss/pkg/etree"
	"golang.org/x/text/encoding/ianaindex"
	"io"
)

type Type int

const (
	TypeUnknown Type = 0
	TypeJSON    Type = 1
	TypeXML     Type = 2
)

func DetectType(r io.Reader) Type {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	bom := header.ReadBom(tee)

	rm := io.MultiReader(&buf, r)

	var b = make([]byte, len(bom)+1)
	for {
		_, err := rm.Read(b)
		if err != nil {
			return TypeUnknown
		}
		if len(b) <= len(bom) {
			continue
		}

		switch b[len(bom)] {
		case ' ', '\n', '\r', '\t':
			continue
		case '{':
			return TypeJSON
		case '<':
			return TypeXML
		default:
			return TypeUnknown
		}
	}
}

func Parse(r io.Reader) (Type, Feed, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	t := DetectType(tee)

	rm := io.MultiReader(&buf, r)
	switch t {
	default:
		return t, nil, fmt.Errorf("unknown type")
	case TypeJSON:
		var f = &JSONFeed{}
		d := json.NewDecoder(rm)
		err := d.Decode(f)
		if err != nil {
			return t, nil, err
		}
		return t, f, nil
	case TypeXML:
		doc := etree.NewDocument()
		doc.ReadSettings.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			enc, err := ianaindex.IANA.Encoding(charset)
			if err != nil {
				return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
			}
			return enc.NewDecoder().Reader(input), nil
		}
		_, err := doc.ReadFrom(rm)
		if err != nil {
			return t, nil, err
		}
		if root := doc.SelectElement("RDF"); root != nil {
			var f = &RssFeed{}
			f.Inner = doc

			err = f.From(root)
			if err != nil {
				return t, nil, err
			}

			return t, f, nil
		}

		if root := doc.SelectElement("rss"); root != nil {
			var f = &RssFeed{}
			f.Inner = doc

			err = f.From(root)
			if err != nil {
				return t, nil, err
			}

			return t, f, nil
		}

		if root := doc.SelectElement("feed"); root != nil {
			var f = &AtomFeed{}
			f.Inner = doc

			err = f.From(root)
			if err != nil {
				return t, nil, err
			}

			return t, f, nil
		}

		return t, nil, fmt.Errorf("unknown xml")
	}
}

func (f *RssFeed) From(root *etree.Element) error {
	channel := root.SelectElement("channel")
	if channel == nil {

	}

	return nil
}

func (f *AtomFeed) From(root *etree.Element) error {

	return nil
}
