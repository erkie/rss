package rss

import (
	"encoding/xml"
	"testing"
)

func TestTitle(t *testing.T) {
	xmls := []string{
		`<title>Hello, World!</title>`,
		`<title type="text">Hello, World!</title>`,
		`<title type="html">Hello, World!</title>`,
		`<title><![CDATA[Hello, World!]]></title>`,
	}
	for _, testCase := range xmls {
		title := Title{}
		err := xml.Unmarshal([]byte(testCase), &title)
		if err != nil {
			t.Fatal(err)
		}

		if title.String() != "Hello, World!" {
			t.Errorf("expected %s, got %s", "Hello, World!", title.String())
		}
	}
}
