package grss

import (
	"encoding/json"
	"github.com/nbio/xml"
	"html"
	"io"
	"sort"
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
	defer ff.Uniform()

	if f.Author != nil {
		ff.Authors = append(ff.Authors, f.Author)
	}

	ff.Authors = append(ff.Authors, f.Authors...)

	return ff
}

func (f *JSONFeed) ToRss() *RssFeed {
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
		Attributes:       nil,
		Title:            "",
		Link:             "",
		Description:      "",
		Language:         "",
		Copyright:        "",
		ManagingEditor:   "",
		WebMaster:        "",
		PubDate:          "",
		LastBuildDate:    "",
		Categories:       nil,
		Generator:        "",
		Docs:             "",
		Cloud:            nil,
		Ttl:              "",
		Image:            nil,
		Rating:           nil,
		TextInput:        nil,
		SkipHours:        nil,
		SkipDays:         nil,
		Items:            nil,
		ExtensionElement: nil,
	}
	defer ff.Uniform()

	ff.Channel.Title = f.Title
	ff.Channel.Link = f.HomePageURL
	ff.Channel.Description = f.Description

	authors := f.Authors
	if f.Author != nil {
		authors = append([]*JSONAuthor{f.Author}, f.Authors...)
	}
	if len(authors) > 0 {
		ff.Channel.WebMaster = authors[0].Name
	}

	if len(authors) > 1 {
		ff.Channel.ManagingEditor = authors[1].Name
	}

	for _, jitem := range f.Items {
		item := &RssItem{
			Title:          "",
			Link:           "",
			Description:    "",
			Author:         nil,
			Categories:     nil,
			Comments:       "",
			Enclosure:      nil,
			Guid:           nil,
			PubDate:        "",
			Source:         nil,
			Content:        nil,
			ContentEncoded: nil,
		}
		ff.Channel.Items = append(ff.Channel.Items, item)

		if jitem.ID != "" {
			item.Guid = &RssGuid{
				IsPermaLink: "",
				Guid:        jitem.ID,
			}
		}

		item.Link = jitem.URL

		if jitem.ExternalURL != "" {
			item.Source = &RssSource{
				Url:  jitem.ExternalURL,
				Text: "",
			}
		}

		item.Title = jitem.Title

		if jitem.Summary != "" {
			item.Description = jitem.Summary
		}

		if jitem.ContentHTML != "" {
			item.ContentEncoded = &RssContent{
				XMLName: xml.Name{
					Space: "http://purl.org/rss/1.0/modules/content/",
					Local: "encoded",
				},
				Content: []byte(html.EscapeString(html.UnescapeString(jitem.ContentHTML))),
			}
		}

		if jitem.ContentText != "" {
			item.ContentEncoded = &RssContent{
				XMLName: xml.Name{
					Space: "http://purl.org/rss/1.0/modules/content/",
					Local: "encoded",
				},
				Content: []byte(html.EscapeString(html.UnescapeString(jitem.ContentText))),
			}
		}

		if jitem.DatePublished == "" {
			item.PubDate = FormatDate(jitem.DateModified, time.RFC1123Z)
		} else {
			item.PubDate = FormatDate(jitem.DatePublished, time.RFC1123Z)
		}

		jauthors := jitem.Authors
		if jitem.Author != nil {
			jauthors = append([]*JSONAuthor{jitem.Author}, jitem.Authors...)
		}
		if len(jauthors) > 0 {
			item.Author = &RssAuthor{Email: jauthors[0].Name}
		}

		if len(jitem.Attachments) > 0 {
			item.Enclosure = &RssEnclosure{
				Length: "",
				Type:   jitem.Attachments[0].MimeType,
				Url:    jitem.Attachments[0].URL,
			}
		}

	}

	return ff
}

