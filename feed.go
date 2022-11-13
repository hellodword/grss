package grss

import (
	"encoding/json"
	"github.com/nbio/xml"
	"html"
	"io"
	"sort"
	"strconv"
	"time"
)

// https://validator.w3.org/feed/docs/howto/declare_namespaces.html
// https://validator.w3.org/feed/#validate_by_input
// https://github.com/w3c/feedvalidator/tree/main/src
// https://rss.com/blog/rss-feed-validators/
// https://podba.se/validate/
// https://www.rssboard.org/rss-validator/
// https://github.com/andre487/feed-validator
// https://www.castfeedvalidator.com/

type Feed interface {
	Uniform()
	Mime(fallback bool) string
	ToJSON() *JSONFeed
	ToRss() *RssFeed
	ToAtom() *AtomFeed
	WriteOut(w io.Writer) error
}

func (f *JSONFeed) Uniform() {
	f.Version = "https://jsonfeed.org/version/1.1"
}

func (f *JSONFeed) Mime(fallback bool) string {
	if fallback {
		return JSONMimeFallback
	} else {
		return JSONMime
	}
}

func (f *JSONFeed) ToJSON() *JSONFeed {
	ff := &JSONFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       f.Title,
		HomePageURL: f.HomePageURL,
		FeedURL:     f.FeedURL,
		Description: f.Description,
		UserComment: f.UserComment,
		NextURL:     f.NextURL,
		Icon:        f.Icon,
		Favicon:     f.Favicon,
		Expired:     f.Expired,
		Items:       f.Items,
		Hub:         f.Hub,
		Extensions:  f.Extensions,
		Language:    f.Language,
	}

	if f.Author != nil {
		ff.Authors = append(ff.Authors, f.Author)
	}

	ff.Authors = append(ff.Authors, f.Authors...)

	ff.Uniform()
	return ff
}

func (f *JSONFeed) ToRss() *RssFeed {
	// https://www.jsonfeed.org/mappingrssandatom/#rss
	ff := &RssFeed{}

	ff.Channel = &RssChannel{}

	// RSS title and description map directly to JSON.
	ff.Channel.Title = XmlText{
		Text: f.Title,
	}

	ff.Channel.Description = XmlText{
		Text: f.Description,
	}

	// An RSS link maps to home_page_url.
	ff.Channel.Link = f.HomePageURL

	// RSS has webmaster and managingEditor items, while JSON Feed has an author item.
	authors := f.Authors
	if f.Author != nil {
		authors = append([]*JSONAuthor{f.Author}, f.Authors...)
	}
	if len(authors) > 0 {
		ff.Channel.WebMaster = authors[0].Name
		ff.Channel.ManagingEditor = authors[0].Name
	}

	if len(authors) > 1 {
		ff.Channel.ManagingEditor = authors[1].Name
	}

	// image is superficially like a JSON Feed icon — except that the JSON Feed version should be square, and the RSS image should be wider than it is tall. These map only in the case that your RSS image is already square.
	if f.Icon != "" {
		ff.Channel.Image = &RssImage{
			Url:         f.Icon,
			Description: "icon",
		}
	}

	for _, jitem := range f.Items {
		// https://www.jsonfeed.org/mappingrssandatom/#item
		item := &RssItem{}
		ff.Channel.Items = append(ff.Channel.Items, item)

		// title maps directly, with the caveat that it must be plain text in JSON Feed. (Note: title is optional in both RSS and JSON Feed.)
		item.Title = jitem.Title

		// link maps to url and to external_url. The url must be a permalink, a link to the content of the item. When an item links to another page — as in a linkblog — then external_url must be the URL of that other page. (See the Daring Fireball JSON feed for an example.)
		if jitem.URL != "" {
			item.Link = jitem.URL
		} else {
			item.Link = jitem.ExternalURL
		}

		// description maps to content_html and to content_text. Choose the one that fits. A Twitter-like service might use context_text, while a blog might use content_html.
		// TODO really?
		if jitem.ContentHTML != "" {
			item.ContentEncoded = &RssContent{
				XMLName: xml.Name{
					Space: "http://purl.org/rss/1.0/modules/content/",
					Local: "encoded",
				},
				XmlText: XmlText{
					Cdata: html.EscapeString(html.UnescapeString(jitem.ContentHTML)),
				},
			}
		}

		if jitem.ContentText != "" {
			item.ContentEncoded = &RssContent{
				XMLName: xml.Name{
					Space: "http://purl.org/rss/1.0/modules/content/",
					Local: "encoded",
				},
				XmlText: XmlText{
					Cdata: html.EscapeString(html.UnescapeString(jitem.ContentText)),
				},
			}
		}

		// guid maps to id. In RSS, guid can have an isPermaLink attribute; in JSON Feed the url must be the permalink, and id may be the same as url, though it doesn’t have to be.
		if jitem.ID != "" {
			item.Guid = &RssGuid{
				Guid: jitem.ID,
			}
		}

		// The author element is a single value, while in JSON Feed it’s an object with name, url, and avatar values.
		jauthors := jitem.Authors
		if jitem.Author != nil {
			jauthors = append([]*JSONAuthor{jitem.Author}, jitem.Authors...)
		}
		if len(jauthors) > 0 {
			item.Author = &RssAuthor{Email: jauthors[0].Name}
		}

		// pubDate maps to date_published, but this date and others in JSON Feed use a different format, the RFC 3339 format. (Example: 2010-02-07T14:04:00-05:00.)
		if jitem.DatePublished != "" {
			item.PubDate = FormatDate(jitem.DatePublished, time.RFC1123Z)
		}

		// enclosure maps to attachments — but JSON Feed allows for multiple attachments. An RSS enclosure has attributes url, length, and type, and the JSON Feed attachment object has corresponding elements url, size_in_bytes, and mime_type. JSON Feed adds title and duration_in_seconds.
		if len(jitem.Attachments) > 0 {
			item.Enclosure = &RssEnclosure{
				Url:    jitem.Attachments[0].URL,
				Length: strconv.FormatUint(jitem.Attachments[0].SizeInBytes, 10),
				Type:   jitem.Attachments[0].MimeType,
			}
		}

		// category maps to tags. In RSS a category can have a domain attribute, and there’s no equivalent in JSON Feed.
		for i := range jitem.Tags {
			item.Categories = append(item.Categories, &RssCategory{
				Text: jitem.Tags[i],
			})
		}
	}

	ff.Uniform()
	return ff
}

