package grss

import "encoding/xml"

// 0.90 https://www.rssboard.org/rss-0-9-0
// 0.91(netscape) https://www.rssboard.org/rss-0-9-1-netscape
// 0.91 https://www.rssboard.org/rss-0-9-1
// 0.92 https://www.rssboard.org/rss-0-9-2
// 1.0 https://web.resource.org/rss/1.0/spec
// 2.0 https://www.rssboard.org/rss-2-0
// 2.0 https://validator.w3.org/feed/docs/rss2.html
// 2.0 https://cyber.harvard.edu/rss/rss.html
// 2.0.1 https://www.rssboard.org/rss-2-0-1

// https://www.w3schools.com/xml/xml_rss.asp

// MRSS 1.5.1 https://www.rssboard.org/media-rss
// https://en.wikipedia.org/wiki/Media_RSS

// http://www.rssboard.org/files/rss-2.0-sample.xml
// http://www.rssboard.org/files/sample-rss-2.xml

// https://www.rssboard.org/rss-history
// https://www.rssboard.org/rss-profile

const (
	// RssMime https://www.rssboard.org/rss-mime-type-application.txt
	RssMime         = "application/rss+xml"
	RssMimeFallback = "application/xml"
)

type RssFeed struct {
	Charset string `xml:"-"`

	XMLName xml.Name // `xml:"rss"`

	Attributes   []xml.Attr `xml:",any,attr,omitempty"`
	Version      string     `xml:"version,attr,omitempty"`
	XmlnsContent string     `xml:"xmlns:content,attr,omitempty"`

	Channel *RssChannel `xml:"channel,omitempty"`

	// 0.90
	Image     *RssImage     `xml:"image,omitempty"`
	Items     []*RssItem    `xml:"item,omitempty"`
	TextInput *RssTextInput `xml:"textinput,omitempty"`
}

type RssChannel struct {
	Attributes []xml.Attr `xml:",any,attr,omitempty"`

	// the required channel elements, each with a brief description, an example, and where available, a pointer to a more complete description.

	// Title The name of the channel. It's how people refer to your service. If you have an HTML website that contains the same information as your RSS file, the title of your channel should be the same as the title of your website.
	Title string `xml:"title"`
	// Link The URL to the HTML website corresponding to the channel.
	Link string `xml:"link"`
	// Description	Phrase or sentence describing the channel.
	Description string `xml:"description"`

	// optional channel elements.

	// Language the channel is written in. This allows aggregators to group all Italian language sites, for example, on a single page. A list of allowable values for this element, as provided by Netscape, is here. You may also use values defined by the W3C.
	Language string `xml:"language,omitempty"`
	// 	Copyright notice for content in the channel.
	Copyright string `xml:"copyright,omitempty"`
	// ManagingEditor Email address for person responsible for editorial content.
	ManagingEditor string `xml:"managingEditor,omitempty"`
	// WebMaster Email address for person responsible for technical issues relating to channel.
	WebMaster string `xml:"webMaster,omitempty"`
	// PubDate The publication date for the content in the channel. For example, the New York Times publishes on a daily basis, the publication date flips once every 24 hours. That's when the pubDate of the channel changes. All date-times in RSS conform to the Date and Time Specification of RFC 822, with the exception that the year may be expressed with two characters or four characters (four preferred).
	PubDate string `xml:"pubDate,omitempty"`
	// LastBuildDate The last time the content of the channel changed.
	LastBuildDate string `xml:"lastBuildDate,omitempty"`
	// Category Specify one or more categories that the channel belongs to. Follows the same rules as the <item>-level category element.
	Category []*RssCategory `xml:"category,omitempty"`
	// Generator A string indicating the program used to generate the channel.
	Generator string `xml:"generator,omitempty"`
	// Docs A URL that points to the documentation for the format used in the RSS file. It's probably a pointer to this page. It's for people who might stumble across an RSS file on a Web server 25 years from now and wonder what it is.
	Docs string `xml:"docs,omitempty"`
	// Cloud Allows processes to register with a cloud to be notified of updates to the channel, implementing a lightweight publish-subscribe protocol for RSS feeds
	Cloud *XmlGeneric `xml:"cloud,omitempty"`
	// Ttl stands for time to live. It's a number of minutes that indicates how long a channel can be cached before refreshing from the source.
	Ttl string `xml:"ttl,omitempty"`
	// Image Specifies a GIF, JPEG or PNG image that can be displayed with the channel.
	Image *RssImage `xml:"image,omitempty"`
	// Rating The PICS rating for the channel.
	Rating *XmlGeneric `xml:"rating,omitempty"`
	// TextInput Specifies a text input box that can be displayed with the channel.
	TextInput *RssTextInput `xml:"textinput,omitempty"`
	// SkipHours A hint for aggregators telling them which hours they can skip.
	SkipHours *RssSkipHours `xml:"skipHours,omitempty"`
	// SkipDays A hint for aggregators telling them which days they can skip.
	SkipDays *RssSkipDays `xml:"skipDays,omitempty"`

	Items []*RssItem `xml:"item,omitempty"`

	ExtensionElement []XmlGeneric `xml:",any"`
}

