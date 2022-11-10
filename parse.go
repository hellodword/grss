package grss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hellodword/grss/common/header"
	"github.com/hellodword/grss/pkg/etree"
	"io"
)

type Type int

const (
	TypeUnknown = 0
	TypeJSON    = 1
	TypeXML     = 2
)

func CheckType(r io.Reader) Type {
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

func Parse(r io.Reader) (Type, interface{}, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	t := CheckType(tee)

	rm := io.MultiReader(&buf, r)
	switch t {
	default:
		return t, nil, fmt.Errorf("unknown type")
	case TypeJSON:
		var j = &JSONFeed{}
		d := json.NewDecoder(rm)
		err := d.Decode(j)
		if err != nil {
			return t, nil, err
		}
		return t, j, nil
	case TypeXML:
		doc := etree.NewDocument()
		_, err := doc.ReadFrom(rm)
		if err != nil {
			return t, nil, err
		}
		if doc.SelectElement("RDF") != nil {
			// TODO
		}

		if doc.SelectElement("rss") != nil {
			// TODO
		}

		if doc.SelectElement("feed") != nil {
			// TODO
		}

	}

	return t, nil, nil
}
