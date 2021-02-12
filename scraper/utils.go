package scraper

import (
	"sort"
	"strconv"
	"strings"
)

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
