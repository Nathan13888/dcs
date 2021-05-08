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
	url := fmt.Sprintf("%s/search?type=movies&keyword=%s", URL, url.QueryEscape(qry))

	res := []DramaInfo{}

	c := getCollector()

	// check if there were results
	c.OnHTML("h3.title", func(e *colly.HTMLElement) {
		name := strings.Trim(e.DOM.Contents().Text(), " \n")

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
			Name:    name,
		}
		res = append(res, obj)
	})
	c.Visit(url)

	// printObj(res)
	// fmt.Println(res)

	return res
}
