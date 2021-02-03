package scraper

import "strings"

// JoinArgs - Join arguments
func JoinArgs(args []string) string {
	return strings.ReplaceAll(strings.Join(args, " "), "\"", "")
}

// FirstSearch - Get link of first search result
func FirstSearch(qry string) string {
	var link string
	res := Search(qry)
	if len(res) > 0 {
		link = res[0].FullURL
	}
	return link
}

// IsLink - Tell whether a string is a link
func IsLink(link string) bool {
	if strings.Contains(link, "/") {
		return true
	}
	return false
}
