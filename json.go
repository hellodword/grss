package grss

import "encoding/json"

// v1.0 https://www.jsonfeed.org/version/1/
// v1.1 https://www.jsonfeed.org/version/1.1/

// https://help.apple.com/itc/podcasts_connect/#/itcb54353390

const (
	// JSONMime https://www.jsonfeed.org/version/1.1/#suggestions-for-publishers-a-name-suggestions-for-publishers-a
	JSONMime         = "application/feed+json"
	JSONMimeFallback = "application/json"
)

// JSONFeed JSONFeed
type JSONFeed struct {

	// v1.0 https://www.jsonfeed.org/version/1/

	// Version version (required, string) is the URL of the version of the format the feed uses. This should appear at the very top, though we recognize that not all JSON generators allow for ordering.
	Version string `json:"version,omitempty"`

	// Title title (required, string) is the name of the feed, which will often correspond to the name of the website (blog, for instance), though not necessarily.
	Title string `json:"title,omitempty"`

	// HomePageURL home_page_url (optional but strongly recommended, string) is the URL of the resource that the feed describes. This resource may or may not actually be a “home” page, but it should be an HTML page. If a feed is published on the public web, this should be considered as required. But it may not make sense in the case of a file created on a desktop computer, when that file is not shared or is shared only privately.
	HomePageURL string `json:"home_page_url,omitempty"`

	// FeedURL feed_url (optional but strongly recommended, string) is the URL of the feed, and serves as the unique identifier for the feed. As with home_page_url, this should be considered required for feeds on the public web.
	FeedURL string `json:"feed_url,omitempty"`

	// Description description (optional, string) provides more detail, beyond the title, on what the feed is about. A feed reader may display this text.
	Description string `json:"description,omitempty"`

	// UserComment user_comment (optional, string) is a description of the purpose of the feed. This is for the use of people looking at the raw JSON, and should be ignored by feed readers.
	UserComment string `json:"user_comment,omitempty"`

	// NextURL next_url (optional, string) is the URL of a feed that provides the next n items, where n is determined by the publisher. This allows for pagination, but with the expectation that reader software is not required to use it and probably won’t use it very often. next_url must not be the same as feed_url, and it must not be the same as a previous next_url (to avoid infinite loops).
	NextURL string `json:"next_url,omitempty"`

	// Icon icon (optional, string) is the URL of an image for the feed suitable to be used in a timeline, much the way an avatar might be used. It should be square and relatively large — such as 512 x 512 — so that it can be scaled-down and so that it can look good on retina displays. It should use transparency where appropriate, since it may be rendered on a non-white background.
	Icon string `json:"icon,omitempty"`

	// Favicon favicon (optional, string) is the URL of an image for the feed suitable to be used in a source list. It should be square and relatively small, but not smaller than 64 x 64 (so that it can look good on retina displays). As with icon, this image should use transparency where appropriate, since it may be rendered on a non-white background.
	Favicon string `json:"favicon,omitempty"`

	// Author author (optional, object) specifies the feed author. The author object has several members. These are all optional — but if you provide an author object, then at least one is required:
	Author *JSONAuthor `json:"author,omitempty"`

	// Expired expired (optional, boolean) says whether or not the feed is finished — that is, whether or not it will ever update again. A feed for a temporary event, such as an instance of the Olympics, could expire. If the value is true, then it’s expired. Any other value, or the absence of expired, means the feed may continue to update.
	Expired *bool `json:"expired,omitempty"`

	// Items is an array, and is required. An item includes:
	Items []*JSONItem `json:"items,omitempty"`

	// Hub Subscribing to Real-time Notifications
	Hub interface{} `json:"-"`

	// Extensions Publishers can use custom objects in JSON Feeds. Names must start with an _ character followed by a letter. Custom objects can appear anywhere in a feed.
	Extensions map[string]interface{} `json:"-"`

	// v1.1 https://www.jsonfeed.org/version/1.1/

	// Authors authors (optional, array of objects) specifies one or more feed authors. The author object has several members. These are all optional — but if you provide an author object, then at least one is required:
	Authors []*JSONAuthor `json:"authors,omitempty"`

	// Language language (optional, string) is the primary language for the feed in the format specified in RFC 5646. The value is usually a 2-letter language tag from ISO 639-1, optionally followed by a region tag. (Examples: en or en-US.)
	Language string `json:"language,omitempty"`
}

