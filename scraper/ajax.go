package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// ScrapeAjax - Get the link to the video source of an episode
func ScrapeAjax(ajax AjaxResult) string {
	client := &http.Client{}
	url := fmt.Sprintf("https://%s%s&refer=none", ajax.Domain, ajax.Ajax)
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	// req.Header.Set("accept-langauge","en-US,en;q=0.9,zh-TW;q=0.8,zh-CN;q=0.7,zh:q=0.6")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var obj AjaxResponse
	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&obj)

	// fmt.Println(obj)

	link := obj.Source[0].File

	return link
}

// TODO: improve link searching
// GetAjax - Find the link for the Ajax
func GetAjax(episode string) AjaxResult {
	res := AjaxResult{
		Found: false,
	}

	c := getCollector()

	c.OnHTML("div.watch-drama h1", func(e *colly.HTMLElement) {
		num := strings.Split(strings.Trim(e.DOM.Text(), " \n"), " ")
		parse := num[len(num)-1]
		conv, err := strconv.ParseFloat(parse, 64)
		if err != nil {
			fmt.Println(err)
			// } else if conv == 0 {
		} else {
			res.Num = conv
		}
	})

	c.OnHTML("div.category a", func(e *colly.HTMLElement) {
		res.Name = e.DOM.Text()
	})

	c.OnHTML("div.watch_video iframe", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		index := strings.Index(src, "streaming")
		if src != "" && index != -1 {
			streaming := src[index:]
			res.Streaming = streaming
			ajax := strings.Replace(streaming, "streaming", "ajax", 1)
			res.Ajax = ajax
			index2 := strings.Index(src, "embed")
			if index2 != -1 {
				res.Found = true
				domain := src[index2:index]
				res.Domain = domain
			}
		}
	})
	c.Visit(episode)

	return res
}

// AjaxResult - The result of scraping the ajax url from an episode link
type AjaxResult struct {
	Found     bool
	Name      string
	Num       float64
	Ajax      string
	Streaming string
	Domain    string
}

// AjaxResponse - The expected JSON response from the DC Ajax endpoint
type AjaxResponse struct {
	Source   []SourceOption `json:"source"`
	SourceBK []SourceOption `json:"source_bk"`
	Track    []TrackOption  `json:"track"`
	// Advertising []string
	LinkIFrame string `json:"linkiframe"`
}

// SourceOption - Option for source of video
type SourceOption struct {
	File    string `json:"file"`
	Label   string `json:"label"`
	Default bool   `json:"default"`
	Type    string `json:"type"`
}

// TrackOption - Part of AjaxResponse
type TrackOption struct {
	File string `json:"file"`
	Kind string `json:"kind"`
}