func (f *JSONFeed) ToAtom() *AtomFeed {
	ff := &AtomFeed{
		XMLName: xml.Name{
			Space: "",
			Local: "feed",
		},
		AtomCommonAttributes: AtomCommonAttributes{},
		Authors:              nil,
		Categories:           nil,
		Contributors:         nil,
		Generator:            nil,
		Icon:                 nil,
		ID: AtomId{
			AtomCommonAttributes: AtomCommonAttributes{},
			AtomUri:              "",
		},
		Links:            nil,
		Logo:             nil,
		Rights:           nil,
		Subtitle:         nil,
		Title:            nil,
		Updated:          nil,
		ExtensionElement: nil,
		Entries:          nil,
	}
	defer ff.Uniform()

	if f.Title != "" {
		ff.Title = &AtomTextConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Type:                 "",
			Text:                 f.Title,
			Div:                  nil,
		}
	}

	if f.HomePageURL != "" {
		ff.ID.AtomUri = AtomUri(f.HomePageURL)

		ff.Links = []*AtomLink{
			{
				AtomCommonAttributes: AtomCommonAttributes{},
				Href:                 AtomUri(f.HomePageURL),
				Rel:                  "",
				Type:                 "",
				Hreflang:             "",
				Title:                "",
				Length:               "",
				UndefinedContent:     nil,
			},
		}
	}

	if f.Icon != "" {
		ff.Icon = &AtomIcon{
			AtomCommonAttributes: AtomCommonAttributes{},
			AtomUri:              AtomUri(f.Icon),
		}
	}

	authors := f.Authors
	if f.Author != nil {
		authors = append([]*JSONAuthor{f.Author}, f.Authors...)
	}
	for _, author := range authors {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Name:                 author.Name,
			Uri:                  AtomUri(author.URL),
			Email:                "",
		})
	}

	ff.Language = AtomLanguageTag(f.Language)

	for _, item := range f.Items {
		entry := &AtomEntry{
			AtomCommonAttributes: AtomCommonAttributes{},
			Authors:              nil,
			Categories:           nil,
			Content:              nil,
			Contributors:         nil,
			ID:                   nil,
			Links:                nil,
			Published:            nil,
			Rights:               nil,
			Source:               nil,
			Summary:              nil,
			Title:                nil,
			Updated:              nil,
			ExtensionElement:     nil,
		}
		ff.Entries = append(ff.Entries, entry)

		if item.ID != "" {
			entry.ID = &AtomId{
				AtomCommonAttributes: AtomCommonAttributes{},
				AtomUri:              AtomUri(item.ID),
			}
		}

		if item.URL != "" {
			entry.Links = []*AtomLink{
				{
					AtomCommonAttributes: AtomCommonAttributes{},
					Href:                 AtomUri(item.URL),
					Rel:                  "",
					Type:                 "",
					Hreflang:             "",
					Title:                "",
					Length:               "",
					UndefinedContent:     nil,
				},
			}
		}

		if item.Title != "" {
			entry.Title = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Text:                 item.Title,
				Div:                  nil,
			}
		}

		if item.ContentText != "" {
			entry.Content = &AtomContent{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Src:                  nil,
				Text:                 item.ContentText,
				Div:                  nil,
				Bytes:                nil,
			}
		}

		if item.ContentHTML != "" {
			entry.Content = &AtomContent{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "html",
				Src:                  nil,
				Text:                 item.ContentHTML,
				Bytes:                nil,
			}
		}

		if item.Summary != "" {
			entry.Summary = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Text:                 item.Summary,
				Div:                  nil,
			}
		}

		if item.DatePublished != "" {
			entry.Published = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             item.DatePublished,
			}
		}

		if item.DateModified != "" {
			entry.Updated = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             item.DateModified,
			}
		}

		for i := range item.Authors {
			entry.Authors = append(entry.Authors, &AtomPersonConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				Name:                 item.Authors[i].Name,
				Uri:                  AtomUri(item.Authors[i].URL),
				Email:                "",
			})
		}

		entry.Language = AtomLanguageTag(item.Language)
	}

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
	ff := &JSONFeed{
		Version: "https://jsonfeed.org/version/1.1",
		//Title:       "",
		//HomePageURL: "",
		FeedURL: "",
		//Description: "",
		UserComment: "",
		NextURL:     "",
		Icon:        "",
		Favicon:     "",
		//Author:      nil,
		Expired: nil,
		//Items:      nil,
		Hub:        nil,
		Extensions: nil,
		//Authors:     nil,
		Language: "",
	}
	defer ff.Uniform()

	if f.Channel == nil {
		return ff
	}

	ff.Title = f.Channel.Title
	ff.HomePageURL = f.Channel.Link
	ff.Description = f.Channel.Description

	if f.Channel.WebMaster != "" {
		ff.Authors = append(ff.Authors, &JSONAuthor{
			Name:   f.Channel.WebMaster,
			URL:    "",
			Avatar: "",
		})
	}

	if f.Channel.ManagingEditor != "" {
		ff.Authors = append(ff.Authors, &JSONAuthor{
			Name:   f.Channel.ManagingEditor,
			URL:    "",
			Avatar: "",
		})
	}

	for _, item := range append(f.Items, f.Channel.Items...) {
		jitem := &JSONItem{
			//ID:            "",
			//URL:           "",
			//ExternalURL: "",
			//Title:         "",
			//ContentHTML:   "",
			//ContentText:   "",
			//Summary:       "",
			Image:       "",
			BannerImage: "",
			//DatePublished: "",
			//DateModified:  "",
			//Author:      nil,
			Tags:        nil,
			Attachments: nil,
			Extensions:  nil,
			//Authors:     nil,
			Language: "",
		}
		ff.Items = append(ff.Items, jitem)

		if item.Guid != nil {
			jitem.ID = item.Guid.Guid
		}

		jitem.URL = item.Link

		if item.Source != nil {
			jitem.ExternalURL = item.Source.Url
		}

		jitem.Title = item.Title

		if item.ContentEncoded != nil {
			jitem.ContentHTML = string(item.ContentEncoded.Content)
		}

		if item.Content != nil {
			jitem.ContentText = string(item.Content.Content)
		}

		jitem.Summary = item.Description

		jitem.DatePublished = item.PubDate
		jitem.DateModified = item.PubDate

		if item.Author != nil {
			jitem.Authors = []*JSONAuthor{
				{
					Name:   item.Author.Email,
					URL:    "",
					Avatar: "",
				},
			}
		}

		if item.Enclosure != nil {
			jitem.Attachments = []*JSONAttachments{
				{
					URL:               item.Enclosure.Url,
					MimeType:          item.Enclosure.Type,
					Title:             "",
					SizeInBytes:       0,
					DurationInSeconds: 0,
				},
			}
		}

	}

	ff.Language = f.Channel.Language

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
	defer ff.Uniform()

	if ff.Channel.Image == nil {
		ff.Channel.Image = f.Image
	}

	if ff.Channel.TextInput == nil {
		ff.Channel.TextInput = f.TextInput
	}

	ff.Channel.Items = append(f.Channel.Items, f.Items...)

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
	defer ff.Uniform()

	if f.Channel == nil {
		return ff
	}

	if f.Channel.WebMaster != "" {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Name:                 "",
			Uri:                  "",
			Email:                AtomEmailAddress(f.Channel.WebMaster),
		})
	}

	if f.Channel.ManagingEditor != "" {
		ff.Authors = append(ff.Authors, &AtomPersonConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Name:                 "",
			Uri:                  "",
			Email:                AtomEmailAddress(f.Channel.ManagingEditor),
		})
	}

	for i := range f.Channel.Categories {
		ff.Categories = append(ff.Categories, &AtomCategory{
			AtomCommonAttributes: AtomCommonAttributes{},
			Term:                 "",
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

	if f.Channel.Title != "" {
		ff.Title = &AtomTextConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			Type:                 "",
			Text:                 f.Channel.Title,
			Div:                  nil,
		}
	}

	if f.Channel.PubDate != "" {
		ff.Updated = &AtomDateConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			DateTime:             f.Channel.PubDate,
		}
	} else if f.Channel.LastBuildDate != "" {
		ff.Updated = &AtomDateConstruct{
			AtomCommonAttributes: AtomCommonAttributes{},
			DateTime:             f.Channel.LastBuildDate,
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
					Name:                 "",
					Uri:                  "",
					Email:                AtomEmailAddress(item.Author.Email),
				},
			}
		}

		for i := range item.Categories {
			ff.Categories = append(ff.Categories, &AtomCategory{
				AtomCommonAttributes: AtomCommonAttributes{},
				Term:                 "",
				Scheme:               AtomUri(item.Categories[i].Domain),
				Label:                item.Categories[i].Text,
				UndefinedContent:     nil,
			})
		}

		if item.ContentEncoded != nil {
			entry.Content = &AtomContent{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "html",
				Src:                  nil,
				Text:                 string(item.ContentEncoded.Content),
				Div:                  nil,
				Bytes:                nil,
			}
		} else if item.Content != nil {
			entry.Content = &AtomContent{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Src:                  nil,
				Text:                 string(item.Content.Content),
				Div:                  nil,
				Bytes:                nil,
			}
		}

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
				DateTime:             item.PubDate,
			}

			entry.Updated = &AtomDateConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				DateTime:             item.PubDate,
			}
		}

		if item.Description != "" {
			entry.Summary = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Text:                 item.Description,
				Div:                  nil,
			}
		}

		if item.Title != "" {
			entry.Title = &AtomTextConstruct{
				AtomCommonAttributes: AtomCommonAttributes{},
				Type:                 "",
				Text:                 item.Title,
				Div:                  nil,
			}
		}

	}

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
				Name:  "Unknown (grss)",
				Uri:   "",
				Email: "",
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

}

