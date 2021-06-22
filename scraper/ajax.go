package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// TODO: implement (string,error) return
// ScrapeAjax - Get the link to the video source of an episode
func ScrapeAjax(ajax AjaxResult) string {
	client := &http.Client{}
	url := fmt.Sprintf("https://%s%s&refer=none", ajax.Domain, ajax.Ajax)
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	req.Header.Set("User-Agent", USERAGENT)
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	// req.Header.Set("accept-langauge","en-US,en;q=0.9,zh-TW;q=0.8,zh-CN;q=0.7,zh:q=0.6")
	// req.Header.Set("sec-fetch-mode", "cors")
	// req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	var obj AjaxResponse
	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&obj)

	if len(obj.Source) < 1 {
		fmt.Println("INVALID AJAX FOUND FROM", url)
		fmt.Println(obj)
	}

	link := obj.Source[0].File

	return link
}

func GetInfo(episode string) (string, float64, string) {
	var name string
	var episodeNum float64
	var streaming string

	var epQry string
	var nameQry string
	var iframeQry string
	if ASIANLOAD {
		epQry = "div.video-info-left h1"
		nameQry = "div.video-details span"
		iframeQry = "div.play-video iframe"
	} else {
		epQry = "div.watch-drama h1"
		nameQry = "div.category a"
		iframeQry = "div.watch_video iframe"
	}

	c := getCollector()

	// Find episode number
	c.OnHTML(epQry, func(e *colly.HTMLElement) {
		num := strings.Split(strings.Trim(e.DOM.Text(), " \n"), " ")
		var parse string
		if ASIANLOAD {
			parse = num[len(num)-3]
		} else {
			parse = num[len(num)-1]
		}
		conv, err := strconv.ParseFloat(parse, 64)
		if err != nil {
			fmt.Println(err)
			// res.Num = -1
		} else {
			episodeNum = conv
		}
	})

	c.OnHTML(nameQry, func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	c.OnHTML(iframeQry, func(e *colly.HTMLElement) {
		src := strings.Trim(e.Attr("src"), " /")
		index := strings.Index(src, "streaming")
		if len(src) > 0 {
			streaming = src[index:]
		}
	})
	c.Visit(episode)

	return name, episodeNum, streaming
}

// GetAjax - Find the link for the Ajax
func GetAjax(episode string) AjaxResult {
	name, episodeNum, streaming := GetInfo(episode)
	ajax := strings.Replace(streaming, "streaming", "ajax", 1)

	res := AjaxResult{
		Name:      name,
		Num:       episodeNum,
		Found:     false,
		Streaming: streaming,
	}

	// TODO: better verification
	if len(streaming) > 0 {
		res.Found = true
	}
	res.Ajax = ajax
	// domain := src[:index]
	// res.Domain = domain

	return res
}

//TODO: GET RID OF AjaxResult
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
