package rss

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParseCategory(t *testing.T) {
	tests := map[string]string{
		"rss_0.92":              "Author,Scripting,Userland",
		"rss_1.0":               "",
		"rss_1.0_space_in_purl": "",
		"rss_2.0":               "RSS,Example",
		"rss_2.0-1":             "",
		"atom_1.0":              "weblog,german",
		"atom_1.0-1":            "",
		"utf8.response":         "",
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

		if strings.Join(feed.Categories, ",") != want {
			t.Fatalf("%s: expected %s, got %s", test, want, strings.Join(feed.Categories, ","))
		}
	}
}
