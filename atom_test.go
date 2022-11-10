package grss

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_AtomLanguageTag(t *testing.T) {
	s := `<AtomLanguageTag>zh-CN</AtomLanguageTag>`
	var a AtomLanguageTag
	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)
}

func Test_AtomLanguageTag_XsdStringNotMatchingPattern(t *testing.T) {
	s := `<AtomLanguageTag>zh99999999-CN</AtomLanguageTag>`
	var a AtomLanguageTag
	err := xml.Unmarshal([]byte(s), &a)
	assert.ErrorIs(t, err, ErrXsdStringNotMatchingPattern)
}

func Test_AtomEmailAddress(t *testing.T) {
	s := `<AtomEmailAddress>admin@exmaple.org</AtomEmailAddress>`
	var a AtomEmailAddress
	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)
}

func Test_AtomEmailAddress_XsdStringNotMatchingPattern(t *testing.T) {
	s := `<AtomEmailAddress>admin#exmaple.org</AtomEmailAddress>`
	var a AtomEmailAddress
	err := xml.Unmarshal([]byte(s), &a)
	assert.ErrorIs(t, err, ErrXsdStringNotMatchingPattern)
}

func Test_AtomTextConstruct_Text(t *testing.T) {
	s := `<title type="text" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">title</title>`

	type title struct {
		XMLName xml.Name `xml:"title"`
		AtomTextConstruct
	}
	var a title

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "text", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.UndefinedAttribute[0].Value, "http://www.w3.org/2005/Atom", a)
	assert.Equal(t, a.Text, "title", a)

}

func Test_AtomTextConstruct_HTML(t *testing.T) {
	s := `<title type="html" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">&lt;h1&gt;title&lt;/h1&gt;</title>`

	type title struct {
		XMLName xml.Name `xml:"title"`
		AtomTextConstruct
	}
	var a title

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "html", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.UndefinedAttribute[0].Value, "http://www.w3.org/2005/Atom", a)
	assert.Equal(t, a.Text, "<h1>title</h1>", a)

}

func Test_AtomTextConstruct_XHTML(t *testing.T) {
	s := `
<title type="xhtml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <div xmlns="http://www.w3.org/1999/xhtml">
        <h1>title</h1>
    </div>
</title>
`

	type title struct {
		XMLName xml.Name `xml:"title"`
		AtomTextConstruct
	}
	var a title

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "xhtml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Div.UndefinedAttribute[0].Value, "http://www.w3.org/1999/xhtml", a)
	assert.Equal(t, string(a.Div.Text), "\n        <h1>title</h1>\n    ", a)

}

func Test_AtomTextConstruct_XHTMLWithText(t *testing.T) {
	s := `
<title type="xhtml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <div xmlns="http://www.w3.org/1999/xhtml">
        title
    </div>
</title>
`

	type title struct {
		XMLName xml.Name `xml:"title"`
		AtomTextConstruct
	}
	var a title

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "xhtml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Div.UndefinedAttribute[0].Value, "http://www.w3.org/1999/xhtml", a)
	assert.Equal(t, string(a.Div.Text), "\n        title\n    ", a)

}

func Test_AtomTextConstruct_XHTMLWithDiv(t *testing.T) {
	// TODO div could be parsed
	s := `
<title type="xhtml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <div xmlns="http://www.w3.org/1999/xhtml">
        <h1>
            <div>title</div>
        </h1>
    </div>
</title>
`

	type title struct {
		XMLName xml.Name `xml:"title"`
		AtomTextConstruct
	}
	var a title

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "xhtml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Div.UndefinedAttribute[0].Value, "http://www.w3.org/1999/xhtml", a)
	assert.Equal(t, string(a.Div.Text), "\n        <h1>\n            <div>title</div>\n        </h1>\n    ", a)

}

func Test_AtomContent_Text(t *testing.T) {
	s := `<content type="text" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">title</content>`

	type content struct {
		XMLName xml.Name `xml:"content"`
		AtomContent
	}
	var a content

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "text", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Text, "title", a)

}

