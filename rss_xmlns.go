package grss

import "encoding/xml"

func (f *RssFeed) PatchXmlns() {

	var getPrefix = func(space string) string {
		if space == "" {
			return ""
		}
		//if space == xmlns {
		//	return xmlnsPrefix
		//}
		for i := range f.Attributes {
			if f.Attributes[i].Name.Space == xmlnsPrefix && f.Attributes[i].Value == space {
				return f.Attributes[i].Name.Local
			}
		}
		for i := range f.Channel.Attributes {
			if f.Channel.Attributes[i].Name.Space == xmlnsPrefix && f.Channel.Attributes[i].Value == space {
				return f.Channel.Attributes[i].Name.Local
			}
		}
		return ""
	}

	if f.XMLName.Space != "" {
		prefix := getPrefix(f.XMLName.Space)
		if prefix != "" {
			f.XMLName.Space = ""
			f.XMLName.Local = prefix + ":" + f.XMLName.Local
		}
	}

	for i := range f.Channel.ExtensionElement {
		prefix := getPrefix(f.Channel.ExtensionElement[i].XMLName.Space)
		if prefix == "" {
			continue
		}
		f.Channel.ExtensionElement[i].XMLName.Space = ""
		f.Channel.ExtensionElement[i].XMLName.Local = prefix + ":" + f.Channel.ExtensionElement[i].XMLName.Local
	}

	for i := range f.Attributes {
		if f.Attributes[i].Name.Space == xmlnsPrefix {
			f.Attributes[i].Name.Local = xmlnsPrefix + ":" + f.Attributes[i].Name.Local
			f.Attributes[i].Name.Space = ""
		}
	}

	for i := range f.Channel.Items {
		if f.Channel.Items[i].ContentUnmarshal == nil {
			continue
		}
		prefix := getPrefix(f.Channel.Items[i].ContentUnmarshal.XMLName.Space)
		if prefix != "content" {
			continue
		}
		f.Channel.Items[i].Content = &RssContent{
			XMLName: xml.Name{
				Space: "",
				Local: "content:encoded",
			},
			Content: f.Channel.Items[i].ContentUnmarshal.Content,
		}
		f.Channel.Items[i].ContentUnmarshal = nil
	}

	for i := range f.Channel.Attributes {
		if f.Channel.Attributes[i].Name.Space == xmlnsPrefix {
			f.Channel.Attributes[i].Name.Local = xmlnsPrefix + ":" + f.Channel.Attributes[i].Name.Local
			f.Channel.Attributes[i].Name.Space = ""
		}
	}

}
