package grss

import (
	"encoding/xml"
	"github.com/hellodword/grss/pkg/etree"
)

// https://datatracker.ietf.org/doc/html/rfc4287
// https://validator.w3.org/feed/docs/atom.html
// https://support.google.com/merchants/answer/160598?hl=en
// https://wink.apache.org/1.4.0/api/index.html
// http://xml.coverpages.org/atom.html

// http://www.atomenabled.org/developers/syndication/

const (
	// AtomMime https://datatracker.ietf.org/doc/html/rfc4287#section-7
	AtomMime         = "application/atom+xml"
	AtomMimeFallback = "application/xml"
)

// AtomUri Unconstrained; it's not entirely clear how IRI fit into
//
//	xsd:anyURI so let's not try to constrain it here
//	AtomUri = text
type AtomUri string

// AtomCommonAttributes Common attributes
//
//	AtomCommonAttributes =
//	   attribute xml:base { AtomUri }?,
//	   attribute xml:lang { AtomLanguageTag }?,
//	   undefinedAttribute*
type AtomCommonAttributes struct {
	Base     AtomUri         `xml:"base,attr,omitempty"`
	Language AtomLanguageTag `xml:"lang,attr,omitempty"`
	// https://github.com/golang/go/issues/3633#issuecomment-328678522
	UndefinedAttribute []etree.Attr `xml:",any,attr,omitempty"`
}

// AtomLanguageTag As defined in RFC 3066
//
//	AtomLanguageTag = xsd:string {
//	   pattern = "[A-Za-z]{1,8}(-[A-Za-z0-9]{1,8})*"
//	}
type AtomLanguageTag string

// AtomEmailAddress Whatever an email address is, it contains at least one @
//
//	AtomEmailAddress = xsd:string { pattern = ".+@.+" }
type AtomEmailAddress string

// AtomTextConstruct = AtomPlainTextConstruct | AtomXHTMLTextConstruct
type AtomTextConstruct struct {
	AtomCommonAttributes
	Type string        `xml:"type,attr,omitempty"`
	Text string        `xml:",chardata"`
	Div  *AtomXhtmlDiv `xml:"div,omitempty"`
}

// AtomPlainTextConstruct =
//
//	AtomCommonAttributes,
//	attribute type { "text" | "html" }?,
//	text
type AtomPlainTextConstruct struct {
	AtomCommonAttributes
	Type string `xml:"type,attr,omitempty"`
	Text string `xml:",chardata"`
}

// AtomXHTMLTextConstruct =
//
//	AtomCommonAttributes,
//	attribute type { "xhtml" },
//	AtomXhtmlDiv
type AtomXHTMLTextConstruct struct {
	AtomCommonAttributes
	Type string        `xml:"type,attr,omitempty"`
	Div  *AtomXhtmlDiv `xml:"div,omitempty"`
}

//	AtomXhtmlDiv anyXHTML = element xhtml:* {
//	     (attribute * { text }
//	      | text
//	      | anyXHTML)*
//	  }
//
//	  xhtmlDiv = element xhtml:div {
//	     (attribute * { text }
//	      | text
//	      | anyXHTML)*
//	  }
type AtomXhtmlDiv struct {
	UndefinedAttribute []xml.Attr `xml:",any,attr,omitempty"`
	// TODO text | anyXHTML
	Text []byte `xml:",innerxml"`
}

// AtomPersonConstruct Person Construct
//
//	AtomPersonConstruct =
//	   AtomCommonAttributes,
//	   (element atom:name { text }
//	    & element atom:uri { AtomUri }?
//	    & element atom:email { AtomEmailAddress }?
//	    & ExtensionElement*)
type AtomPersonConstruct struct {
	AtomCommonAttributes
	Name  string           `xml:"name,omitempty"`
	Uri   AtomUri          `xml:"uri,omitempty"`
	Email AtomEmailAddress `xml:"email,omitempty"`
}