func (f *JSONFeed) ToAtom() *AtomFeed {
	// https://www.jsonfeed.org/mappingrssandatom/#atom
	ff := &AtomFeed{}

	// Atom’s title maps directly to JSON Feed.
	if f.Title != "" {
		ff.Title = &AtomTextConstruct{
			XmlText: XmlText{
				Text: f.Title,
			},
		}
	}

	// Atom uses link elements to show the relationship between the feed and other URLs. Without a rel attribute, or with rel="alternate", Atom’s link is equivalent to JSON Feed’s home_page_url.
	if f.HomePageURL != "" {
		ff.Links = append(ff.Links, &AtomLink{
			Href: AtomUri(f.HomePageURL),
			Rel:  "alternate",
		})
	}

	// With rel="self", it maps to JSON Feed’s feed_url.
	if f.FeedURL != "" {
		ff.Links = append(ff.Links, &AtomLink{
			Href: AtomUri(f.FeedURL),
			Rel:  "self",
		})
	}

	// Atom has subtitle and id elements. There is no mapping for these in JSON Feed, although description in JSON could be used instead of subtitle.
	if f.Description != "" {
		ff.Subtitle = &AtomTextConstruct{
			XmlText: XmlText{
				Text: f.Description,
			},
		}
	}

	if f.Icon != "" {
		ff.Icon = &AtomIcon{
			AtomUri: AtomUri(f.Icon),
		}
	}

	authors := f.Authors
	if f.Author != nil {
		authors = append([]*JSONAuthor{f.Author}, f.Authors...)
	}
	for _, author := range authors {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			Name: author.Name,
			Uri:  AtomUri(author.URL),
		})
	}

	ff.Language = AtomLanguageTag(f.Language)

	// https://www.jsonfeed.org/mappingrssandatom/#entry
	for _, item := range f.Items {
		entry := &AtomEntry{}
		ff.Entries = append(ff.Entries, entry)

		// Atom’s title, id, and summary map directly to JSON. In JSON, however, title and summary are always plain text.
		if item.ID != "" {
			entry.ID = &AtomId{
				AtomUri: AtomUri(item.ID),
			}
		}

		if item.Title != "" {
			entry.Title = &AtomTextConstruct{
				XmlText: XmlText{
					Text: item.Title,
				},
			}
		}

		if item.Summary != "" {
			entry.Summary = &AtomTextConstruct{
				XmlText: XmlText{
					Text: item.Summary,
				},
			}
		}

		// Atom uses a content’s type attribute to declare whether the text is plain text or HTML. type="html" or type="xhtml" in an Atom feed maps to content_html in JSON. type="text" maps to content_text in JSON.
		if item.ContentText != "" {
			entry.Content = &AtomContent{
				XmlText: XmlText{
					Cdata: item.ContentText,
				},
			}
		}

		if item.ContentHTML != "" {
			entry.Content = &AtomContent{
				Type: "html",
				XmlText: XmlText{
					Cdata: html.EscapeString(html.UnescapeString(item.ContentHTML)),
				},
			}
		}

		// link with rel="alternate" maps to url in JSON.
		if item.URL != "" {
			entry.Links = append(entry.Links, &AtomLink{
				Href: AtomUri(item.URL),
				Rel:  "alternate",
			})
		}

		// If rel="related" is used for links to an external site, in JSON Feed those map to external_url.
		if item.ExternalURL != "" {
			entry.Links = append(entry.Links, &AtomLink{
				Href: AtomUri(item.ExternalURL),
				Rel:  "related",
			})
		}

		// The author element contains name, uri, and email, while in JSON Feed it’s an object with name, url, and avatar values.
		for i := range item.Authors {
			entry.Authors = append(entry.Authors, &AtomPersonConstruct{
				Name: item.Authors[i].Name,
				Uri:  AtomUri(item.Authors[i].URL),
			})
		}

		// Atom’s published and updated dates map to date_published and date_modified in JSON. Both Atom and JSON Feed use the same date format.
		if item.DatePublished != "" {
			entry.Published = &AtomDateConstruct{
				DateTime: FormatDate(item.DatePublished, time.RFC3339),
			}
		}

		if item.DateModified != "" {
			entry.Updated = &AtomDateConstruct{
				DateTime: FormatDate(item.DateModified, time.RFC3339),
			}
		}

		// Atom’s link with rel="enclosure" maps to attachments in JSON Feed. An Atom enclosure has attributes href, length, and type, and the JSON Feed attachment object has corresponding elements url, size_in_bytes, and mime_type. JSON Feed adds title and duration_in_seconds.
		for i := range item.Attachments {
			entry.Links = append(entry.Links, &AtomLink{
				Rel:    "enclosure",
				Href:   AtomUri(item.Attachments[i].URL),
				Length: strconv.FormatUint(item.Attachments[i].SizeInBytes, 10),
				Type:   AtomMediaType(item.Attachments[i].MimeType),
				Title:  item.Attachments[i].Title,
			})
		}

		entry.Language = AtomLanguageTag(item.Language)
	}

	ff.Uniform()
	return ff
}

