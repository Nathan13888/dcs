package scraper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
)

var ldBase = "https://asianload.io/download?id="

type dlink struct {
	Label string
	Link  string
}

func ScrapeLD(id string) string {
	u := ldBase + id
	fmt.Println("Scraping for download links in ", u)
	c := getCollector()

	n := 0
	found := false
	url := ""
	links := make([]dlink, 0)

	c.OnHTML("div .dowload a[href]", func(e *colly.HTMLElement) {
		n++
		src := strings.TrimSpace(e.Attr("href"))
		links = append(links, dlink{
			Label: strings.TrimSpace(e.DOM.Text()),
			Link:  src,
		})
		if strings.HasSuffix(src, "mp4") && strings.Contains(src, "storage.googleapis.com") {
			found = true
			url = src
		}
	})

	c.Visit(u)

	fmt.Println("Number of Links found:", n)
	if found {
		fmt.Println("Found Link:", url)
	} else {
		fmt.Println("Link NOT FOUND")
		templates := &promptui.SelectTemplates{
			Label:    "{{ . | white | bold }}",
			Active:   "\U0001F449 {{ .Label | green | bold }} ({{ .Link | green | bold }})",
			Inactive: "  {{ .Label | red }}",
			Selected: "\U0001F449 {{ .Label | green }}",
			Details: `

{{ .Label | blue | bold}}
{{ "----------------------------" | white }}
{{ "Link:" | faint }}	{{ .Link | yellow }}`,
		}
		searcher := func(input string, index int) bool {
			l := links[index]
			return len(fuzzy.FindNormalizedFold(input, []string{
				l.Label, l.Link,
			})) > 0
		}

		// settings for the prompt
		p := promptui.Select{
			Label:     "Select a download link:",
			Items:     links,
			Templates: templates,
			Size:      9,
			Searcher:  searcher,
		}

		i, _, err := p.Run()
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}
		link := links[i].Link
		fmt.Println("Selected link: ", link)
		// url, err = prompt.String("Please manually enter download link:")
		pt := promptui.Prompt{
			Label: "Please manually enter download link:",
			Validate: func(input string) error {
				if len(strings.TrimSpace(input)) == 0 {
					return errors.New("invalid empty imput")
				}
				return nil
			},
		}
		url, err = pt.Run()
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}
	}

	return url
}
