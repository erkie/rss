package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func parseAtom(data []byte) (*Feed, error) {
	warnings := false
	feed := atomFeed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	p.Strict = false
	p.CharsetReader = charsetReader
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}

	out := new(Feed)
	out.Title = strings.TrimSpace(feed.Title)
	out.Description = strings.TrimSpace(feed.Description)
	for _, link := range feed.Links {
		if link.Rel == "alternate" || link.Rel == "" {
			out.Link = strings.TrimSpace(link.Href)
			break
		}
	}

	if feed.Items == nil {
		return nil, fmt.Errorf("Error: no feeds found in %q.", string(data))
	}

	out.Items = make([]*Item, 0, len(feed.Items))

	// Process items.
	for _, item := range feed.Items {

		next := new(Item)
		next.Title = strings.TrimSpace(item.Title)
		next.Summary = strings.TrimSpace(item.Summary)
		next.Content = strings.TrimSpace(item.Content)

		if next.Content == "" && next.Summary == "" && len(item.Media) > 0 {
			for _, media := range item.Media {
				if media.Description != "" {
					next.Content = media.Description
				}
			}
		}

		next.Date = defaultTime()
		if item.Date != "" {
			next.Date, err = parseTime(item.Date)
			if err != nil {
				return nil, err
			}
		}
		next.ID = strings.TrimSpace(item.ID)
		for _, link := range item.Links {
			if link.Rel == "alternate" || link.Rel == "" {
				next.Link = link.Href
			} else {
				next.Enclosures = append(next.Enclosures, &Enclosure{
					Url:    strings.TrimSpace(link.Href),
					Type:   link.Type,
					Length: link.Length,
				})
			}
		}
		next.Read = false
		if len(next.Link) == 0 && (strings.HasPrefix(next.ID, "http://") || strings.HasPrefix(next.ID, "https://")) {
			next.Link = next.ID
		}

		out.Items = append(out.Items, next)
	}

	if warnings && debug {
		fmt.Printf("[i] Encountered warnings:\n%s\n", data)
	}

	return out, nil
}

type atomFeed struct {
	XMLName     xml.Name   `xml:"feed"`
	Title       string     `xml:"title"`
	Description string     `xml:"subtitle"`
	Links       []atomLink `xml:"link"`
	Items       []atomItem `xml:"entry"`
	Updated     string     `xml:"updated"`
}

type atomItem struct {
	XMLName xml.Name    `xml:"entry"`
	Title   string      `xml:"title"`
	Summary string      `xml:"summary"`
	Content string      `xml:"content"`
	Links   []atomLink  `xml:"link"`
	Date    string      `xml:"updated"`
	ID      string      `xml:"id"`
	Media   []atomMedia `xml:"group"` // <media:group> from http://search.yahoo.com/mrss/
}

type atomLink struct {
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr"`
	Type   string `xml:"type,attr"`
	Length int    `xml:"length,attr"`
}

type atomMedia struct {
	Description string `xml:"description"`
}
