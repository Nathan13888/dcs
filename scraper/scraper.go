package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// TODO: debug level

// TODO: multiple backup URLs
const URL = "https://watchasian.cc"

// DramaInfo - Information about a drama
type DramaInfo struct {
	FullURL string
	SubURL  string
	Domain  string
	Name    string
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

// EpisodeInfo - Information about episodes
type EpisodeInfo struct {
	Number int
	Date   string
	Link   string
}

// GetEpisodesByLink - GetEpisodes by just give a link
func GetEpisodesByLink(link string) []EpisodeInfo {
	// TODO: use GetInfo() instead
	c := getCollector()

	var name string

	c.OnHTML("div.info h1", func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	c.Visit(link)
	return GetEpisodes(DramaInfo{
		FullURL: link,
		Name:    name,
	})
}

// GetEpisodes - Tells you all the available episodes
func GetEpisodes(drama DramaInfo) []EpisodeInfo {
	// fmt.Printf("\nFetching episodes of `%s`\n\n", drama.Name)
	episodes := []EpisodeInfo{}

	c := getCollector()
	// TODO: cache page

	c.OnHTML("ul.all-episode", func(e *colly.HTMLElement) {
		e.ForEach("li a.img h3.title", func(i int, ee *colly.HTMLElement) {
			parent := ee.DOM.Parent()
			link, _ := parent.Attr("href")
			fullname := strings.Trim(ee.DOM.Contents().Text(), " \n")
			time := strings.Trim(parent.ChildrenFiltered("span.time").Text(), " \n")
			// fmt.Printf("%s was posted %s\n", fullname, time)
			num, err2 := strconv.Atoi(fullname[len(drama.Name)+9:])
			if err2 != nil {
				fmt.Println("ERROR: The episode number could not be interpreted...")
				fmt.Print(err2)
			} else {
				obj := EpisodeInfo{
					Number: num,
					Date:   time, //TODO: finish implementing time
					Link:   link,
				}
				episodes = append(episodes, obj)
			}
		})
	})

	c.Visit(drama.FullURL)

	return episodes
}

// AjaxResult - The result of scraping the ajax url from an episode link
type AjaxResult struct {
	Found     bool
	Name      string
	Num       int
	Ajax      string
	Streaming string
	Domain    string
}

// GetAjax - Find the link for the Ajax
func GetAjax(episode string) AjaxResult {
	res := AjaxResult{
		Found: false,
	}

	c := getCollector()

	c.OnHTML("div.watch-drama h1", func(e *colly.HTMLElement) {
		num := strings.Split(strings.Trim(e.DOM.Text(), " \n"), " ")
		// fmt.Println(num)
		conv, err := strconv.Atoi(num[len(num)-1])
		if err != nil {
			fmt.Println(err)
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

// AjaxResponse - The expected JSON response from the DC Ajax endpoint
type AjaxResponse struct {
	Source   []SourceOption `json:"source"`
	SourceBK []SourceOption `json:"source_bk"`
	Track    []TrackOption  `json:"track"`
	// Advertising []string
	LinkIFrame string `json:"linkiframe"`
}

// ScrapeEpisode - Get the link to the video source of an episode
func ScrapeEpisode(ajax AjaxResult) string {
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

func printObj(obj interface{}) {
	resJSON, _ := json.MarshalIndent(obj, "  ", "    ")
	fmt.Println(string(resJSON))
}

func getCollector() *colly.Collector {
	c := colly.NewCollector() // colly.Debugger(&debug.LogDebugger{})

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Printf("\nVisiting: %s\n\n", r.URL.String())
	})

	return c
}