func Test_AtomContent_HTML(t *testing.T) {
	s := `<content type="html" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">&lt;h1&gt;title&lt;/h1&gt;</content>`

	type content struct {
		XMLName xml.Name `xml:"content"`
		AtomContent
	}
	var a content

	err := xml.Unmarshal([]byte(s), &a)

	assert.Nil(t, err)

	assert.Equal(t, a.Type, "html", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.UndefinedAttribute[0].Value, "http://www.w3.org/2005/Atom", a)
	assert.Equal(t, a.Text, "<h1>title</h1>", a)

}

func Test_AtomContent_XHTML(t *testing.T) {
	s := `
<content type="xhtml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <div xmlns="http://www.w3.org/1999/xhtml">
        <h1>title</h1>
    </div>
</content>
`

	type content struct {
		XMLName xml.Name `xml:"content"`
		AtomContent
	}
	var a content

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "xhtml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Div.UndefinedAttribute[0].Value, "http://www.w3.org/1999/xhtml", a)
	assert.Equal(t, string(a.Div.Text), "\n        <h1>title</h1>\n    ", a)

}

func Test_AtomContent_XHTMLWithText(t *testing.T) {
	s := `
<content type="xhtml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <div xmlns="http://www.w3.org/1999/xhtml">
        title
    </div>
</content>
`

	type content struct {
		XMLName xml.Name `xml:"content"`
		AtomContent
	}
	var a content

	err := xml.Unmarshal([]byte(s), &a)

	assert.Nil(t, err)

	assert.Equal(t, a.Type, "xhtml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, a.Div.UndefinedAttribute[0].Value, "http://www.w3.org/1999/xhtml", a)
	assert.Equal(t, string(a.Div.Text), "\n        title\n    ", a)
}

func Test_AtomContent_XML(t *testing.T) {
	s := `
<content type="application/xml" xml:lang="en-us" xml:base="http://title/base" xmlns="http://www.w3.org/2005/Atom">
    <x xmlns="http://x/">title</x>
</content>
`

	type content struct {
		XMLName xml.Name `xml:"content"`
		AtomContent
	}
	var a content

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Type, "application/xml", a)
	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.Base, "http://title/base", a)
	assert.Equal(t, string(a.Bytes), "\n    <x xmlns=\"http://x/\">title</x>\n", a)
}

func Test_AtomEntry_001(t *testing.T) {
	s := `
<entry xml:lang="en-us" xml:base="http://entry/base" anyAttr="anyAttrValue" xmlns="http://www.w3.org/2005/Atom">
    <id>1</id>
    <updated>2022-11-08T20:48:11Z</updated>
    <title type="text">title</title>
    <summary type="text">summary</summary>
    <published>2022-11-07T22:54:20Z</published>
    <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
    <author>
        <email>author@hp.com</email>
        <name>author</name>
        <uri>http://uri</uri>
    </author>
    <contributor>
        <email>cont@hp.com</email>
        <name>cont</name>
        <uri>http://uri</uri>
    </contributor>
    <category label="label" scheme="scheme" term="term"/>
</entry>
`
	type entry struct {
		XMLName xml.Name `xml:"entry"`
		AtomEntry
	}
	var a entry

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.UndefinedAttribute[0].Name.Local, "anyAttr", a)
	assert.Equal(t, a.ID.AtomUri, "1", a)
	assert.Equal(t, a.Updated.DateTime, "2022-11-08T20:48:11Z", a)
	assert.Equal(t, a.Title.Text, "title", a)
	assert.Equal(t, a.Summary.Text, "summary", a)
	assert.Equal(t, a.Published.DateTime, "2022-11-07T22:54:20Z", a)
	assert.Equal(t, a.Links[0].Href, "href", a)
	assert.Equal(t, a.Authors[0].Email, "author@hp.com", a)
	assert.Equal(t, a.Contributors[0].Uri, "http://uri", a)
	assert.Equal(t, a.Categories[0].Term, "term", a)

}

