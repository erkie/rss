package rss

import "testing"

func TestURLsDetectPaths(t *testing.T) {
	a := "https://github.com/syncthing/syncthing/releases.atom"
	b := "https://github.com/a/b/c/"

	urls := [][]string{
		[]string{a, "/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{a, "syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{b, "d/e/f/tag/v0.12.18"},
		[]string{a, "https://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{a, "http://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{a, "//github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{a, "file://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
	}

	results := []string{
		"https://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"https://github.com/syncthing/syncthing/syncthing/syncthing/releases/tag/v0.12.18",
		"https://github.com/a/b/c/d/e/f/tag/v0.12.18",
		"https://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"http://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"http://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"file://github.com/syncthing/syncthing/releases/tag/v0.12.18",
	}

	for index, baseAndURL := range urls {
		normalizedURL := normalizeURL(baseAndURL[1], baseAndURL[0])

		if normalizedURL != results[index] {
			t.Errorf("%s != %s", normalizedURL, results[index])
		}
	}
}
