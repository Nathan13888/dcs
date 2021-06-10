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