func Test_AtomEntry_002(t *testing.T) {
	s := `
<entry xml:lang="en-us" xml:base="http://entry/base" anyAttr="anyAttrValue" xmlns="http://www.w3.org/2005/Atom">
    <id>2</id>
    <updated>2022-11-08T20:48:11Z</updated>
    <title type="text">title</title>
    <summary type="text">summary</summary>
    <published>2022-11-07T22:54:20Z</published>
    <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
    <author>
        <email>author@hp.com</email>
        <name>author</name>
        <uri>http://uri</uri>
    </author>
    <contributor>
        <email>cont@hp.com</email>
        <name>cont</name>
        <uri>http://uri</uri>
    </contributor>
    <category label="label" scheme="scheme" term="term"/>
    <content type="text/plain">Gustaf's Knäckebröd</content>
</entry>
`
	type entry struct {
		XMLName xml.Name `xml:"entry"`
		AtomEntry
	}
	var a entry

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.UndefinedAttribute[0].Value, "anyAttrValue", a)
	assert.Equal(t, a.ID.AtomUri, "2", a)
	assert.Equal(t, a.Updated.DateTime, "2022-11-08T20:48:11Z", a)
	assert.Equal(t, a.Title.Text, "title", a)
	assert.Equal(t, a.Summary.Text, "summary", a)
	assert.Equal(t, a.Published.DateTime, "2022-11-07T22:54:20Z", a)
	assert.Equal(t, a.Links[0].Href, "href", a)
	assert.Equal(t, a.Authors[0].Email, "author@hp.com", a)
	assert.Equal(t, a.Contributors[0].Uri, "http://uri", a)
	assert.Equal(t, a.Categories[0].Term, "term", a)
	assert.Equal(t, a.Content.Type, "text/plain", a)
	assert.Equal(t, a.Content.Text, "Gustaf's Knäckebröd", a)

}

func Test_AtomEntry_003(t *testing.T) {
	s := `
<?xml version="1.0" encoding="UTF-8"?>
<entry xmlns="http://www.w3.org/2005/Atom" xmlns:ns2="http://a9.com/-/spec/opensearch/1.1/" xmlns:ns3="http://www.w3.org/1999/xhtml" anyAttr="anyAttrValue" xml:base="http://entry/base" xml:lang="en-us">
    <id>3</id>
    <updated>1970-01-01T02:20:34.567+02:00</updated>
    <title type="text">title</title>
    <summary type="text">summary</summary>
    <published>1970-01-01T02:20:34.567+02:00</published>
    <link href="href" hreflang="en-us" length="10" rel="rel" title="title" type="text/plain"/>
    <author>
        <email>author@hp.com</email>
        <name>author</name>
        <uri>http://uri</uri>
    </author>
    <contributor>
        <email>cont@hp.com</email>
        <name>cont</name>
        <uri>http://uri</uri>
    </contributor>
    <category label="label" scheme="scheme" term="term"/>
    <content type="application/xml">
        <x:x xmlns="http://x/" xmlns:ns2="http://www.w3.org/2005/Atom" xmlns:ns3="http://a9.com/-/spec/opensearch/1.1/" xmlns:ns4="http://www.w3.org/1999/xhtml" xmlns:x="http://x/">Gustaf's Knäckebröd</x:x>
    </content>
</entry>
`
	type entry struct {
		XMLName xml.Name `xml:"entry"`
		AtomEntry
	}
	var a entry

	err := xml.Unmarshal([]byte(s), &a)
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.UndefinedAttribute[1].Name.Local, "ns2", a)
	assert.Equal(t, a.ID.AtomUri, "3", a)
	assert.Equal(t, a.Updated.DateTime, "1970-01-01T02:20:34.567+02:00", a)
	assert.Equal(t, a.Title.Text, "title", a)
	assert.Equal(t, a.Summary.Text, "summary", a)
	assert.Equal(t, a.Published.DateTime, "1970-01-01T02:20:34.567+02:00", a)
	assert.Equal(t, a.Links[0].Href, "href", a)
	assert.Equal(t, a.Authors[0].Email, "author@hp.com", a)
	assert.Equal(t, a.Contributors[0].Uri, "http://uri", a)
	assert.Equal(t, a.Categories[0].Term, "term", a)
	assert.Equal(t, a.Content.Type, "application/xml", a)
	assert.Equal(t, string(a.Content.Bytes), "\n        <x:x xmlns=\"http://x/\" xmlns:ns2=\"http://www.w3.org/2005/Atom\" xmlns:ns3=\"http://a9.com/-/spec/opensearch/1.1/\" xmlns:ns4=\"http://www.w3.org/1999/xhtml\" xmlns:x=\"http://x/\">Gustaf's Knäckebröd</x:x>\n    ", a.Content)

}

