package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func parseAtom(data []byte, options ParseOptions) (*Feed, error) {
	warnings := false
	feed := atomFeed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	p.Strict = false
	p.CharsetReader = CharsetReader
	if options.CharsetReader != nil {
		p.CharsetReader = options.CharsetReader
	}
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}

	out := new(Feed)
	out.Type = "atom"
	out.Title = feed.Title
	out.Description = strings.TrimSpace(feed.Description)
	for _, link := range feed.Links {
		if link.Rel == "alternate" || link.Rel == "" {
			if link.Href == "" && link.Contents != "" {
				out.Link = strings.TrimSpace(link.Contents)
			} else {
				out.Link = strings.TrimSpace(link.Href)
			}
			break
		}
	}
	out.Categories = fetchCategoriesFromArray(feed.Categories, false)

	if feed.Items == nil {
		feed.Items = make([]atomItem, 0)
	}

	out.Items = make([]*Item, 0, len(feed.Items))

	for _, item := range feed.Items {
		next := &Item{
			ID:      strings.TrimSpace(item.ID),
			Title:   item.Title,
			Summary: strings.TrimSpace(item.Summary.String()),
			Content: strings.TrimSpace(item.Content.String()),
		}

		if next.Content == "" {
			next.Content = strings.TrimSpace(item.Description)
		}

		next.SetMetadata(item.Metadata)

		for _, link := range item.Links {
			if link.Rel == "alternate" || link.Rel == "" {
				next.Link = link.Href
			} else {
				next.Enclosures = append(next.Enclosures, &Enclosure{
					URL:    strings.TrimSpace(link.Href),
					Type:   link.Type,
					Length: link.Length,
				})
			}
		}

		next.EnsureLinks()

		out.Items = append(out.Items, next)
	}

	out.Links = make([]*Link, len(feed.Links))
	for i, link := range feed.Links {
		out.Links[i] = &Link{
			URL: link.Href,
			Rel: link.Rel,
		}
	}

	if warnings && debug {
		fmt.Printf("[i] Encountered warnings:\n%s\n", data)
	}

	return out, nil
}

type atomFeed struct {
	XMLName     xml.Name   `xml:"feed"`
	Title       Title      `xml:"title"`
	Description string     `xml:"subtitle"`
	Links       []atomLink `xml:"link"`
	Items       []atomItem `xml:"entry"`
	Updated     string     `xml:"updated"`
	Categories  []Category `xml:"category"`
}

type atomItem struct {
	XMLName     xml.Name    `xml:"entry"`
	Title       Title       `xml:"title"`
	Summary     atomContent `xml:"summary"`
	Content     atomContent `xml:"content"`
	Description string      `xml:"description"`
	Links       []atomLink  `xml:"link"`

	ID string `xml:"id"`

	Metadata
}

type atomContent struct {
	Content string `xml:",innerxml"`
	CData   string `xml:",cdata"`
}

func (a atomContent) String() string {
	cData := strings.TrimSpace(a.CData)
	if cData != "" {
		return cData
	}
	return strings.TrimSpace(a.Content)
}

type atomLink struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Length   string `xml:"length,attr"`
	Contents string `xml:",chardata"`
}