// AtomFeed atom:feed
//
//	atomFeed =
//	   [
//	      s:rule [
//	         context = "atom:feed"
//	         s:assert [
//	            test = "atom:author or not(atom:entry[not(atom:author)])"
//	            "An atom:feed must have an atom:author unless all "
//	            ~ "of its atom:entry children have an atom:author."
//	         ]
//	      ]
//	   ]
//	   element atom:feed {
//	      AtomCommonAttributes,
//	      (atomAuthor*
//	       & AtomCategory*
//	       & atomContributor*
//	       & AtomGenerator?
//	       & AtomIcon?
//	       & AtomId
//	       & AtomLink*
//	       & AtomLogo?
//	       & atomRights?
//	       & atomSubtitle?
//	       & atomTitle
//	       & atomUpdated
//	       & ExtensionElement*),
//	      AtomEntry*
//	   }
type AtomFeed struct {
	Inner *etree.Document `xml:"-"`

	XMLName xml.Name `xml:"feed"`

	AtomCommonAttributes

	Authors      []*AtomPersonConstruct `xml:"author"`
	Categories   []*AtomCategory        `xml:"category,omitempty"`
	Contributors []*AtomPersonConstruct `xml:"contributor,omitempty"`
	Generator    *AtomGenerator         `xml:"generator,omitempty"`
	Icon         *AtomUriExt            `xml:"icon,omitempty"`
	ID           AtomUriExt             `xml:"id"`
	Links        []*AtomLink            `xml:"link,omitempty"`
	Logo         *AtomUriExt            `xml:"logo,omitempty"`
	Rights       *AtomTextConstruct     `xml:"rights,omitempty"`
	Subtitle     *AtomTextConstruct     `xml:"subtitle,omitempty"`
	Title        *AtomTextConstruct     `xml:"title"`
	Updated      *AtomDateConstruct     `xml:"updated,omitempty"`

	ExtensionElements []XmlGeneric `xml:",any"`

	Entries []*AtomEntry `xml:"entry"`
}

// ExtensionElement = simpleExtensionElement | structuredExtensionElement
// simpleExtensionElement =
//
//	   element * - atom:* {
//	      text
//	   }
//
//	# Structured Extension
//
//	structuredExtensionElement =
//	   element * - atom:* {
//	      (attribute * { text }+,
//	         (text|anyElement)*)
//	    | (attribute * { text }*,
//	       (text?, anyElement+, (text|anyElement)*))
//	   }
//
// TODO ExtensionElement

// UndefinedContent undefinedContent = (text|anyForeignElement)*
//
//	 anyElement =
//	   element * {
//	      (attribute * { text }
//	       | text
//	       | anyElement)*
//	   }
//
//	anyForeignElement =
//	   element * - atom:* {
//	      (attribute * { text }
//	       | text
//	       | anyElement)*
//	   }
//
// TODO UndefinedContent
type UndefinedContent struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr,omitempty"`
	Content    string     `xml:",innerxml"`
}

// AtomCategory The "atom:category" element conveys information about a category
//
//	associated with an entry or feed.  This specification assigns no
//	meaning to the content (if any) of this element.
//	atomCategory =
//	   element atom:category {
//	      AtomCommonAttributes,
//	      attribute term { text },
//	      attribute scheme { AtomUri }?,
//	      attribute label { text }?,
//	      undefinedContent
//	   }
type AtomCategory struct {
	AtomCommonAttributes
	Term   string  `xml:"term,attr,omitempty"`
	Scheme AtomUri `xml:"scheme,attr,omitempty"`
	Label  string  `xml:"label,attr,omitempty"`

	//UndefinedContent []UndefinedContent `xml:",any,omitempty"`
}

