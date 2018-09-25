package rss

import "testing"

func TestURLsDetectPaths(t *testing.T) {
	finalURL := "https://gitlab.com/test"
	a := "https://github.com/syncthing/syncthing/releases.atom"
	b := "https://github.com/a/b/c/"

	urls := [][]string{
		[]string{finalURL, a, "/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, a, "syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, b, "d/e/f/tag/v0.12.18"},
		[]string{finalURL, a, "https://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, a, "http://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, a, "//github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, a, "file://github.com/syncthing/syncthing/releases/tag/v0.12.18"},
		[]string{finalURL, "", "/hello/world"},
	}

	results := []string{
		"https://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"https://github.com/syncthing/syncthing/syncthing/syncthing/releases/tag/v0.12.18",
		"https://github.com/a/b/c/d/e/f/tag/v0.12.18",
		"https://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"http://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"http://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"file://github.com/syncthing/syncthing/releases/tag/v0.12.18",
		"https://gitlab.com/hello/world",
	}

	for index, baseAndURL := range urls {
		normalizedURL := normalizeURL(baseAndURL[2], baseAndURL[1], baseAndURL[0])

		if normalizedURL != results[index] {
			t.Errorf("%s != %s", normalizedURL, results[index])
		}
	}
}
