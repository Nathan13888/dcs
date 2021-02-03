package scraper

import "strings"

// JoinArgs - Join arguments
func JoinArgs(args []string) string {
	return strings.ReplaceAll(strings.Join(args, " "), "\"", "")
}