// AtomGenerator The "atom:generator" element's content identifies the agent used to
//
//	generate a feed, for debugging and other purposes.
//	AtomGenerator = element atom:generator {
//	   AtomCommonAttributes,
//	   attribute uri { AtomUri }?,
//	   attribute version { text }?,
//	   text
//	}
type AtomGenerator struct {
	AtomCommonAttributes
	URI     AtomUri `xml:"uri,attr,omitempty"`
	Version string  `xml:"version,attr,omitempty"`
	Text    string  `xml:",chardata"`
}

type AtomUriExt struct {
	AtomCommonAttributes
	AtomUri `xml:",chardata"`
}

//// AtomIcon atom:icon
////
////	AtomIcon = element atom:icon {
////	   AtomCommonAttributes,
////	   (AtomUri)
////	}
//type AtomIcon AtomUriExt
//
//// AtomId atom:id
////
////	AtomId = element atom:id {
////	   AtomCommonAttributes,
////	   (AtomUri)
////	}
//type AtomId AtomUriExt
//
//// AtomLogo atom:logo
////
////	AtomLogo = element atom:logo {
////	   AtomCommonAttributes,
////	   (AtomUri)
////	}
//type AtomLogo AtomUriExt

// AtomLink The "atom:link" element defines a reference from an entry or feed to
//
//	a Web resource.  This specification assigns no meaning to the content
//	(if any) of this element.
//	AtomLink =
//	   element atom:link {
//	      AtomCommonAttributes,
//	      attribute href { AtomUri },
//	      attribute rel { AtomNCName | AtomUri }?,
//	      attribute type { AtomMediaType }?,
//	      attribute hreflang { AtomLanguageTag }?,
//	      attribute title { text }?,
//	      attribute length { text }?,
//	      undefinedContent
//	   }
type AtomLink struct {
	AtomCommonAttributes
	Href AtomUri `xml:"href,attr"`
	// TODO rel { atomNCName | atomUri }
	Rel      string          `xml:"rel,attr,omitempty"`
	Type     AtomMediaType   `xml:"type,attr,omitempty"`
	Hreflang AtomLanguageTag `xml:"hreflang,attr,omitempty"`
	Title    string          `xml:"title,attr,omitempty"`
	Length   string          `xml:"length,attr,omitempty"`

	//UndefinedContent []UndefinedContent `xml:",any"`
}

// AtomMediaType Whatever a media type is, it contains at least one slash
//
//	AtomMediaType = xsd:string { pattern = ".+/.+" }
type AtomMediaType string

// AtomNCName AtomNCName = xsd:string { minLength = "1" pattern = "[^:]*" }
type AtomNCName string

// AtomDateConstruct Date Construct
//
//	AtomDateConstruct =
//	   AtomCommonAttributes,
//	   xsd:dateTime
type AtomDateConstruct struct {
	AtomCommonAttributes
	// TODO xsd:dateTime
	DateTime string `xml:",chardata"`
}

// AtomEntry atom:entry
//
//	AtomEntry =
//	   [
//	      s:rule [
//	         context = "atom:entry"
//	         s:assert [
//	            test = "atom:link[@rel='alternate'] "
//	            ~ "or atom:link[not(@rel)] "
//	            ~ "or atom:content"
//	            "An atom:entry must have at least one atom:link element "
//	            ~ "with a rel attribute of 'alternate' "
//	            ~ "or an atom:content."
//	         ]
//	      ]
//	      s:rule [
//	         context = "atom:entry"
//	         s:assert [
//	            test = "atom:author or "
//	            ~ "../atom:author or atom:source/atom:author"
//	            "An atom:entry must have an atom:author "
//	            ~ "if its feed does not."
//	         ]
//	      ]
//	   ]
//	   element atom:entry {
//	      AtomCommonAttributes,
//	      (AtomAuthor*
//	       & AtomCategory*
//	       & AtomContent?
//	       & atomContributor*
//	       & AtomId
//	       & AtomLink*
//	       & AtomPublished?
//	       & atomRights?
//	       & AtomSource?
//	       & atomSummary?
//	       & atomTitle
//	       & atomUpdated
//	       & ExtensionElement*)
//	   }
type AtomEntry struct {
	AtomCommonAttributes

	Authors      []*AtomPersonConstruct `xml:"author"`
	Categories   []*AtomCategory        `xml:"category,omitempty"`
	Content      *AtomContent           `xml:"content,omitempty"`
	Contributors []*AtomPersonConstruct `xml:"contributor,omitempty"`
	ID           AtomUriExt             `xml:"id"`
	Links        []*AtomLink            `xml:"link,omitempty"`
	Published    *AtomDateConstruct     `xml:"published,omitempty"`
	Rights       *AtomTextConstruct     `xml:"rights,omitempty"`
	Source       *AtomSource            `xml:"source,omitempty"`
	Summary      *AtomTextConstruct     `xml:"summary,omitempty"`
	Title        AtomTextConstruct      `xml:"title"`
	Updated      *AtomDateConstruct     `xml:"updated,omitempty"`

	ExtensionElements []XmlGeneric `xml:",any"`
}

