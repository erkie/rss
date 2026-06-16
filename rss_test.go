package rss

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseTitle(t *testing.T) {
	tests := map[string]string{
		"rss_0.92":              "Dave Winer: Grateful Dead",
		"rss_1.0":               "Golem.de",
		"rss_1.0_space_in_purl": "ECB | Swedish krona (SEK) - Euro foreign exchange reference rates",
		"rss_2.0":               "RSS Title",
		"rss_2.0-1":             "Liftoff News",
		"atom_1.0":              "Titel des Weblogs",
		"atom_1.0-1":            "Golem.de",
	}

	for test, want := range tests {
		data, err := ioutil.ReadFile("testdata/" + test)
		if err != nil {
			t.Fatalf("Reading %s: %v", test, err)
		}

		feed, err := Parse(data, ParseOptions{})
		if err != nil {
			t.Fatalf("Parsing %s: %v", test, err)
		}

		if feed.Title.String() != want {
			t.Fatalf("%s: expected %s, got %s", test, want, feed.Title)
		}
	}
}

func TestEnclosure(t *testing.T) {
	tests := map[string]Enclosure{
		"rss_1.0":   Enclosure{URL: "http://foo.bar/baz.mp3", Type: "audio/mpeg", Length: "65535"},
		"rss_2.0":   Enclosure{URL: "http://example.com/file.mp3", Type: "audio/mpeg", Length: "65535"},
		"rss_2.0-1": Enclosure{URL: "http://gdb.voanews.com/6C49CA6D-C18D-414D-8A51-2B7042A81010_cx0_cy29_cw0_w800_h450.jpg", Type: "image/jpeg", Length: "3123"},
		"atom_1.0":  Enclosure{URL: "http://example.org/audio.mp3", Type: "audio/mpeg", Length: "1234"},
	}

	for test, want := range tests {
		data, err := ioutil.ReadFile("testdata/" + test + "_enclosure")
		if err != nil {
			t.Fatalf("Reading %s: %v", test, err)
		}

		feed, err := Parse(data, ParseOptions{})
		if err != nil {
			t.Fatalf("Parsing %s: %v", test, err)
		}

		enclosureFound := false
		for _, item := range feed.Items {
			for _, enc := range item.Enclosures {
				enclosureFound = true
				if !reflect.DeepEqual(*enc, want) {
					t.Errorf("%s: expected %#v, got %#v", test, want, *enc)
				}
			}
		}
		if !enclosureFound {
			t.Errorf("No enclosures parsed in test %v", test)
		}
	}
}

func TestComments(t *testing.T) {
	tests := map[string]string{
		"rss_1.0_comments": "http://www.xul.fr/feed/RSS-1.0.html#comments",
		"rss_2.0_comments": "http://www.wikipedia.org/comments",
	}

	for test, want := range tests {
		data, err := ioutil.ReadFile("testdata/" + test)
		if err != nil {
			t.Fatalf("Reading %s: %v", test, err)
		}

		feed, err := Parse(data, ParseOptions{})
		if err != nil {
			t.Fatalf("Parsing %s: %v", test, err)
		}

		if len(feed.Items) == 0 {
			t.Fatalf("%s: no items parsed", test)
		}

		for _, item := range feed.Items {
			if item.Comments != want {
				t.Errorf("%s: expected comments %q, got %q", test, want, item.Comments)
			}
		}
	}
}

// TestCommentsWithSlashCount guards against the slash-namespaced
// <slash:comments> count clobbering the core <comments> URL. Both elements
// share the local name "comments", so without namespace-qualified parsing
// Go's decoder maps both to the same field and the count wins.
func TestCommentsWithSlashCount(t *testing.T) {
	type want struct {
		url   string
		count string
	}
	tests := map[string]want{
		"rss_2.0_wordpress_comments": {
			url:   "https://blog.example.com/a-post/#comments",
			count: "19",
		},
		"rss_1.0_slash_comments": {
			url:   "http://www.xul.fr/feed/RSS-1.0.html#comments",
			count: "42",
		},
	}

	for test, w := range tests {
		data, err := ioutil.ReadFile("testdata/" + test)
		if err != nil {
			t.Fatalf("Reading %s: %v", test, err)
		}

		feed, err := Parse(data, ParseOptions{})
		if err != nil {
			t.Fatalf("Parsing %s: %v", test, err)
		}

		if len(feed.Items) == 0 {
			t.Fatalf("%s: no items parsed", test)
		}

		for _, item := range feed.Items {
			if item.Comments != w.url {
				t.Errorf("%s: expected comments URL %q, got %q", test, w.url, item.Comments)
			}
			if item.CommentCount != w.count {
				t.Errorf("%s: expected comment count %q, got %q", test, w.count, item.CommentCount)
			}
		}
	}
}

func TestEnclosureLink(t *testing.T) {
	tests := map[string]string{
		"rss_1.0_enclosurelink":   "http://foo.bar/baz.mp3",
		"rss_2.0_enclosurelink":   "http://example.com/file.mp3",
		"rss_2.0-1_enclosurelink": "http://gdb.voanews.com/6C49CA6D-C18D-414D-8A51-2B7042A81010_cx0_cy29_cw0_w800_h450.jpg",
		"atom_1.0_enclosurelink":  "http://example.org/audio.mp3",
	}

	for test, want := range tests {
		data, err := ioutil.ReadFile("testdata/" + test)
		if err != nil {
			t.Fatalf("Reading %s: %v", test, err)
		}

		feed, err := Parse(data, ParseOptions{})
		if err != nil {
			t.Fatalf("Parsing %s: %v", test, err)
		}

		for _, item := range feed.Items {
			if item.Link != want {
				t.Errorf("Incorrect link %s != %s on %s", item.Link, want, test)
			}
		}
	}
}