type RssItem struct {
	// Title of the item.
	Title string `xml:"title,omitempty"`
	// Link The URL of the item.
	Link string `xml:"link,omitempty"`
	// Description	The item synopsis.
	Description string `xml:"description,omitempty"`
	// Email address of the author of the item.
	Author *RssAuthor `xml:"author,omitempty"`
	// Category Includes the item in one or more categories.
	Category []*RssCategory `xml:"category,omitempty"`
	// Comments URL of a page for comments relating to the item.
	Comments string `xml:"comments,omitempty"`
	// Enclosure Describes a media object that is attached to the item.
	Enclosure *RssEnclosure `xml:"enclosure,omitempty"`
	// Guid A string that uniquely identifies the item.
	Guid *RssGuid `xml:"guid,omitempty"`
	// PubDate Indicates when the item was published.
	PubDate string `xml:"pubDate,omitempty"`
	// Source The RSS channel that the item came from.
	Source *RssSource `xml:"source,omitempty"`

	ContentUnmarshal *RssContent `xml:"encoded,omitempty"`
	Content          *RssContent `xml:"content:encoded,omitempty"`
}

type RssContent struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
}

type RssEnclosure struct {
	// Length Defines the length (in bytes) of the media file
	Length string `xml:"length,attr,omitempty"`
	// Type Defines the type of media file
	Type string `xml:"type,attr,omitempty"`
	// Url Defines the URL to the media file
	Url string `xml:"url,attr,omitempty"`
}

// RssGuid stands for globally unique identifier. It's a string that uniquely identifies the item. When present, an aggregator may choose to use this string to determine if an item is new. There are no rules for the syntax of a guid. Aggregators must view them as a string. It's up to the source of the feed to establish the uniqueness of the string.
// 2.0.1
type RssGuid struct {
	// IsPermaLink If the guid element has an attribute named isPermaLink with a value of true, the reader may assume that it is a permalink to the item, that is, a url that can be opened in a Web browser, that points to the full item described by the <item> element. IsPermaLink is optional, its default value is true. If its value is false, the guid may not be assumed to be a url, or a url to anything in particular.
	IsPermaLink string `xml:"isPermaLink,attr,omitempty"`
	Guid        string `xml:",chardata"`
}

// RssAuthor It's the email address of the author of the item. For newspapers and magazines syndicating via RSS, the author is the person who wrote the article that the <item> describes. For collaborative weblogs, the author of the item might be different from the managing editor or webmaster. For a weblog authored by a single individual it would make sense to omit the <author> element.
type RssAuthor struct {
	Email string `xml:",chardata"`
}

// RssSource Its value is the name of the RSS channel that the item came from, derived from its <title>. It has one required attribute, url, which links to the XMLization of the source. The purpose of this element is to propogate credit for links, to publicize the sources of news items. It's used in the post command in the Radio UserLand aggregator. It should be generated automatically when forwarding an item from an aggregator to a weblog authoring tool.
type RssSource struct {
	// Url Specifies the link to the source
	Url  string `xml:"url,attr,omitempty"`
	Text string `xml:",chardata"`
}

// RssCategory You may include as many category elements as you need to, for different domains, and to have an item cross-referenced in different parts of the same domain.
type RssCategory struct {
	// Domain a string that identifies a categorization taxonomy.
	Domain string `xml:"domain,attr,omitempty"`
	// Text The value of the element is a forward-slash-separated string that identifies a hierarchic location in the indicated taxonomy. Processors may establish conventions for the interpretation of categories.
	Text string `xml:",chardata"`
}

// RssImage is an optional sub-element of <channel>, which contains three required and three optional sub-elements.
type RssImage struct {
	// Url is the URL of a GIF, JPEG or PNG image that represents the channel.
	Url string `xml:"url,omitempty"`
	// Title describes the image, it's used in the ALT attribute of the HTML <img> tag when the channel is rendered in HTML.
	Title string `xml:"title,omitempty"`
	// Link is the URL of the site, when the channel is rendered, the image is a link to the site. (Note, in practice the image <title> and <link> should have the same value as the channel's <title> and <link>.
	Link string `xml:"link,omitempty"`
	// Width indicating the width of the image in pixels. Maximum value for width is 144, default value is 88.
	Width string `xml:"width,omitempty"`
	// Height indicating the height of the image in pixels. Maximum value for height is 400, default value is 31.
	Height string `xml:"height,omitempty"`
	// Description contains text that is included in the TITLE attribute of the link formed around the image in the HTML
	Description string `xml:"description,omitempty"`
}

// RssTextInput A channel may optionally contain a <textInput> sub-element, which contains four required sub-elements. The purpose of the <textInput> element is something of a mystery. You can use it to specify a search engine box. Or to allow a reader to provide feedback. Most aggregators ignore it.
type RssTextInput struct {
	// Title The label of the Submit button in the text input area.
	Title string `xml:"title,omitempty"`
	// Description Explains the text input area.
	Description string `xml:"description,omitempty"`
	// Name The name of the text object in the text input area.
	Name string `xml:"name,omitempty"`
	// Link The URL of the CGI script that processes text input requests.
	Link string `xml:"link,omitempty"`
}

// RssSkipHours An XML element that contains up to 24 <hour> sub-elements whose value is a number between 0 and 23, representing a time in GMT, when aggregators, if they support the feature, may not read the channel on hours listed in the skipHours element. The hour beginning at midnight is hour zero.
type RssSkipHours struct {
	Hours []string `xml:"hour,omitempty"`
}

// RssSkipDays An XML element that contains up to seven <day> sub-elements whose value is Monday, Tuesday, Wednesday, Thursday, Friday, Saturday or Sunday. Aggregators may not read the channel during days listed in the skipDays element.
type RssSkipDays struct {
	Days []string `xml:"day,omitempty"`
}
