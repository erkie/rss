package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

func parseRSS1(data []byte, options ParseOptions) (*Feed, error) {
	warnings := false
	feed := rss1_0Feed{}
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
	out.Type = "rss1.0"
	out.Title = channel.Title
	out.Description = channel.Description
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

	if feed.Items == nil {
		feed.Items = make([]rss1_0Item, 0)
	}

	out.Items = make([]*Item, 0, len(feed.Items))

	for _, item := range feed.Items {
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

		next := &Item{
			ID:      strings.TrimSpace(item.ID),
			Title:   item.Title,
			Summary: strings.TrimSpace(item.Description),
			Content: strings.TrimSpace(item.Content),
		}

		if item.Links != nil {
			for _, link := range item.Links {
				if link != "" {
					next.Link = strings.TrimSpace(link)
					break
				}
			}
		}

		next.SetMetadata(item.Metadata)

		if next.ID == "" && item.RDFAbout != "" {
			next.ID = strings.TrimSpace(item.RDFAbout)
			if next.Meta == nil {
				next.Meta = make(map[string]string)
			}
			next.Meta["id_from_rdf_about"] = "1"
		}

		if len(item.Enclosures) > 0 {
			next.Enclosures = make([]*Enclosure, len(item.Enclosures))
			for i := range item.Enclosures {
				next.Enclosures[i] = item.Enclosures[i].Enclosure()
			}
		}

		next.EnsureLinks()

		out.Items = append(out.Items, next)
	}

	out.Links = make([]*Link, len(channel.Links))
	for i, link := range channel.Links {
		out.Links[i] = &Link{
			URL: link.Href,
			Rel: link.Rel,
		}
	}

	mapItemsBySequence(out, feed)

	if warnings && debug {
		fmt.Printf("[i] Encountered warnings:\n%s\n", data)
	}

	return out, nil
}

func mapItemsBySequence(out *Feed, feed rss1_0Feed) {
	if len(out.Items) != len(feed.Channel.Sequence) || len(feed.Channel.Sequence) == 0 {
		log.Println("Not equal", len(out.Items), len(feed.Channel.Sequence))
		return
	}

	byID := make(map[string]*Item)

	for _, item := range out.Items {
		if item.ID == "" {
			log.Println("1. aborting because invalid item.ID")
			return
		}
		if _, exists := byID[item.ID]; exists {
			log.Println("2. aborting because it already exists, causing duplicate sequence")
			return
		}
		byID[item.ID] = item
	}

	newItems := make([]*Item, len(out.Items))
	for index, sequenceItem := range feed.Channel.Sequence {
		if targetItem, ok := byID[sequenceItem.Resource]; ok {
			newItems[index] = targetItem
		} else {
			log.Printf("3. Target item for resource %s not found", sequenceItem.Resource)
			return
		}
	}
	out.Items = newItems
}

type rss1_0Feed struct {
	XMLName xml.Name       `xml:"RDF"`
	Channel *rss1_0Channel `xml:"channel"`
	Items   []rss1_0Item   `xml:"item"`
}

type rss1_0Channel struct {
	XMLName     xml.Name         `xml:"channel"`
	Title       Title            `xml:"title"`
	Description string           `xml:"description"`
	Links       []rss1_0Link     `xml:"link"`
	MinsToLive  string           `xml:"ttl"`
	SkipHours   []string         `xml:"skipHours>hour"`
	SkipDays    []string         `xml:"skipDays>day"`
	Categories  []Category       `xml:"category"`
	Sequence    []rss1_0Sequence `xml:"items>Seq>li"`
}

type rss1_0Sequence struct {
	XMLName  xml.Name `xml:"li"`
	Resource string   `xml:"resource,attr"`
}

type rss1_0Item struct {
	XMLName     xml.Name          `xml:"item"`
	Title       Title             `xml:"title"`
	Description string            `xml:"description"`
	Content     string            `xml:"encoded"`
	Links       []string          `xml:"link"`
	ID          string            `xml:"guid"`
	RDFAbout    string            `xml:"about,attr"`
	Enclosures  []rss1_0Enclosure `xml:"enclosure"`

	Metadata
}

type rss1_0Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"resource,attr"`
	Type    string   `xml:"type,attr"`
	Length  string   `xml:"length,attr"`
}

type rss1_0Link struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	Length   string `xml:"length,attr"`
	Contents string `xml:",chardata"`
}

func (r *rss1_0Enclosure) Enclosure() *Enclosure {
	out := new(Enclosure)
	out.URL = r.URL
	out.Type = r.Type
	out.Length = r.Length
	return out
}
