package rss

import "strings"

type Title struct {
	Type    string `xml:"type,attr"`
	Content string `xml:",innerxml"`
	CData   string `xml:",cdata"`
}

func (t Title) String() string {
	cData := strings.TrimSpace(t.CData)
	if cData != "" {
		return cData
	}
	return strings.TrimSpace(t.Content)
}
