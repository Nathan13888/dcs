package scraper

import (
	"fmt"
	"net"
	"time"
)

// PingDC - Ping Dramacool
func PingDC() bool {
	return Ping("watchasian.cc", 443)
}

// TODO: update what is this about
// Ping - Ping a website to see if it's online (defaults to HTTPS)
func Ping(url string, port int) bool {
	// DEFAULTS to HTTPS port
	timeout := time.Duration(3 * time.Second)
	_, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%d", url, port),
		timeout,
	)
	return err == nil
}
