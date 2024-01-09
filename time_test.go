package rss

import (
	"testing"
	"time"
)

const customLayout = "2006-01-02T15:04Z07:00"

var (
	timeVal         = time.Date(2015, 7, 1, 9, 27, 0, 0, time.UTC)
	originalLayouts = TimeLayouts
)

func TestParseTimeUsingOnlyDefaultLayouts(t *testing.T) {
	// Positive cases
	for _, layout := range originalLayouts {
		s := timeVal.Format(layout)
		if tv := parseTime(s); !tv.Equal(timeVal) {
			t.Errorf("expected no err and times to equal, and time value %v", tv)
		}
	}

	// Negative cases
	parseTime("")
	parseTime("abc")

	custom := timeVal.Format(customLayout)
	parseTime(custom)
}

func TestParseTimeUsingCustomLayoutsPrepended(t *testing.T) {
	TimeLayouts = append([]string{customLayout}, originalLayouts...)
	custom := timeVal.Format(customLayout)
	if tv := parseTime(custom); !tv.Equal(timeVal) {
		t.Errorf("expected no err and times to equal, and time value %v", tv)
	}
	TimeLayouts = originalLayouts
}

func TestParseTimeUsingCustomLayoutsAppended(t *testing.T) {
	TimeLayouts = append(originalLayouts, customLayout)
	custom := timeVal.Format(customLayout)
	if tv := parseTime(custom); !tv.Equal(timeVal) {
		t.Errorf("expected no err and times to equal, and time value %v", tv)
	}
	TimeLayouts = originalLayouts
}

func TestParseWithTwoDigitYear(t *testing.T) {
	s := "Sun, 18 Dec 16 18:25:00 +0100"
	if tv := parseTime(s); tv.Year() != 2016 {
		t.Errorf("expected no err and year to be 2016, and year %d", tv.Year())
	}
}

func TestParser(t *testing.T) {
	s := []string{
		"2016-06-28T00:00:00",
		"Fri, 02 Sep 2022 02:38:39 PDT",
		"Tue, 14 Mar 2023 14:05:19 Z",
		"09-Jan-2024 14:00:08",
	}
	res := []time.Time{
		time.Date(2016, 6, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 9, 2, 2, 38, 39, 0, time.UTC),
		time.Date(2023, 3, 14, 14, 5, 19, 0, time.UTC),
		time.Date(2024, 1, 9, 14, 0, 8, 0, time.UTC),
	}
	for i, form := range s {
		tv := parseTime(form)
		if tv.UTC() != res[i] {
			t.Errorf("expected no err and year to be %s = %s", res[i].String(), tv.UTC().String())
		}
	}
}
