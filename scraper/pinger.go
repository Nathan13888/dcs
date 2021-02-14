package scraper

import (
	"fmt"
	"net"
	"time"
)

// PingDC - Ping Dramacool
func PingDC() bool {
	return Ping("watchasian.cc")
}

// Ping - Ping a website to see if it's online (defaults to HTTPS)
func Ping(url string) bool {
	// DEFAULTS to HTTPS port
	port := 443
	timeout := time.Duration(3 * time.Second)
	_, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", url, port),
		timeout,
	)
	if err != nil {
		return false
	}
	return true
}