func Test_AtomFeed_001(t *testing.T) {
	s := `
<feed xml:lang="en-us" xml:base="http://feed/base" anyAttr="anyAttrValue" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/" xmlns="http://www.w3.org/2005/Atom">
    <id>id</id>
    <updated>2022-11-08T20:48:11Z</updated>
    <title type="text">title</title>
    <subtitle type="text">subtitle</subtitle>
    <opensearch:itemsPerPage>5</opensearch:itemsPerPage>
    <opensearch:startIndex>6</opensearch:startIndex>
    <opensearch:totalResults>7</opensearch:totalResults>
    <opensearch:Query searchTerms="query 1"/>
    <opensearch:Query searchTerms="query 2"/>
    <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
    <author>
        <email>author@hp.com</email>
        <name>author</name>
        <uri>http://uri</uri>
    </author>
    <contributor>
        <email>cont@hp.com</email>
        <name>cont</name>
        <uri>http://uri</uri>
    </contributor>
    <category label="label" scheme="scheme" term="term"/>
    <generator version="1.0" uri="http://generator/uri" xml:lang="en-us" xml:base="http://generator/base">wink</generator>
    <icon>icon</icon>
    <logo>logo</logo>
    <rights type="text">rights</rights>
    <entry xml:lang="en-us" xml:base="http://entry/base" anyAttr="anyAttrValue">
        <id>1</id>
        <updated>2022-11-07T20:48:11Z</updated>
        <title type="text">title</title>
        <summary type="text">summary</summary>
        <published>2022-11-06T20:48:11Z</published>
        <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
        <author>
            <email>author@hp.com</email>
            <name>author</name>
            <uri>http://uri</uri>
        </author>
        <contributor>
            <email>cont@hp.com</email>
            <name>cont</name>
            <uri>http://uri</uri>
        </contributor>
        <category label="label" scheme="scheme" term="term"/>
        <content type="application/xml">
            <x:x xmlns="http://x/" xmlns:ns2="http://www.w3.org/2005/Atom" xmlns:ns3="http://a9.com/-/spec/opensearch/1.1/" xmlns:ns4="http://www.w3.org/1999/xhtml" xmlns:x="http://x/">Gustaf's Knäckebröd</x:x>
        </content>
    </entry>
</feed>
`

	a, err := AtomParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.UndefinedAttribute[0].Name.Local, "anyAttr", a)
	assert.Equal(t, a.ID.AtomUri, "id", a)
	assert.Equal(t, a.Updated.DateTime, "2022-11-08T20:48:11Z", a)
	assert.Equal(t, a.Title.Text, "title", a)
	assert.Equal(t, a.Subtitle.Text, "subtitle", a)
	assert.Equal(t, a.ExtensionElement[0].Content, "5", a.ExtensionElement)
	assert.Equal(t, a.ExtensionElement[1].XMLName.Local, "opensearch:startIndex", a.ExtensionElement)
	assert.Equal(t, a.Links[0].Href, "href", a)
	assert.Equal(t, a.Authors[0].Email, "author@hp.com", a)
	assert.Equal(t, a.Contributors[0].Uri, "http://uri", a)
	assert.Equal(t, a.Categories[0].Term, "term", a)
	assert.Equal(t, a.Icon.AtomUri, "icon", a)
	assert.Equal(t, a.Logo.AtomUri, "logo", a)
	assert.Equal(t, a.Rights.Text, "rights", a)
	assert.Equal(t, a.Entries[0].Content.Type, "application/xml", a)
	assert.Equal(t, string(a.Entries[0].Content.Bytes), "\n            <x:x xmlns=\"http://x/\" xmlns:ns2=\"http://www.w3.org/2005/Atom\" xmlns:ns3=\"http://a9.com/-/spec/opensearch/1.1/\" xmlns:ns4=\"http://www.w3.org/1999/xhtml\" xmlns:x=\"http://x/\">Gustaf's Knäckebröd</x:x>\n        ", a.Entries[0])

}

