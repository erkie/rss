package rss

import "time"

type Metadata struct {
	Authors    []Author   `xml:"author"`
	Categories []Category `xml:"category"`

	MediaGroups  []MediaGroup   `xml:"group"`
	MediaContent []MediaContent `xml:"content"`

	PubDate   string `xml:"pubDate"`
	Date      string `xml:"date"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}

func (m *Metadata) PublishedDate() time.Time {
	date := defaultTime()

	if m.PubDate != "" {
		date = ParseTime(m.PubDate)
	}

	if date.IsZero() && m.Published != "" {
		date = ParseTime(m.Published)
	}

	if date.IsZero() && m.Date != "" {
		date = ParseTime(m.Date)
	}

	if date.IsZero() && m.Updated != "" {
		date = ParseTime(m.Updated)
	}

	return date
}

func (m *Metadata) CategoriesAsString() []string {
	ret := make([]string, len(m.Categories))
	for i, category := range m.Categories {
		ret[i] = category.Contents
	}
	return ret
}

func (m *Metadata) MediaContents() string {
	for _, media := range m.MediaGroups {
		if media.Description != nil {
			return media.Description.Content
		}
	}

	for _, media := range m.MediaContent {
		if media.Description != nil {
			return media.Description.Content
		}
	}
	return ""
}
