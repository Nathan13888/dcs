package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// DramaInfo - Information about a drama
type DramaInfo struct {
	FullURL string `json:"fullurl"`
	SubURL  string `json:"suburl"`
	Domain  string `json:"domain"`
	Name    string `json:"name"`
}

// Search - Search for something...
func Search(qry string) []DramaInfo {
	qe := url.QueryEscape(qry)
	var url string
	var nameQry string
	if ASIANLOAD {
		url = fmt.Sprintf("%s/search.html?&keyword=%s", URL, qe)
		nameQry = "div.name"
	} else {
		url = fmt.Sprintf("%s/search?type=movies&keyword=%s", URL, qe)
		nameQry = "h3.title"
	}

	res := []DramaInfo{}

	c := getCollector()

	// check if there were results
	c.OnHTML(nameQry, func(e *colly.HTMLElement) {
		name := strings.Trim(e.DOM.Contents().Text(), " \n")
		if len(name) > 10 && ASIANLOAD { // trim away " Episode X" from the end
			end := len(name)
			for i := len(name) - 1; i >= 0; i-- { // remove all numbers
				c := name[i]
				if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
					break
				}
				if ('0' <= c && c <= '9') || (c == ' ') {
					end--
				}
			}
			if end >= 8 && strings.EqualFold(name[end-7:end], "episode") {
				end -= 8
			}
			name = name[:end]
		}

		subURL, _ := e.DOM.Parent().Attr("href")
		fullURL := e.Request.AbsoluteURL(subURL)

		// TODO: filter nil results
		// if name == "" || subURL || fullURL {
		// 	panic("")
		// }
		var obj = DramaInfo{
			FullURL: fullURL,
			SubURL:  subURL,
			Domain:  URL,
			Name:    EscapeName(name),
		}
		res = append(res, obj)
	})
	c.Visit(url)

	// printObj(res)
	// fmt.Println(res)

	return res
}