func (f *JSONFeed) WriteOut(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetEscapeHTML(true)
	e.SetIndent("", "    ")
	return e.Encode(f)
}

func (f *RssFeed) Uniform() {
	f.Version = "2.0"
	f.XMLName = xml.Name{
		Space: "",
		Local: "rss",
	}

	pre := [][3]string{
		{"http://www.w3.org/2000/xmlns/", "media", "http://search.yahoo.com/mrss/"},
		{"http://www.w3.org/2000/xmlns/", "dc", "http://purl.org/dc/elements/1.1/"},
		{"http://www.w3.org/2000/xmlns/", "content", "http://purl.org/rss/1.0/modules/content/"},
		{"http://www.w3.org/2000/xmlns/", "atom", "http://www.w3.org/2005/Atom"},
	}

	f.Attributes = append(f.Attributes, addAttrs(pre, f.Attributes)...)

	for _, item := range f.Channel.Items {
		if item.ContentEncoded != nil && item.Description == "" {
			if item.Title != "" {
				item.Description = item.Title
			} else if item.Link != "" {
				item.Description = item.Link
			} else {
				item.Description = "Unknown"
			}
		}
	}
}

func (f *RssFeed) Mime(fallback bool) string {
	if fallback {
		return RssMimeFallback
	} else {
		return RssMime
	}
}

