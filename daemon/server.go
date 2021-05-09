package daemon

import (
	"dcs/config"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// StartTime - start time of the daemon
var StartTime time.Time

// Start - start a HTTP API service
func Start() {
	_, port := config.DaemonURL()
	log.Printf("Starting DCS Daemon at port %d\n\n", port)

	r := mux.NewRouter()
	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	r.HandleFunc("/", handlePing).Methods("GET")
	r.HandleFunc("/ping", handlePing).Methods("GET")
	r.HandleFunc("/status", handlePing).Methods("GET")
	r.HandleFunc("/api/recentdownloads", getRecentDownloads).Methods("GET")
	r.HandleFunc("/api/recentdownload", postRecentDownload).Methods("POST")
	r.HandleFunc("/api/download", postDownload).Methods("POST")

	StartTime = time.Now()

	log.Fatal(server.ListenAndServe())
}
