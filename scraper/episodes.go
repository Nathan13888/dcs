package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// EpisodeInfo - Information about episodes
type EpisodeInfo struct {
	Number float64
	Date   string
	Link   string
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

// GetEpisodesByLink - GetEpisodes by just give a link
func GetEpisodesByLink(link string) []EpisodeInfo {
	// TODO: use GetInfo() instead
	c := getCollector()

	var name string

	// TODO: FIX for asianload `div.video-details span.date`
	c.OnHTML("div.info h1", func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	c.Visit(link)
	return GetEpisodes(DramaInfo{
		FullURL: link,
		Name:    EscapeName(name),
	})
}

// GetEpisodes - Tells you all the available episodes
func GetEpisodes(drama DramaInfo) []EpisodeInfo {
	// fmt.Printf("\nFetching episodes of `%s`\n\n", drama.Name)
	episodes := []EpisodeInfo{}

	c := getCollector()
	// TODO: cache page

	var ulQry string
	var liQry string
	if ASIANLOAD {
		ulQry = "ul.listing.lists"
		liQry = "li a div.name"
	} else {
		ulQry = "ul.all-episode"
		liQry = "li a.img h3.title"
	}

	c.OnHTML(ulQry, func(e *colly.HTMLElement) {
		e.ForEach(liQry, func(i int, ee *colly.HTMLElement) {
			parent := ee.DOM.Parent()
			link, _ := parent.Attr("href")
			fullname := strings.Trim(ee.DOM.Contents().Text(), " \n")
			time := strings.Trim(parent.ChildrenFiltered("span.time").Text(), " \n")
			// fmt.Printf("%s was posted %s\n", fullname, time)
			split := strings.Split(fullname, " ")
			parse := split[len(split)-1]
			num, err := strconv.ParseFloat(parse, 64)
			if err != nil {
				fmt.Println("ERROR: The episode number could not be interpreted...")
				fmt.Println(err)
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
