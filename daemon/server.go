package daemon

import (
	"log"
	"net/http"
)

var port = ":6969"

// Start - start a HTTP API service
func Start() {
	log.Printf("Starting DCS Daemon at port %s\n\n", port[1:])

	http.HandleFunc("/", handlePing)
	http.HandleFunc("/ping", handlePing)
	log.Fatal(http.ListenAndServe(port, nil))
}
