package rss

import (
	"os"
	"testing"
)

func TestAuthors(t *testing.T) {
	file, _ := os.ReadFile("testdata/atom_1.0-1")

	out, err := Parse(file, ParseOptions{})
	if err != nil {
		t.Error(err)
	}

	expectedAuthors := []string{
		"Jörg Thoma",
		"Achim Sawall",
		"Hanno Böck",
		"Jörg Thoma",
	}

	for index, expectedAuthor := range expectedAuthors {
		item := out.Items[index]

		if len(item.Metadata.Authors) != 1 {
			t.Errorf("%d contained %d authors", index, len(item.Metadata.Authors))
		}

		if expectedAuthor != item.Metadata.Authors[0].Name() {
			t.Errorf("%d expected %s but got %s", index, expectedAuthor, item.Metadata.Authors[0])
		}
	}
}
