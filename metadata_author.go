package rss

type Author struct {
	NameAttribute string `xml:"name"`
	URI           string `xml:"uri"`
	Content       string `xml:",chardata"`
}

func (a Author) Name() string {
	if a.NameAttribute != "" {
		return a.NameAttribute
	}

	return a.Content
}