func (f *RssFeed) ToJSON() *JSONFeed {
	// https://www.jsonfeed.org/mappingrssandatom/#rss
	ff := &JSONFeed{}

	if f.Channel == nil {
		return ff
	}

	// RSS title and description map directly to JSON.
	ff.Title = f.Channel.Title.String()
	ff.Description = f.Channel.Description.String()

	// An RSS link maps to home_page_url.
	ff.HomePageURL = f.Channel.Link

	// RSS has webmaster and managingEditor items, while JSON Feed has an author item.
	if f.Channel.WebMaster != "" {
		ff.Authors = append(ff.Authors, &JSONAuthor{
			Name: f.Channel.WebMaster,
		})
	}

	if f.Channel.ManagingEditor != "" {
		ff.Authors = append(ff.Authors, &JSONAuthor{
			Name: f.Channel.ManagingEditor,
		})
	}

	// image is superficially like a JSON Feed icon — except that the JSON Feed version should be square, and the RSS image should be wider than it is tall. These map only in the case that your RSS image is already square.
	if f.Channel.Image != nil {
		ff.Icon = f.Channel.Image.Url
	}

	for _, item := range append(f.Items, f.Channel.Items...) {
		// https://www.jsonfeed.org/mappingrssandatom/#item
		jitem := &JSONItem{}
		ff.Items = append(ff.Items, jitem)

		// title maps directly, with the caveat that it must be plain text in JSON Feed. (Note: title is optional in both RSS and JSON Feed.)
		jitem.Title = item.Title

		// link maps to url and to external_url. The url must be a permalink, a link to the content of the item. When an item links to another page — as in a linkblog — then external_url must be the URL of that other page. (See the Daring Fireball JSON feed for an example.)
		jitem.URL = item.Link

		// description maps to content_html and to content_text. Choose the one that fits. A Twitter-like service might use context_text, while a blog might use content_html.
		// TODO really?
		if item.ContentEncoded != nil {
			jitem.ContentHTML = item.ContentEncoded.XmlText.String()
		} else if item.Description != "" {
			jitem.ContentText = item.Description
		}

		// guid maps to id. In RSS, guid can have an isPermaLink attribute; in JSON Feed the url must be the permalink, and id may be the same as url, though it doesn’t have to be.
		if item.Guid != nil {
			jitem.ID = item.Guid.Guid
		}

		// The author element is a single value, while in JSON Feed it’s an object with name, url, and avatar values.
		if item.Author != nil {
			jitem.Authors = []*JSONAuthor{
				{
					Name: item.Author.Email,
				},
			}
		}

		// pubDate maps to date_published, but this date and others in JSON Feed use a different format, the RFC 3339 format. (Example: 2010-02-07T14:04:00-05:00.)
		jitem.DatePublished = item.PubDate
		//jitem.DateModified = item.PubDate

		// enclosure maps to attachments — but JSON Feed allows for multiple attachments. An RSS enclosure has attributes url, length, and type, and the JSON Feed attachment object has corresponding elements url, size_in_bytes, and mime_type. JSON Feed adds title and duration_in_seconds.
		if item.Enclosure != nil {
			attachment := &JSONAttachments{
				URL:      item.Enclosure.Url,
				MimeType: item.Enclosure.Type,
			}
			jitem.Attachments = append(jitem.Attachments, attachment)
			if length, err := strconv.ParseUint(item.Enclosure.Length, 10, 0); err == nil {
				attachment.SizeInBytes = length
			}
		}

		//if item.Content != nil {
		//	jitem.ContentText = item.Content.XmlText.String()
		//}

		// category maps to tags. In RSS a category can have a domain attribute, and there’s no equivalent in JSON Feed.
		for i := range item.Categories {
			jitem.Tags = append(jitem.Tags, item.Categories[i].Text)
		}
	}

	ff.Language = f.Channel.Language

	ff.Uniform()
	return ff
}

func (f *RssFeed) ToRss() *RssFeed {
	ff := &RssFeed{
		XMLName: xml.Name{
			Space: "",
			Local: "rss",
		},
		Attributes: f.Attributes,
		Version:    "2.0",
		//XmlnsContent: f.XmlnsContent,
		Channel:   nil,
		Image:     nil,
		Items:     nil,
		TextInput: nil,
	}

	if f.Channel == nil {
		return ff
	}

	ff.Channel = &RssChannel{
		Attributes:       f.Channel.Attributes,
		Title:            f.Channel.Title,
		Link:             f.Channel.Link,
		Description:      f.Channel.Description,
		Language:         f.Channel.Language,
		Copyright:        f.Channel.Copyright,
		ManagingEditor:   f.Channel.ManagingEditor,
		WebMaster:        f.Channel.WebMaster,
		PubDate:          f.Channel.PubDate,
		LastBuildDate:    f.Channel.LastBuildDate,
		Categories:       f.Channel.Categories,
		Generator:        f.Channel.Generator,
		Docs:             f.Channel.Docs,
		Cloud:            f.Channel.Cloud,
		Ttl:              f.Channel.Ttl,
		Image:            f.Channel.Image,
		Rating:           f.Channel.Rating,
		TextInput:        f.Channel.TextInput,
		SkipHours:        f.Channel.SkipHours,
		SkipDays:         f.Channel.SkipDays,
		Items:            nil,
		ExtensionElement: f.Channel.ExtensionElement,
	}

	if ff.Channel.Image == nil {
		ff.Channel.Image = f.Image
	}

	if ff.Channel.TextInput == nil {
		ff.Channel.TextInput = f.TextInput
	}

	ff.Channel.Items = append(f.Channel.Items, f.Items...)

	ff.Uniform()
	return ff
}

