package grss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dimchansky/utfbom"
	"github.com/nbio/xml"
	"io"
)

type Type int

const (
	TypeUnknown Type = 0
	TypeJSON    Type = 1
	TypeXML     Type = 2
	TypeXMLAtom Type = 4
	TypeXMLRss  Type = 8
)

func DetectType(r io.Reader) Type {
	var b = make([]byte, 1)

	for {
		n, err := r.Read(b)
		if err != nil || n != 1 {
			return TypeUnknown
		}

		switch b[0] {
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

	sr, _ := utfbom.Skip(r)

	var buf bytes.Buffer
	tee := io.TeeReader(sr, &buf)
	mr := io.MultiReader(&buf, sr)
	t := DetectType(tee)
	switch t {
	default:
		return t, nil, fmt.Errorf("unknown type")
	case TypeJSON:
		var f = &JSONFeed{}
		d := json.NewDecoder(mr)
		err := d.Decode(f)
		if err != nil {
			return t, nil, err
		}
		return t, f, nil
	case TypeXML:
		var bufx bytes.Buffer
		teex := io.TeeReader(mr, &bufx)
		mrx := io.MultiReader(&bufx, mr)

		var s struct {
			XMLName xml.Name
		}

		err := newXmlDecoder(teex).Decode(&s)
		if err != nil {
			return t, nil, err
		}

		switch s.XMLName.Local {
		case "feed":
			t |= TypeXMLAtom
			var f = &AtomFeed{}
			err := newXmlDecoder(mrx).Decode(f)
			if err != nil {
				return t, nil, err
			}
			return t, f, nil
		case "rss", "RDF":
			t |= TypeXMLRss
			var f = &RssFeed{}
			err := newXmlDecoder(mrx).Decode(f)
			if err != nil {
				return t, nil, err
			}
			return t, f, nil
		default:
			return t, nil, fmt.Errorf("unknown xml")
		}
	}
}
