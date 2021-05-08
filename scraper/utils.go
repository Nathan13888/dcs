package scraper

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

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

// GetRange - Determine which numbers are included in a "range"; returns [] if range is invalid
func GetRange(r string) []float64 {
	// TODO: check if range is valid; valid characters: ,-0123456789
	var res []float64
	exps := strings.Split(r, ",")
	var candidates []float64
	for _, x := range exps {
		isNum, num := CheckNumber(x)
		if isNum {
			candidates = append(candidates, num)
		} else {
			// TODO: expression must be two numbers separated by a -
			split := strings.Split(x, "-")
			isAValid, a := CheckNumber(split[0])
			isBValid, b := CheckNumber(split[1])
			if isAValid && isBValid {
				candidates = append(candidates, a) // add starting position
				candidates = append(candidates, b) // add ending position
				for i := math.Ceil(a + 1); i < b; i++ {
					candidates = append(candidates, i)
				}
			}
		}
	}
	// filter duplicate values
	for _, y := range candidates {
		if sort.Search(len(candidates), func(i int) bool { return candidates[i] == y }) >= len(res) {
			res = append(res, y)
		}
	}
	sort.Float64s(res)
	return res
}

// JoinURL combines several parts of URLs together into a string
func JoinURL(a string, b string) string {
	u, _ := url.Parse(a)
	u.Path = path.Join(u.Path, b)
	return u.String()
}

// CheckNumber - Check if a string is a number
func CheckNumber(num string) (bool, float64) {
	i, err := strconv.ParseFloat(num, 64)
	if err == nil {
		return true, i
	}
	return false, i
}

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
