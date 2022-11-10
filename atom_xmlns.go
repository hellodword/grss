package grss

// https://github.com/golang/go/issues/9519

const (
	xmlnsPrefix = "xmlns"
)

func (f *AtomFeed) PatchXmlns() {
	//var xmlns string
	//for i := range f.UndefinedAttribute {
	//	if f.UndefinedAttribute[i].Name.Local == xmlnsPrefix {
	//		xmlns = f.UndefinedAttribute[i].Value
	//		break
	//	}
	//}

	var getPrefix = func(space string) string {
		if space == "" {
			return ""
		}
		//if space == xmlns {
		//	return xmlnsPrefix
		//}
		for i := range f.UndefinedAttribute {
			if f.UndefinedAttribute[i].Name.Space == xmlnsPrefix && f.UndefinedAttribute[i].Value == space {
				return f.UndefinedAttribute[i].Name.Local
			}
		}
		return ""
	}

	for i := range f.ExtensionElement {
		prefix := getPrefix(f.ExtensionElement[i].XMLName.Space)
		if prefix == "" {
			continue
		}
		f.ExtensionElement[i].XMLName.Space = ""
		f.ExtensionElement[i].XMLName.Local = prefix + ":" + f.ExtensionElement[i].XMLName.Local
	}

	for i := range f.Entries {
		for j := range f.Entries[i].ExtensionElement {
			prefix := getPrefix(f.Entries[i].ExtensionElement[j].XMLName.Space)
			if prefix == "" {
				continue
			}
			f.Entries[i].ExtensionElement[j].XMLName.Space = ""
			f.Entries[i].ExtensionElement[j].XMLName.Local = prefix + ":" + f.Entries[i].ExtensionElement[j].XMLName.Local
		}
	}

	for i := range f.UndefinedAttribute {
		if f.UndefinedAttribute[i].Name.Space == xmlnsPrefix {
			f.UndefinedAttribute[i].Name.Local = xmlnsPrefix + ":" + f.UndefinedAttribute[i].Name.Local
			f.UndefinedAttribute[i].Name.Space = ""
		}
	}
}
