package scraper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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
func GetRange(r string) []int {
	// TODO: check if range is valid; valid characters: ,-0123456789
	var res []int
	exps := strings.Split(r, ",")
	var candidates []int
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
				for i := a; i <= b; i++ {
					candidates = append(candidates, i)
				}
			}
		}
	}
	for _, y := range candidates {
		if sort.Search(len(candidates), func(i int) bool { return candidates[i] == y }) >= len(res) {
			res = append(res, y)
		}
	}
	sort.Ints(res)
	return res
}

// CheckNumber - Check if a string is a number
func CheckNumber(num string) (bool, int) {
	i, err := strconv.Atoi(num)
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
