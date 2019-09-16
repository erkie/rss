package rss

import "strings"

type genericCategory struct {
	Term       string `xml:"term,attr"`
	TermInBody string `xml:",chardata"`
}

func fetchCategoriesFromArray(categories []genericCategory, splitBySlashes bool) []string {
	var ret []string
	for _, category := range categories {
		var theTerm string
		if category.Term == "" && category.TermInBody != "" {
			theTerm = category.TermInBody
		} else if category.Term != "" {
			theTerm = category.Term
		}

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