func Test_AtomFeed_002(t *testing.T) {
	s := `
<feed xml:lang="en-us" xml:base="http://feed/base" anyAttr="anyAttrValue" xmlns="http://www.w3.org/2005/Atom">
    <id>id</id>
    <updated>2022-11-08T20:48:11Z</updated>
    <title type="text">title</title>
    <subtitle type="text">subtitle</subtitle>
    <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
    <author>
        <email>author@hp.com</email>
        <name>author</name>
        <uri>http://uri</uri>
    </author>
    <contributor>
        <email>cont@hp.com</email>
        <name>cont</name>
        <uri>http://uri</uri>
    </contributor>
    <category label="label" scheme="scheme" term="term"/>
    <generator version="1.0" uri="http://generator/uri" xml:lang="en-us" xml:base="http://generator/base">wink</generator>
    <icon>icon</icon>
    <logo>logo</logo>
    <rights type="text">rights</rights>
    <entry xml:lang="en-us" xml:base="http://entry/base" anyAttr="anyAttrValue">
        <id>1</id>
        <updated>2022-11-07T20:48:11Z</updated>
        <title type="text">title</title>
        <summary type="text">summary</summary>
        <published>2022-11-06T20:48:11Z</published>
        <link href="href" type="text/plain" rel="rel" hreflang="en-us" title="title" length="10"/>
        <author>
            <email>author@hp.com</email>
            <name>author</name>
            <uri>http://uri</uri>
        </author>
        <contributor>
            <email>cont@hp.com</email>
            <name>cont</name>
            <uri>http://uri</uri>
        </contributor>
        <category label="label" scheme="scheme" term="term"/>
        <content type="application/xml">
            <x1 xmlns="xxx" xmlns:y="yyy">
                <x2>
                    <y:y1>Gustaf's Knäckebröd</y:y1>
                </x2>
            </x1>
        </content>
    </entry>
</feed>
`

	a, err := AtomParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-us", a)
	assert.Equal(t, a.UndefinedAttribute[0].Name.Local, "anyAttr", a)
	assert.Equal(t, a.ID.AtomUri, "id", a)
	assert.Equal(t, a.Updated.DateTime, "2022-11-08T20:48:11Z", a)
	assert.Equal(t, a.Title.Text, "title", a)
	assert.Equal(t, a.Subtitle.Text, "subtitle", a)
	assert.Equal(t, a.Links[0].Href, "href", a)
	assert.Equal(t, a.Authors[0].Email, "author@hp.com", a)
	assert.Equal(t, a.Contributors[0].Uri, "http://uri", a)
	assert.Equal(t, a.Categories[0].Term, "term", a)
	assert.Equal(t, a.Icon.AtomUri, "icon", a)
	assert.Equal(t, a.Logo.AtomUri, "logo", a)
	assert.Equal(t, a.Rights.Text, "rights", a)
	assert.Equal(t, a.Entries[0].Content.Type, "application/xml", a)
	assert.Equal(t, string(a.Entries[0].Content.Bytes), "\n            <x1 xmlns=\"xxx\" xmlns:y=\"yyy\">\n                <x2>\n                    <y:y1>Gustaf's Knäckebröd</y:y1>\n                </x2>\n            </x1>\n        ", a.Entries[0])

}

