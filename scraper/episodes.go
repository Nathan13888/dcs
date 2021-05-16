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

	c.OnHTML("ul.all-episode", func(e *colly.HTMLElement) {
		e.ForEach("li a.img h3.title", func(i int, ee *colly.HTMLElement) {
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