// JSONItem items is an array, and is required. An item includes:
type JSONItem struct {
	// v1.0 https://www.jsonfeed.org/version/1/

	// ID id (required, string) is unique for that item for that feed over time. If an item is ever updated, the id should be unchanged. New items should never use a previously-used id. If an id is presented as a number or other type, a JSON Feed reader must coerce it to a string. Ideally, the id is the full URL of the resource described by the item, since URLs make great unique identifiers.
	ID string `json:"id,omitempty"`

	// URL url (optional, string) is the URL of the resource described by the item. It’s the permalink. This may be the same as the id — but should be present regardless.
	URL string `json:"url,omitempty"`

	// ExternalURL external_url (very optional, string) is the URL of a page elsewhere. This is especially useful for linkblogs. If url links to where you’re talking about a thing, then external_url links to the thing you’re talking about.
	ExternalURL string `json:"external_url,omitempty"`

	// Title title (optional, string) is plain text. Microblog items in particular may omit titles.
	Title string `json:"title,omitempty"`

	// ContentHTML content_html and content_text are each optional strings — but one or both must be present. This is the HTML or plain text of the item. Important: the only place HTML is allowed in this format is in content_html. A Twitter-like service might use content_text, while a blog might use content_html. Use whichever makes sense for your resource. (It doesn’t even have to be the same for each item in a feed.)
	ContentHTML string `json:"content_html,omitempty"`

	// ContentText same as ContentHTML
	ContentText string `json:"content_text,omitempty"`

	// Summary summary (optional, string) is a plain text sentence or two describing the item. This might be presented in a timeline, for instance, where a detail view would display all of content_html or content_text.
	Summary string `json:"summary,omitempty"`

	// Image image (optional, string) is the URL of the main image for the item. This image may also appear in the content_html — if so, it’s a hint to the feed reader that this is the main, featured image. Feed readers may use the image as a preview (probably resized as a thumbnail and placed in a timeline).
	Image string `json:"image,omitempty"`

	// BannerImage banner_image (optional, string) is the URL of an image to use as a banner. Some blogging systems (such as Medium) display a different banner image chosen to go with each post, but that image wouldn’t otherwise appear in the content_html. A feed reader with a detail view may choose to show this banner image at the top of the detail view, possibly with the title overlaid.
	BannerImage string `json:"banner_image,omitempty"`

	// DatePublished date_published (optional, string) specifies the date in RFC 3339 format. (Example: 2010-02-07T14:04:00-05:00.)
	DatePublished string `json:"date_published,omitempty"`

	// DateModified date_modified (optional, string) specifies the modification date in RFC 3339 format.
	DateModified string `json:"date_modified,omitempty"`

	// Author author (optional, object) has the same structure as the top-level author. If not specified in an item, then the top-level author, if present, is the author of the item.
	Author *JSONAuthor `json:"author,omitempty"`

	// Tags tags (optional, array of strings) can have any plain text values you want. Tags tend to be just one word, but they may be anything. Note: they are not the equivalent of Twitter hashtags. Some blogging systems and other feed formats call these categories.
	Tags []string `json:"tags,omitempty"`

	// Attachments An individual item may have one or more attachments.
	Attachments []*JSONAttachments `json:"attachments,omitempty"`

	// Extensions Publishers can use custom objects in JSON Feeds. Names must start with an _ character followed by a letter. Custom objects can appear anywhere in a feed.
	Extensions map[string]interface{} `json:"-"`

	// v1.1 https://www.jsonfeed.org/version/1.1/

	// Authors authors (optional, array of objects) has the same structure as the top-level authors. If not specified in an item, then the top-level authors, if present, are the authors of the item.
	Authors []*JSONAuthor `json:"authors,omitempty"`

	// Language language (optional, string) is the language for this item, using the same format as the top-level language field. The value can be different than the primary language for the feed when a specific item is written in a different language than other items in the feed.
	Language string `json:"language,omitempty"`
}

// JSONAuthor author (optional, object) specifies the feed author. The author object has several members. These are all optional — but if you provide an author object, then at least one is required:
type JSONAuthor struct {
	// Name name (optional, string) is the author’s name.
	Name string `json:"name,omitempty"`

	// URL url (optional, string) is the URL of a site owned by the author. It could be a blog, micro-blog, Twitter account, and so on. Ideally the linked-to page provides a way to contact the author, but that’s not required. The URL could be a mailto: link, though we suspect that will be rare.
	URL string `json:"url,omitempty"`

	// Avatar avatar (optional, string) is the URL for an image for the author. As with icon, it should be square and relatively large — such as 512 x 512 — and should use transparency where appropriate, since it may be rendered on a non-white background.
	Avatar string `json:"avatar,omitempty"`
}

// JSONAttachments attachments (optional, array) lists related resources. Podcasts, for instance, would include an attachment that’s an audio or video file. Each attachment has several members:
type JSONAttachments struct {
	// URL url (required, string) specifies the location of the attachment.
	URL string `json:"url,omitempty"`

	// MimeType mime_type (required, string) specifies the type of the attachment, such as “audio/mpeg.”
	MimeType string `json:"mime_type,omitempty"`

	// Title title (optional, string) is a name for the attachment. Important: if there are multiple attachments, and two or more have the exact same title (when title is present), then they are considered as alternate representations of the same thing. In this way a podcaster, for instance, might provide an audio recording in different formats.
	Title string `json:"title,omitempty"`

	// SizeInBytes size_in_bytes (optional, number) specifies how large the file is.
	SizeInBytes uint64 `json:"size_in_bytes,omitempty"`

	// DurationInSeconds duration_in_seconds (optional, number) specifies how long it takes to listen to or watch, when played at normal speed.
	DurationInSeconds uint64 `json:"duration_in_seconds,omitempty"`
}

func (f *JSONFeed) UnmarshalJSON(b []byte) error {
	var m = map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	type inner JSONFeed
	err = json.Unmarshal(b, (*inner)(f))
	if err != nil {
		return err
	}

	f.Extensions = map[string]interface{}{}

	for k, v := range m {
		if len(k) >= 2 &&
			k[0] == '_' &&
			(('a' <= k[1] && k[1] <= 'z') || ('A' <= k[1] && k[1] <= 'Z')) {
			f.Extensions[k] = v
		}
	}

	switch items := m["items"].(type) {
	case []interface{}:
		{
			if len(items) != len(f.Items) {
				break
			}
			for i := range f.Items {
				f.Items[i].Extensions = map[string]interface{}{}
				switch item := items[i].(type) {
				case map[string]interface{}:
					{
						for k, v := range item {
							if len(k) >= 2 &&
								k[0] == '_' &&
								(('a' <= k[1] && k[1] <= 'z') || ('A' <= k[1] && k[1] <= 'Z')) {
								f.Items[i].Extensions[k] = v
							}
						}
					}
				}
			}
		}
	}

	return nil
}
