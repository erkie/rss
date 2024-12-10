package rss

import "strings"

type Categories struct {
	Category `xml:"category"`
}

type Category struct {
	Term     string `xml:"term,attr"`
	Contents string `xml:",chardata"`
}

func (c Category) Name() string {
	if c.Term == "" {
		return c.Contents
	} else {
		return c.Term
	}
}

func fetchCategoriesFromArray(categories []Category, splitBySlashes bool) []string {
	var ret []string
	for _, category := range categories {
		theTerm := category.Name()

		if theTerm != "" {
			if splitBySlashes {
				pieces := strings.Split(theTerm, "/")
				ret = append(ret, pieces...)
			} else {
				ret = append(ret, theTerm)
			}
		}
	}
	for index, val := range ret {
		ret[index] = strings.TrimSpace(val)
	}
	return ret
}