func (f *RssFeed) ToAtom() *AtomFeed {
	ff := &AtomFeed{
		XMLName: xml.Name{
			Space: "",
			Local: "feed",
		},
		AtomCommonAttributes: AtomCommonAttributes{},
		//Authors:              nil,
		//Categories:       nil,
		Contributors: nil,
		Generator:    nil,
		Icon:         nil,
		ID: AtomId{
			AtomCommonAttributes: AtomCommonAttributes{},
			AtomUri:              "",
		},
		//Links:            nil,
		Logo:     nil,
		Rights:   nil,
		Subtitle: nil,
		//Title:            nil,
		//Updated:          nil,
		ExtensionElement: nil,
		//Entries:          nil,
	}

	if f.Channel == nil {
		return ff
	}

	if f.Channel.WebMaster != "" {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Name:                 f.Channel.WebMaster,
			Uri:                  "",
			Email:                AtomEmailAddress(f.Channel.WebMaster),
		})
	}

	if f.Channel.ManagingEditor != "" {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Name:                 f.Channel.ManagingEditor,
			Uri:                  "",
			Email:                AtomEmailAddress(f.Channel.ManagingEditor),
		})
	}

	for i := range f.Channel.Categories {
		ff.Categories = append(ff.Categories, &AtomCategory{
			AtomCommonAttributes: AtomCommonAttributes{},
			Term:                 f.Channel.Categories[i].Text,
			Scheme:               AtomUri(f.Channel.Categories[i].Domain),
			Label:                f.Channel.Categories[i].Text,
			UndefinedContent:     nil,
		})
	}

	if f.Channel.Link != "" {
		ff.ID.AtomUri = AtomUri(f.Channel.Link)

		ff.Links = []*AtomLink{
			{
				AtomCommonAttributes: AtomCommonAttributes{},
				Href:                 AtomUri(f.Channel.Link),
				Rel:                  "",
				Type:                 "",
				Hreflang:             "",
				Title:                "",
				Length:               "",
				UndefinedContent:     nil,
			},
		}
	}

	if f.Channel.Title.String() != "" {
		ff.Title = &AtomTextConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			XmlText:              f.Channel.Title,
		}
	}

	if f.Channel.PubDate != "" {
		ff.Updated = &AtomDateConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			DateTime:             FormatDate(f.Channel.PubDate, time.RFC3339),
		}
	} else if f.Channel.LastBuildDate != "" {
		ff.Updated = &AtomDateConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			DateTime:             FormatDate(f.Channel.LastBuildDate, time.RFC3339),
		}
	}

	for _, item := range append(f.Items, f.Channel.Items...) {
		entry := &AtomEntry{
			AtomCommonAttributes: AtomCommonAttributes{},
			//Authors:              nil,
			//Categories:       nil,
			//Content:          nil,
			Contributors: nil,
			//ID:               nil,
			//Links:            nil,
			//Published:        nil,
			Rights: nil,
			Source: nil,
			//Summary: nil,
			//Title: nil,
			//Updated:          nil,
			ExtensionElement: nil,
		}
		ff.Entries = append(ff.Entries, entry)

		if item.Author != nil {
			entry.Authors = []*AtomPersonConstruct{
				{
					AtomCommonAttributes: AtomCommonAttributes{},
					Name:                 item.Author.Email,
				},
			}
		}

		for i := range item.Categories {
			ff.Categories = append(ff.Categories, &AtomCategory{
				AtomCommonAttributes: AtomCommonAttributes{},
				Term:                 item.Categories[i].Text,
				Scheme:               AtomUri(item.Categories[i].Domain),
				Label:                item.Categories[i].Text,
				UndefinedContent:     nil,
			})
		}

		if item.ContentEncoded != nil {
			entry.Content = &AtomContent{
				AtomCommonAttributes: AtomCommonAttributes{},
				XmlText:              item.ContentEncoded.XmlText,
			}
		}
		//else if item.Content != nil {
		//	entry.Content = &AtomContent{
		//		AtomCommonAttributes: AtomCommonAttributes{},
		//		XmlText:              item.Content.XmlText,
		//	}
		//}

		if item.Guid != nil {
			entry.ID = &AtomId{
				AtomCommonAttributes: AtomCommonAttributes{},
				AtomUri:              AtomUri(item.Guid.Guid),
			}
		}

		if item.Link != "" {
			entry.Links = []*AtomLink{
				{
					AtomCommonAttributes: AtomCommonAttributes{},
					Href:                 AtomUri(item.Link),
					Rel:                  "",
					Type:                 "",
					Hreflang:             "",
					Title:                "",
					Length:               "",
					UndefinedContent:     nil,
				},
			}
		}

		if item.PubDate != "" {
			entry.Published = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             FormatDate(item.PubDate, time.RFC3339),
			}

			entry.Updated = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             FormatDate(item.PubDate, time.RFC3339),
			}
		}

		if entry.Content == nil && item.Description != "" {
			entry.Summary = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				XmlText: XmlText{
					Cdata: item.Description,
				},
			}
		}

		if item.Title != "" {
			entry.Title = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				XmlText: XmlText{
					Text: item.Title,
				},
			}
		}

	}

	ff.Uniform()
	return ff
}

