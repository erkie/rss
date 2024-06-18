package rss

import (
	"strings"
	"time"
)

var pdtFixTimeCutOff = time.Date(2024, 6, 19, 10, 0, 0, 0, time.UTC)

// This map contains the incorrect timezone offset in seconds for the timezone
// Based on timezones we've seen in the wild
var timeZoneMap = map[string]int{
	"PDT":  -7 * 60 * 60,  // -7 hours in seconds
	"PST":  -8 * 60 * 60,  // -8 hours in seconds
	"EDT":  -4 * 60 * 60,  // -4 hours in seconds
	"EST":  -5 * 60 * 60,  // -5 hours in seconds
	"CDT":  -5 * 60 * 60,  // -5 hours in seconds
	"CST":  -6 * 60 * 60,  // -6 hours in seconds
	"MDT":  -6 * 60 * 60,  // -6 hours in seconds
	"MST":  -7 * 60 * 60,  // -7 hours in seconds
	"GMT":  0 * 60 * 60,   // 0 hours in seconds
	"BST":  1 * 60 * 60,   // 1 hour in seconds
	"AKDT": -8 * 60 * 60,  // -8 hours in seconds
	"AKST": -9 * 60 * 60,  // -9 hours in seconds
	"HADT": -9 * 60 * 60,  // -9 hours in seconds
	"HAST": -10 * 60 * 60, // -10 hours in seconds
}

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

// This post explains why this is needed: https://stackoverflow.com/a/66606191/224732
func fixAmbiguousTimeLocation(t time.Time) time.Time {
	if t.Before(pdtFixTimeCutOff) {
		return t
	}

	offset, exists := timeZoneMap[t.Location().String()]
	if !exists {
		return t
	}

	localTime := t.Add(-time.Duration(offset) * time.Second)
	return localTime
}
