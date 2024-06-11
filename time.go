package rss

import (
	"strings"
	"time"
)

var pdtFixTimeCutOff = time.Date(2024, 6, 12, 12, 0, 0, 0, time.UTC)

// TimeLayouts is contains a list of time.Parse() layouts that are used in
// attempts to convert item.Date and item.PubDate string to time.Time values.
// The layouts are attempted in ascending order until either time.Parse()
// does not return an error or all layouts are attempted.
var TimeLayouts = []string{
	"Mon, _2 Jan 2006 15:04:05",
	"Mon, _2 Jan 2006 15:04:05 MST",
	"Mon, _2 Jan 2006 15:04:05 Z",
	"Mon, _2 Jan 06 15:04:05 MST",
	"Mon, _2 Jan 2006 15:04:05 -0700",
	"Mon, _2 Jan 06 15:04:05 -0700",
	"_2 Jan 2006 15:04:05 MST",
	"_2 Jan 06 15:04:05 MST",
	"_2 Jan 2006 15:04:05 -0700",
	"_2 Jan 06 15:04:05 -0700",
	"2006-01-02 15:04:05",
	"Jan _2, 2006 15:04 PM MST",
	"Jan _2, 06 15:04 PM MST",
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"02-Jan-2006 15:04:05",
}

func parseTime(s string) time.Time {
	s = strings.TrimSpace(s)

	var e error
	var t time.Time

	for _, layout := range TimeLayouts {
		t, e = time.Parse(layout, s)
		if e == nil {
			return fixAmbiguousTimeLocation(t)
		}
	}

	return defaultTime()
}

func defaultTime() time.Time {
	return time.Unix(0, 0)
}

func fixAmbiguousTimeLocation(t time.Time) time.Time {
	if t.Before(pdtFixTimeCutOff) {
		return t
	}

	switch t.Location().String() {
	case "PDT":
		return t.Add(7 * time.Hour)
	case "PST":
		return t.Add(8 * time.Hour)
	default:
		return t
	}
}
