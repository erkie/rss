package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

func parseRSS2(data []byte, options ParseOptions) (*Feed, error) {
	warnings := false
	feed := rss2_0Feed{}
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
	if feed.Channel == nil {
		return nil, fmt.Errorf("Error: no channel found in %q", string(data))
	}
	channel := feed.Channel

	out := new(Feed)
	out.Type = "rss2.0"
	out.Title = strings.TrimSpace(channel.Title)
	out.Description = strings.TrimSpace(channel.Description)
	out.Categories = fetchCategoriesFromArray(channel.Categories, true)
	for _, link := range channel.Links {
		if link.Rel == "alternate" || link.Rel == "" {
			if link.Href == "" && link.Contents != "" {
				out.Link = strings.TrimSpace(link.Contents)
			} else {
				out.Link = strings.TrimSpace(link.Href)
			}
			break
		}
	}

	itemsToUse := channel.Items
	if itemsToUse == nil {
		itemsToUse = feed.Items
	}

	if itemsToUse == nil {
		itemsToUse = make([]rss2_0Item, 0)
	}

	out.Items = make([]*Item, 0, len(itemsToUse))

	// Process items.
	for _, item := range itemsToUse {

		if item.ID == "" {
			if len(item.Links) == 0 && len(item.Description) == 0 && len(item.Content) == 0 {
				if debug {
					fmt.Printf("[w] Item %q has no ID or link and will be ignored.\n", item.Title)
					fmt.Printf("[w] %#v\n", item)
				}
				warnings = true
				continue
			}
		}

		next := new(Item)

		if item.Links != nil {
			for _, link := range item.Links {
				if link != "" {
					next.Link = strings.TrimSpace(link)
					break
				}
			}
		}

		if next.Link == "" && len(item.Href) > 0 {
			next.Link = strings.TrimSpace(item.Href)
		}

		next.Title = strings.TrimSpace(item.Title)
		next.Summary = strings.TrimSpace(item.Description)
		next.Content = strings.TrimSpace(item.Content)

		if next.Content == "" && len(item.Media) > 0 {
			for _, media := range item.Media {
				if media.Description != "" {
					next.Content = media.Description
				}
			}
		}

		next.Date = defaultTime()
		if item.Date != "" {
			next.Date = parseTime(item.Date)
		} else if item.PubDate != "" {
			next.Date = parseTime(item.PubDate)
		}
		next.ID = strings.TrimSpace(item.ID)
		if len(item.Enclosures) > 0 {
			next.Enclosures = make([]*Enclosure, len(item.Enclosures))
			for i := range item.Enclosures {
				next.Enclosures[i] = item.Enclosures[i].Enclosure()
			}
		}
		if len(next.Link) == 0 && (strings.HasPrefix(next.ID, "http://") || strings.HasPrefix(next.ID, "https://")) {
			next.Link = next.ID
		}
		if len(next.Link) == 0 && len(next.Enclosures) > 0 {
			next.Link = next.Enclosures[0].URL
		}

		out.Items = append(out.Items, next)
	}

	out.Links = make([]*Link, len(channel.Links))
	for i, link := range channel.Links {
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

type rss2_0Feed struct {
	XMLName xml.Name       `xml:"rss"`
	Channel *rss2_0Channel `xml:"channel"`
	Items   []rss2_0Item   `xml:"item"`
}

type rss2_0Channel struct {
	XMLName     xml.Name          `xml:"channel"`
	Title       string            `xml:"title"`
	Description string            `xml:"description"`
	Links       []rss2_0Link      `xml:"link"`
	Items       []rss2_0Item      `xml:"item"`
	MinsToLive  string            `xml:"ttl"`
	SkipHours   []string          `xml:"skipHours>hour"`
	SkipDays    []string          `xml:"skipDays>day"`
	Categories  []genericCategory `xml:"category"`
}

type rss2_0Item struct {
	XMLName     xml.Name          `xml:"item"`
	Title       string            `xml:"title"`
	Description string            `xml:"description"`
	Content     string            `xml:"encoded"`
	Category    string            `xml:"category"`
	Links       []string          `xml:"link"`
	Href        string            `xml:"href"` // Non-standard but found in the wild...
	PubDate     string            `xml:"pubDate"`
	Date        string            `xml:"date"`
	ID          string            `xml:"guid"`
	Enclosures  []rss2_0Enclosure `xml:"enclosure"`
	Media       []rss2_0Media     `xml:"group"` // <media:group> from http://search.yahoo.com/mrss/
}

type rss2_0Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Type    string   `xml:"type,attr"`
	Length  string   `xml:"length,attr"`
}

type rss2_0Media struct {
	Description string `xml:"description"`
}

type rss2_0Link struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Length   string `xml:"length,attr"`
	Contents string `xml:",chardata"`
}

func (r *rss2_0Enclosure) Enclosure() *Enclosure {
	out := new(Enclosure)
	out.URL = strings.TrimSpace(r.URL)
	out.Type = r.Type
	out.Length = r.Length
	return out
}