func (f *RssFeed) WriteOut(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "    ")
	return e.Encode(f)
}

func (f *AtomFeed) Uniform() {
	now := time.Now()

	f.XMLName = xml.Name{
		Space: "http://www.w3.org/2005/Atom",
		Local: "feed",
	}

	pre := [][3]string{
		{"", "xmlns", "http://www.w3.org/2005/Atom"},
		{"http://www.w3.org/2000/xmlns/", "media", "http://search.yahoo.com/mrss/"},
	}

	f.UndefinedAttribute = append(f.UndefinedAttribute, addAttrs(pre, f.UndefinedAttribute)...)

	for _, entry := range f.Entries {
		if entry.Updated != nil {
			entry.Updated.DateTime = FormatDate(entry.Updated.DateTime, time.RFC3339)
		}
		if entry.Published != nil {
			entry.Published.DateTime = FormatDate(entry.Published.DateTime, time.RFC3339)
		}
	}

	if f.Updated == nil {
		var ts sort.StringSlice
		for _, entry := range f.Entries {
			if entry.Updated != nil {
				ts = append(ts, FormatDate(entry.Updated.DateTime, time.RFC3339))
			} else if entry.Published != nil {
				ts = append(ts, FormatDate(entry.Published.DateTime, time.RFC3339))
			}
		}

		sort.Sort(sort.Reverse(ts))

		if len(ts) > 0 {
			f.Updated = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             ts[0],
			}
		} else {
			f.Updated = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             now.Format(time.RFC3339),
			}
		}

	}

	for _, entry := range f.Entries {
		if len(entry.Authors) > 0 {
			continue
		}

		if len(f.Authors) > 0 {
			entry.Authors = []*AtomPersonConstruct{f.Authors[0]}
			continue
		}

		entry.Authors = []*AtomPersonConstruct{
			{
				Name: "Unknown (grss)",
			},
		}
	}

	for _, entry := range f.Entries {
		if entry.Updated != nil {
			continue
		}

		if entry.Published != nil {
			entry.Updated = entry.Published
			continue
		}

		entry.Updated = &AtomDateConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			DateTime:             now.Format(time.RFC3339),
		}
	}

	for _, author := range f.Authors {
		if author.Name == "" {
			if author.Email != "" {
				author.Name = string(author.Email)
			} else {
				author.Name = "Unknown (grss)"
			}
		}
	}

	for _, entry := range f.Entries {
		for _, author := range entry.Authors {
			if author.Name == "" {
				if author.Email != "" {
					author.Name = string(author.Email)
				} else {
					author.Name = "Unknown (grss)"
				}
			}
		}
	}

	for _, entry := range f.Entries {
		if entry.ID == nil {
			if len(entry.Links) > 0 && entry.Links[0].Href != "" {
				entry.ID = &AtomId{
					AtomCommonAttributes: AtomCommonAttributes{},
					AtomUri:              entry.Links[0].Href,
				}
			}
		}
	}

}