func Test_AtomFeed_003(t *testing.T) {
	s := `
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>tag:github.com,2008:/torvalds</id>
  <link type="text/html" rel="alternate" href="https://github.com/torvalds"/>
  <link type="application/atom+xml" rel="self" href="https://github.com/torvalds.atom"/>
  <title>GitHub Public Timeline Feed</title>
  <updated>2022-11-08T20:48:11Z</updated>
  <entry>
    <id>tag:github.com,2008:PushEvent/25111313160</id>
    <published>2022-11-08T20:48:11Z</published>
    <updated>2022-11-08T20:48:11Z</updated>
    <link type="text/html" rel="alternate" href="https://github.com/torvalds/linux/compare/59f2f4b8a7...f141df3713"/>
    <title type="html">torvalds pushed to master in torvalds/linux</title>
    <author>
      <name>torvalds</name>
      <uri>https://github.com/torvalds</uri>
    </author>
    <media:thumbnail height="30" width="30" url="https://avatars.githubusercontent.com/u/1024025?s=30&amp;v=4"/>
    <content type="html">&lt;div class=&quot;push js-feed-item-view&quot; data-hydro-view=&quot;{&amp;quot;event_type&amp;quot;:&amp;quot;news_feed.event.view&amp;quot;,&amp;quot;payload&amp;quot;:{&amp;quot;event&amp;quot;:{&amp;quot;repo_id&amp;quot;:2325298,&amp;quot;actor_id&amp;quot;:1024025,&amp;quot;public&amp;quot;:true,&amp;quot;type&amp;quot;:&amp;quot;PushEvent&amp;quot;,&amp;quot;target_id&amp;quot;:null,&amp;quot;id&amp;quot;:25111313160,&amp;quot;additional_details_shown&amp;quot;:false,&amp;quot;grouped&amp;quot;:false},&amp;quot;event_group&amp;quot;:null,&amp;quot;org_id&amp;quot;:null,&amp;quot;target_type&amp;quot;:&amp;quot;event&amp;quot;,&amp;quot;user_id&amp;quot;:null,&amp;quot;feed_card&amp;quot;:{&amp;quot;card_retrieved_id&amp;quot;:&amp;quot;688e5ecf-1153-441c-b9ad-f197c3ceb067&amp;quot;},&amp;quot;originating_url&amp;quot;:&amp;quot;https://github.com/torvalds.atom&amp;quot;}}&quot; data-hydro-view-hmac=&quot;fdd03fb524ea3f1eea1ce30fc90ec1c0564728b329bc36c74efc63faba15beb1&quot;&gt;&lt;div class=&quot;body&quot;&gt;
&lt;!-- push --&gt;
&lt;div class=&quot;d-flex flex-items-baseline border-bottom color-border-muted py-3&quot;&gt;
    &lt;span class=&quot;mr-2&quot;&gt;&lt;a class=&quot;d-inline-block&quot; href=&quot;/torvalds&quot; rel=&quot;noreferrer&quot;&gt;&lt;img class=&quot;avatar avatar-user&quot; src=&quot;https://avatars.githubusercontent.com/u/1024025?s=64&amp;amp;v=4&quot; width=&quot;32&quot; height=&quot;32&quot; alt=&quot;@torvalds&quot;&gt;&lt;/a&gt;&lt;/span&gt;
  &lt;div class=&quot;d-flex flex-column width-full&quot;&gt;
    &lt;div class=&quot;&quot;&gt;
      &lt;a class=&quot;Link--primary no-underline wb-break-all text-bold d-inline-block&quot; href=&quot;/torvalds&quot; rel=&quot;noreferrer&quot;&gt;torvalds&lt;/a&gt;
      
      pushed to
        &lt;a class=&quot;branch-name&quot; href=&quot;/torvalds/linux/tree/master&quot; rel=&quot;noreferrer&quot;&gt;master&lt;/a&gt;
        in
      &lt;a class=&quot;Link--primary no-underline wb-break-all text-bold d-inline-block&quot; href=&quot;/torvalds/linux&quot; rel=&quot;noreferrer&quot;&gt;torvalds/linux&lt;/a&gt;
        &lt;span class=&quot;color-fg-muted no-wrap f6 ml-1&quot;&gt;
          &lt;relative-time datetime=&quot;2022-11-08T20:48:11Z&quot; class=&quot;no-wrap&quot;&gt;Nov 8, 2022&lt;/relative-time&gt;
        &lt;/span&gt;

        &lt;div class=&quot;Box p-3 mt-2 &quot;&gt;
          &lt;span&gt;2 commits to&lt;/span&gt;
          &lt;a class=&quot;branch-name&quot; href=&quot;/torvalds/linux/tree/master&quot; rel=&quot;noreferrer&quot;&gt;master&lt;/a&gt;

          &lt;div class=&quot;commits pusher-is-only-committer&quot;&gt;
            &lt;ul class=&quot;list-style-none&quot;&gt;
                &lt;li class=&quot;d-flex flex-items-baseline&quot;&gt;
                  &lt;span title=&quot;torvalds&quot;&gt;
                    &lt;a class=&quot;d-inline-block&quot; href=&quot;/torvalds&quot; rel=&quot;noreferrer&quot;&gt;&lt;img class=&quot;mr-1 avatar-user&quot; src=&quot;https://avatars.githubusercontent.com/u/1024025?s=32&amp;amp;v=4&quot; width=&quot;16&quot; height=&quot;16&quot; alt=&quot;@torvalds&quot;&gt;&lt;/a&gt;
                  &lt;/span&gt;
                  &lt;code&gt;&lt;a class=&quot;mr-1&quot; href=&quot;/torvalds/linux/commit/f141df371335645ce29a87d9683a3f79fba7fd67&quot; rel=&quot;noreferrer&quot;&gt;f141df3&lt;/a&gt;&lt;/code&gt;
                  &lt;div class=&quot;dashboard-break-word lh-condensed&quot;&gt;
                    &lt;blockquote&gt;
                      Merge tag &#39;audit-pr-20221107&#39; of git://git.kernel.org/pub/scm/linux/k…
                    &lt;/blockquote&gt;
                  &lt;/div&gt;
                &lt;/li&gt;
                &lt;li class=&quot;d-flex flex-items-baseline&quot;&gt;
                  &lt;span title=&quot;torvalds&quot;&gt;
                    &lt;a class=&quot;d-inline-block&quot; href=&quot;/torvalds&quot; rel=&quot;noreferrer&quot;&gt;&lt;img class=&quot;mr-1 avatar-user&quot; src=&quot;https://avatars.githubusercontent.com/u/1024025?s=32&amp;amp;v=4&quot; width=&quot;16&quot; height=&quot;16&quot; alt=&quot;@torvalds&quot;&gt;&lt;/a&gt;
                  &lt;/span&gt;
                  &lt;code&gt;&lt;a class=&quot;mr-1&quot; href=&quot;/torvalds/linux/commit/f49b2d89fb10ef5fa5fa1993f648ec5daa884bef&quot; rel=&quot;noreferrer&quot;&gt;f49b2d8&lt;/a&gt;&lt;/code&gt;
                  &lt;div class=&quot;dashboard-break-word lh-condensed&quot;&gt;
                    &lt;blockquote&gt;
                      Merge tag &#39;lsm-pr-20221107&#39; of git://git.kernel.org/pub/scm/linux/ker…
                    &lt;/blockquote&gt;
                  &lt;/div&gt;
                &lt;/li&gt;


                &lt;li class=&quot;f6 mt-2&quot;&gt;
                  &lt;a class=&quot;Link--secondary&quot; href=&quot;/torvalds/linux/compare/59f2f4b8a7...f141df3713&quot; rel=&quot;noreferrer&quot;&gt;2 more commits »&lt;/a&gt;
                &lt;/li&gt;
            &lt;/ul&gt;
          &lt;/div&gt;
        &lt;/div&gt;
    &lt;/div&gt;
  &lt;/div&gt;
&lt;/div&gt;
&lt;/div&gt;&lt;/div&gt;</content>
  </entry>
</feed>
`

	a, err := AtomParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.Language, "en-US", a)
	assert.Equal(t, a.UndefinedAttribute[1].Name.Local, "xmlns:media", a)
	assert.Equal(t, a.Entries[0].ExtensionElement[0].XMLName.Local, "media:thumbnail", a)

}

