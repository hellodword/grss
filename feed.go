package grss

import (
	"encoding/json"
	"io"
)

type Feed interface {
	ToJSON() (*JSONFeed, error)
	ToRss2() (*RssFeed, error)
	ToAtom() (*AtomFeed, error)
	WriteOut(w io.Writer) error
}

func (f *JSONFeed) ToJSON() (*JSONFeed, error) {
	return f, nil
}

func (f *JSONFeed) ToRss2() (*RssFeed, error) {
	return nil, nil
}

func (f *JSONFeed) ToAtom() (*AtomFeed, error) {
	return nil, nil
}

func (f *JSONFeed) WriteOut(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(true)
	e.SetIndent("", "    ")
	if f.Inner != nil {
		return e.Encode(f.Inner)
	} else {
		return e.Encode(f)
	}
}

func (f *RssFeed) ToJSON() (*JSONFeed, error) {
	return nil, nil
}

func (f *RssFeed) ToRss2() (*RssFeed, error) {
	return f, nil
}

func (f *RssFeed) ToAtom() (*AtomFeed, error) {
	return nil, nil
}

func (f *RssFeed) WriteOut(w io.Writer) error {
	if f.Inner != nil {
		_, err := f.Inner.WriteTo(w)
		return err
	}
	return nil
}

func (f *AtomFeed) ToJSON() (*JSONFeed, error) {
	return nil, nil
}

func (f *AtomFeed) ToRss2() (*RssFeed, error) {
	return nil, nil
}

func (f *AtomFeed) ToAtom() (*AtomFeed, error) {
	return f, nil
}

func (f *AtomFeed) WriteOut(w io.Writer) error {
	if f.Inner != nil {
		_, err := f.Inner.WriteTo(w)
		return err
	}
	return nil
}