func (f *AtomFeed) Mime(fallback bool) string {
	if fallback {
		return AtomMimeFallback
	} else {
		return AtomMime
	}
}

func (f *AtomFeed) ToJSON() *JSONFeed {
	// https://www.jsonfeed.org/mappingrssandatom/#atom
	ff := &JSONFeed{}

	// Atom’s title maps directly to JSON Feed.
	if f.Title != nil {
		ff.Title = f.Title.String()
	}

	// Atom uses link elements to show the relationship between the feed and other URLs. Without a rel attribute, or with rel="alternate", Atom’s link is equivalent to JSON Feed’s home_page_url.
	// With rel="self", it maps to JSON Feed’s feed_url.
	for i := range f.Links {
		switch f.Links[i].Rel {
		case "self":
			ff.FeedURL = string(f.Links[i].Href)
		case "", "alternate":
			ff.HomePageURL = string(f.Links[i].Href)
		}
	}

	// Atom has subtitle and id elements. There is no mapping for these in JSON Feed, although description in JSON could be used instead of subtitle.
	if f.Subtitle != nil {
		ff.Description = f.Subtitle.XmlText.String()
	}

	if f.Icon != nil {
		ff.Icon = string(f.Icon.AtomUri)
	}

	for i := range f.Authors {
		author := &JSONAuthor{
			Name: f.Authors[i].Name,
			URL:  string(f.Authors[i].Uri),
		}
		if author.Name == "" {
			author.Name = string(f.Authors[i].Email)
		}
		ff.Authors = append(ff.Authors, author)
	}

	ff.Language = string(f.Language)

	// https://www.jsonfeed.org/mappingrssandatom/#entry
	for _, entry := range f.Entries {
		item := &JSONItem{}
		ff.Items = append(ff.Items, item)

		// Atom’s title, id, and summary map directly to JSON. In JSON, however, title and summary are always plain text.
		if entry.ID != nil {
			item.ID = string(entry.ID.AtomUri)
		}

		if entry.Title != nil {
			item.Title = entry.Title.String()
		}

		if entry.Summary != nil {
			item.Summary = entry.Summary.String()
		}

		// Atom uses a content’s type attribute to declare whether the text is plain text or HTML. type="html" or type="xhtml" in an Atom feed maps to content_html in JSON. type="text" maps to content_text in JSON.
		if entry.Content != nil {
			switch entry.Content.Type {
			case "html", "xhtml":
				item.ContentHTML = entry.Content.String()
			default:
				item.ContentText = entry.Content.String()
			}
		}

		// link with rel="alternate" maps to url in JSON.
		// If rel="related" is used for links to an external site, in JSON Feed those map to external_url.
		// Atom’s link with rel="enclosure" maps to attachments in JSON Feed. An Atom enclosure has attributes href, length, and type, and the JSON Feed attachment object has corresponding elements url, size_in_bytes, and mime_type. JSON Feed adds title and duration_in_seconds.
		for i := range entry.Links {
			switch entry.Links[i].Rel {
			case "alternate":
				item.URL = string(entry.Links[i].Href)
			case "related":
				item.ExternalURL = string(entry.Links[i].Href)
			case "enclosure":
				attachment := &JSONAttachments{
					URL:      string(entry.Links[i].Href),
					MimeType: string(entry.Links[i].Type),
					Title:    entry.Links[i].Title,
				}
				item.Attachments = append(item.Attachments, attachment)
				if entry.Links[i].Length != "" {
					if length, err := strconv.ParseUint(entry.Links[i].Length, 10, 0); err == nil {
						attachment.SizeInBytes = length
					}
				}
			}
		}

		// The author element contains name, uri, and email, while in JSON Feed it’s an object with name, url, and avatar values.
		for i := range entry.Authors {
			item.Authors = append(item.Authors, &JSONAuthor{
				Name: entry.Authors[i].Name,
				URL:  string(entry.Authors[i].Uri),
			})
		}

		// Atom’s published and updated dates map to date_published and date_modified in JSON. Both Atom and JSON Feed use the same date format.
		if entry.Published != nil {
			item.DatePublished = entry.Published.DateTime
		}

		if entry.Updated != nil {
			item.DateModified = entry.Updated.DateTime
		}

		item.Language = string(entry.Language)
	}

	ff.Uniform()
	return ff
}