func Test_AtomFeed_004(t *testing.T) {
	// https://support.google.com/merchants/answer/160598?hl=en
	s := `
<?xml version="1.0"?>
<feed version="0.3" xmlns="http://purl.org/atom/ns#"
xmlns:g="http://base.google.com/ns/1.0">
  <title>The name of your data feed</title>
  <link href="http://www.example.com" rel="alternate" type="text/html" />
  <modified>2005-10-11T18:30:02Z</modified>
  <author>
    <name>Google</name>
  </author>
  <id>tag:google.com,2005-10-15:/support/products</id>
  <entry>
    <title>Red wool sweater</title>
    <link href="http://www.example.com/item1-info-page.html" />
    <summary>Comfortable and soft, this sweater will keep you warm on those cold winter nights.</summary>
    <id>tag:google.com,2005-10-15:/support/products</id>
    <issued>2005-10-13T18:30:02Z</issued>
    <modified>2005-10-13T18:30:02Z</modified>
    <g:image_link>http://www.example.com/image1.jpg</g:image_link>
    <g:price>25</g:price>
    <g:condition>new</g:condition>
  </entry>
</feed>
`

	a, err := AtomParse(strings.NewReader(s))
	assert.Nil(t, err)

	assert.Equal(t, a.AtomCommonAttributes.UndefinedAttribute[0].Name.Local, "version", a)
	assert.Equal(t, a.Links[0].Rel, "alternate", a)
	assert.Equal(t, a.Entries[0].ExtensionElement[0].XMLName.Local, "issued", a.Entries[0].ExtensionElement)
}
