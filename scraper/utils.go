package scraper

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

var ASIANLOAD bool = false
var PROXIES []string = make([]string, 0)
var NOPROXY bool

func ConfigProxies(proxies []string) {
	if !NOPROXY {
		PROXIES = proxies
		// for _, p := range proxies {
		// 	// TODO: proxy validation
		// 	PROXIES = append(PROXIES, p)
		// }
	}

	fmt.Println("Using Proxies:", PROXIES)
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
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
		r.Save("tmp.error.html")
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", USERAGENT)
	})

	if len(PROXIES) > 0 {
		if p, err := proxy.RoundRobinProxySwitcher(PROXIES...); err == nil {
			c.SetProxyFunc(p)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("NO PROXIES are being used.")
	}

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

func EscapeName(name string) string {
	var sb strings.Builder
	runes := []rune{' ', '(', ')', '[', ']', '.'}
	for _, c := range name {
		// search `runes`
		found := false
		for _, r := range runes {
			if c == r {
				found = true
				break
			}
		}
		if found || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
			sb.WriteRune(c)
		}
	}
	escaped := strings.TrimSpace(sb.String())
	escaped = regexp.MustCompile(`\s+`).ReplaceAllString(escaped, " ")
	return escaped
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