// AtomSource atom:source
//
//	AtomSource =
//	   element atom:source {
//	      AtomCommonAttributes,
//	      (atomAuthor*
//	       & AtomCategory*
//	       & atomContributor*
//	       & AtomGenerator?
//	       & AtomIcon?
//	       & AtomId?
//	       & AtomLink*
//	       & AtomLogo?
//	       & atomRights?
//	       & atomSubtitle?
//	       & atomTitle?
//	       & atomUpdated?
//	       & ExtensionElement*)
//	   }
type AtomSource struct {
	AtomCommonAttributes

	Authors      []*AtomPersonConstruct `xml:"author"`
	Categories   []*AtomCategory        `xml:"category,omitempty"`
	Contributors []*AtomPersonConstruct `xml:"contributor,omitempty"`
	Generator    *AtomGenerator         `xml:"generator,omitempty"`
	Icon         *AtomUriExt            `xml:"icon,omitempty"`
	ID           AtomUriExt             `xml:"id"`
	Links        []*AtomLink            `xml:"link,omitempty"`
	Logo         *AtomUriExt            `xml:"logo,omitempty"`
	Rights       *AtomTextConstruct     `xml:"rights,omitempty"`
	Subtitle     *AtomTextConstruct     `xml:"subtitle,omitempty"`
	Title        *AtomTextConstruct     `xml:"title"`
	Updated      *AtomDateConstruct     `xml:"updated,omitempty"`

	ExtensionElements []XmlGeneric `xml:",any"`
}

// AtomContent atom:content
//
//	    atomInlineTextContent =
//		   element atom:content {
//		      AtomCommonAttributes,
//		      attribute type { "text" | "html" }?,
//		      (text)*
//		   }
//
//		atomInlineXHTMLContent =
//		   element atom:content {
//		      AtomCommonAttributes,
//		      attribute type { "xhtml" },
//		      AtomXhtmlDiv
//		   }
//
//		atomInlineOtherContent =
//		   element atom:content {
//		      AtomCommonAttributes,
//		      attribute type { AtomMediaType }?,
//		      (text|anyElement)*
//		   }
//
//		atomOutOfLineContent =
//		   element atom:content {
//		      AtomCommonAttributes,
//		      attribute type { AtomMediaType }?,
//		      attribute src { AtomUri },
//		      empty
//		   }
//
//		AtomContent = atomInlineTextContent
//		 | atomInlineXHTMLContent
//		 | atomInlineOtherContent
//		 | atomOutOfLineContent
type AtomContent struct {
	AtomCommonAttributes
	Type  string        `xml:"type,attr,omitempty"`
	Src   *AtomUri      `xml:"src,attr,omitempty"`
	Text  string        `xml:",chardata"`
	Div   *AtomXhtmlDiv `xml:"div,omitempty"`
	Bytes []byte        `xml:",innerxml"`
}
