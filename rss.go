package rss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"
)

type ParseOptions struct {
	ResponseHeaders http.Header
	FinalURL        string
}

// ParserFunc is the interface for a parser
type ParserFunc func(data []byte, options ParseOptions) (*Feed, error)

// Parse RSS or Atom data.
func Parse(data []byte, options ParseOptions) (*Feed, error) {
	data = DiscardInvalidUTF8IfUTF8(data, options.ResponseHeaders)

	var feed *Feed
	var err error

	possibleParsers := make([]ParserFunc, 0)

	if strings.Contains(string(data), "\"http://purl.org/rss/1.0/\"") || strings.Contains(string(data), "<rdf:RDF") {
		if debug {
			fmt.Println("[i] Parsing as RSS 1.0")
		}
		possibleParsers = append(possibleParsers, parseRSS1)
	}

	if strings.Contains(string(data), "<rss") {
		if debug {
			fmt.Println("[i] Parsing as RSS 2.0")
		}
		possibleParsers = append(possibleParsers, parseRSS2)
		feed, err = parseRSS2(data, options)
	}

	if debug {
		fmt.Println("[i] Parsing as Atom")
	}
	possibleParsers = append(possibleParsers, parseAtom)

	for _, parser := range possibleParsers {
		feed, err = parser(data, options)

		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	normalizeURLsInFeed(feed, options.FinalURL)

	return feed, err
}

// Feed is the top-level structure.
type Feed struct {
	Type        string
	Title       string
	Description string
	Link        string // Link to the creator's website.
	UpdateURL   string // URL of the feed itself.
	Items       []*Item
	Links       []*Link
	Categories  []string
}

// Link as defined inside RSS feeds that can contain various information
type Link struct {
	URL string
	Rel string
}

func (f *Feed) String() string {
	buf := new(bytes.Buffer)
	if debug {
		w := tabwriter.NewWriter(buf, 0, 8, 0, '\t', tabwriter.StripEscape)
		fmt.Fprintf(w, "Feed {\n")
		fmt.Fprintf(w, "\xff\t\xffType:\t%q\n", f.Type)
		fmt.Fprintf(w, "\xff\t\xffTitle:\t%q\n", f.Title)
		fmt.Fprintf(w, "\xff\t\xffDescription:\t%q\n", f.Description)
		fmt.Fprintf(w, "\xff\t\xffLink:\t%q\n", f.Link)
		fmt.Fprintf(w, "\xff\t\xffUpdateURL:\t%q\n", f.UpdateURL)
		fmt.Fprintf(w, "\xff\t\xffItems:\t(%d) {\n", len(f.Items))
		for _, item := range f.Items {
			fmt.Fprintf(w, "%s\n", item.Format(2))
		}
		fmt.Fprintf(w, "\xff\t\xff}\n}\n")
		w.Flush()
	} else {
		w := buf
		fmt.Fprintf(w, "Feed %q\n", f.Title)
		fmt.Fprintf(w, "\t%q\n", f.Description)
		fmt.Fprintf(w, "\t%q\n", f.Link)
		fmt.Fprintf(w, "\tItems:\n")
		for _, item := range f.Items {
			fmt.Fprintf(w, "\t%s\n", item.Format(2))
		}
	}
	return buf.String()
}

// Item represents a single story.
type Item struct {
	Title      string            `json:"title"`
	Summary    string            `json:"summary"`
	Content    string            `json:"content"`
	Category   string            `json:"category"`
	Link       string            `json:"link"`
	Date       time.Time         `json:"date"`
	ID         string            `json:"id"`
	Enclosures []*Enclosure      `json:"enclosures"`
	Meta       map[string]string `json:"meta"`
}

func (i *Item) String() string {
	return i.Format(0)
}

// Format format an item nicely
func (i *Item) Format(indent int) string {
	buf := new(bytes.Buffer)
	single := strings.Repeat("\t", indent)
	double := single + "\t"
	if debug {
		w := tabwriter.NewWriter(buf, 0, 8, 0, '\t', tabwriter.StripEscape)
		fmt.Fprintf(w, "\xff%s\xffItem {\n", single)
		fmt.Fprintf(w, "\xff%s\xffTitle:\t%q\n", double, i.Title)
		fmt.Fprintf(w, "\xff%s\xffSummary:\t%q\n", double, i.Summary)
		fmt.Fprintf(w, "\xff%s\xffCategory:\t%q\n", double, i.Category)
		fmt.Fprintf(w, "\xff%s\xffLink:\t%s\n", double, i.Link)
		fmt.Fprintf(w, "\xff%s\xffDate:\t%s\n", double, i.Date.Format(DATE))
		fmt.Fprintf(w, "\xff%s\xffID:\t%s\n", double, i.ID)
		fmt.Fprintf(w, "\xff%s\xffContent:\t%q\n", double, i.Content)
		fmt.Fprintf(w, "\xff%s\xff}\n", single)
		w.Flush()
	} else {
		w := buf
		fmt.Fprintf(w, "%sItem %q\n", single, i.Title)
		fmt.Fprintf(w, "%s%q\n", double, i.Link)
		fmt.Fprintf(w, "%s%s\n", double, i.Date.Format(DATE))
		fmt.Fprintf(w, "%s%q\n", double, i.ID)
		fmt.Fprintf(w, "%s%q\n", double, i.Content)
	}
	return buf.String()
}

// Enclosure holds enclosure data
type Enclosure struct {
	URL    string `json:"url"`
	Type   string `json:"type"`
	Length string `json:"length"`
}

// Get returns an io.Reader for the data held by the Enclosure
func (e *Enclosure) Get() (io.ReadCloser, error) {
	if e == nil || e.URL == "" {
		return nil, errors.New("No enclosure")
	}

	res, err := http.Get(e.URL)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