func (f *AtomFeed) Mime(fallback bool) string {
	if fallback {
		return AtomMimeFallback
	} else {
		return AtomMime
	}
}

func (f *AtomFeed) ToJSON() *JSONFeed {
	ff := &JSONFeed{
		Version: "https://jsonfeed.org/version/1.1",
		//Title:       "",
		//HomePageURL: "",
		FeedURL:     "",
		Description: "",
		UserComment: "",
		NextURL:     "",
		//Icon:        "",
		Favicon: "",
		//Author:     nil,
		Expired: nil,
		//Items:      nil,
		Hub:        nil,
		Extensions: nil,
		//Authors:    nil,
		//Language: "",
	}
	defer ff.Uniform()

	if f.Title != nil {
		ff.Title = f.Title.String()
	}

	if len(f.Links) > 0 {
		ff.HomePageURL = string(f.Links[0].Href)
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

	for _, entry := range f.Entries {
		item := &JSONItem{
			//ID:            "",
			//URL:           "",
			ExternalURL: "",
			//Title:         "",
			//ContentHTML:   "",
			//ContentText:   "",
			//Summary:       "",
			Image:       "",
			BannerImage: "",
			//DatePublished: "",
			//DateModified:  "",
			//Author:      nil,
			Tags:        nil,
			Attachments: nil,
			Extensions:  nil,
			//Authors:     nil,
			//Language: "",
		}
		ff.Items = append(ff.Items, item)

		if entry.ID != nil {
			item.ID = string(entry.ID.AtomUri)
		}

		if len(entry.Links) > 0 {
			item.URL = string(entry.Links[0].Href)
		}

		if entry.Title != nil {
			item.Title = entry.Title.String()
		}

		if entry.Content != nil {
			switch entry.Content.Type {
			case "html", "xhtml":
				item.ContentHTML = entry.Content.String()
			default:
				item.ContentText = entry.Content.String()
			}
		}

		if entry.Summary != nil {
			item.Summary = entry.Summary.String()
		}

		if entry.Published != nil {
			item.DatePublished = entry.Published.DateTime
		}

		if entry.Updated != nil {
			item.DateModified = entry.Updated.DateTime
		}

		for i := range entry.Authors {
			item.Authors = append(item.Authors, &JSONAuthor{
				Name: entry.Authors[i].Name,
				URL:  string(entry.Authors[i].Uri),
			})
		}

		item.Language = string(entry.Language)
	}

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
	defer ff.Uniform()

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
				Content: []byte(html.EscapeString(html.UnescapeString(entry.Content.String()))),
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
		ff.Channel.Title = f.Title.String()
	}

	if len(f.Links) > 0 {
		ff.Channel.Link = string(f.Links[0].Href)
	}

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
	defer ff.Uniform()

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
