package rss

import (
	"net/url"
	"regexp"
	"strings"
)

func normalizeURLsInFeed(feed *Feed) {
	for _, item := range feed.Items {
		item.Link = normalizeURL(item.Link, feed.Link)
	}
}

func normalizeURL(link string, base string) string {
	// Don't touch regular http links
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	}
	// Add http to protocol independent links. Client can deal with http vs https
	if strings.HasPrefix(link, "//") {
		return "http:" + link
	}
	// Find other protocols. Maybe we should limit to certain allowed protocol. This being the web HTTP makes sense...
	pieces := strings.SplitN(link, ":", 2)

	if len(pieces) == 2 && (pieces[0] == "http" || pieces[0] == "https" || isValidScheme(pieces[0])) {
		return link
	}
	// Parse base url and use that for link
	baseURL, err := url.Parse(base)
	// Could not parse baseURL, not much we can do but return original
	if err != nil {
		return link
	}

	baseHost := baseURL.Scheme + "://" + baseURL.Host

	// Simple case, URL has leading slash, otherwise we need to resolve the path
	if strings.HasPrefix(link, "/") {
		return baseHost + link
	}

	basePath := baseURL.EscapedPath()
	if !strings.HasSuffix(basePath, "/") {
		slashPieces := strings.Split(basePath, "/")
		basePath = strings.Join(slashPieces[0:len(slashPieces)-1], "/") + "/"
	}

	builtURL, err := url.Parse(baseHost + basePath + link)

	// Our attempt at building a URL failed, yolo, and try simple case
	if err != nil {
		return baseHost + link
	}

	return builtURL.String()
}

var schemeRegex *regexp.Regexp

func isValidScheme(s string) bool {
	if schemeRegex == nil {
		schemeRegex = regexp.MustCompile("^[A-Za-z][A-Za-z0-9+-.]+$")
	}

	// scheme = alpha *( alpha | digit | "+" | "-" | "." )
	return schemeRegex.MatchString(s)
}