func (f *AtomFeed) ToRss() *RssFeed {
	ff := &RssFeed{
		XMLName: xml.Name{
			Space: "",
			Local: "rss",
		},
		Attributes: nil,
		Version:    "2.0",
		//XmlnsContent: "",
		Channel:   nil,
		Image:     nil,
		Items:     nil,
		TextInput: nil,
	}

	ff.Channel = &RssChannel{
		Attributes: nil,
		//Title:       "",
		//Link:        "",
		//Description: "",
		Language:  "",
		Copyright: "",
		//ManagingEditor:   "",
		//WebMaster:        "",
		//PubDate:       "",
		//LastBuildDate: "",
		//Categories:       nil,
		Generator: "",
		Docs:      "",
		Cloud:     nil,
		Ttl:       "",
		Image:     nil,
		Rating:    nil,
		TextInput: nil,
		SkipHours: nil,
		SkipDays:  nil,
		//Items:            nil,
		ExtensionElement: nil,
	}

	if len(f.Authors) > 0 {
		ff.Channel.WebMaster = string(f.Authors[0].Email)
	}

	if len(f.Authors) > 1 {
		ff.Channel.ManagingEditor = string(f.Authors[0].Email)
	}

	for i := range f.Categories {
		ff.Channel.Categories = append(ff.Channel.Categories, &RssCategory{
			Domain: string(f.Categories[i].Scheme),
			Text:   f.Categories[i].Label,
		})
	}

	if f.Updated != nil {
		ff.Channel.PubDate = FormatDate(f.Updated.DateTime, time.RFC1123Z)
		ff.Channel.LastBuildDate = FormatDate(f.Updated.DateTime, time.RFC1123Z)
	}

	for _, entry := range f.Entries {
		item := &RssItem{
			//Title: "",
			//Link:        "",
			Description: "",
			//Author:         nil,
			//Categories:     nil,
			Comments:  "",
			Enclosure: nil,
			//Guid:      nil,
			//PubDate: "",
			Source: nil,
			//Content:        nil,
			ContentEncoded: nil,
		}
		ff.Channel.Items = append(ff.Channel.Items, item)

		if len(entry.Authors) > 0 {
			item.Author = &RssAuthor{Email: string(entry.Authors[0].Email)}
		}

		for i := range entry.Categories {
			item.Categories = append(item.Categories, &RssCategory{
				Domain: string(entry.Categories[i].Scheme),
				Text:   entry.Categories[i].Label,
			})
		}

		if entry.Content != nil {
			item.ContentEncoded = &RssContent{
				XMLName: xml.Name{
					Space: "http://purl.org/rss/1.0/modules/content/",
					Local: "encoded",
				},
				XmlText: entry.Content.XmlText,
			}
		}

		if entry.ID != nil {
			item.Guid = &RssGuid{
				IsPermaLink: "",
				Guid:        string(entry.ID.AtomUri),
			}
		}

		if len(entry.Links) > 0 {
			item.Link = string(entry.Links[0].Href)
		}

		if entry.Published != nil {
			item.PubDate = FormatDate(entry.Published.DateTime, time.RFC1123Z)
		} else if entry.Updated != nil {
			item.PubDate = FormatDate(entry.Updated.DateTime, time.RFC1123Z)
		}

		if entry.Summary != nil {
			item.Description = entry.Summary.String()
		}

		if entry.Title != nil {
			item.Title = entry.Title.String()
		}
	}

	if f.Title != nil {
		ff.Channel.Title = f.Title.XmlText
	}

	if len(f.Links) > 0 {
		ff.Channel.Link = string(f.Links[0].Href)
	}

	ff.Uniform()
	return ff
}

func (f *AtomFeed) ToAtom() *AtomFeed {
	ff := &AtomFeed{
		XMLName: xml.Name{
			Space: "",
			Local: "feed",
		},
		AtomCommonAttributes: f.AtomCommonAttributes,
		Authors:              f.Authors,
		Categories:           f.Categories,
		Contributors:         f.Contributors,
		Generator:            f.Generator,
		Icon:                 f.Icon,
		ID:                   f.ID,
		Links:                f.Links,
		Logo:                 f.Logo,
		Rights:               f.Rights,
		Subtitle:             f.Subtitle,
		Title:                f.Title,
		Updated:              f.Updated,
		ExtensionElement:     f.ExtensionElement,
		Entries:              f.Entries,
	}

	ff.Uniform()
	return ff
}

func (f *AtomFeed) WriteOut(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "    ")
	return e.Encode(f)
}
