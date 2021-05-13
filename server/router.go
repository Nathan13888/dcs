package server

import (
	"context"
	"dcs/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// StartTime - start time of the daemon
var StartTime time.Time
var DEBUGMODE = false
var processedRequests int64 = 0

// Start - start a HTTP API service
func Start(debug bool) {
	if debug {
		DEBUGMODE = true
	}

	_, port := config.DaemonURL()

	r := mux.NewRouter()
	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	r.HandleFunc("/", getStatus).Methods("GET")
	r.HandleFunc("/ping", getPing).Methods("GET")
	r.HandleFunc("/status", getStatus).Methods("GET")
	r.HandleFunc("/api/recentdownloads", getRecentDownloads).Methods("GET")
	r.HandleFunc("/api/recentdownload", postRecentDownload).Methods("POST")
	r.HandleFunc("/api/download", postDownload).Methods("POST")

	r.Use(loggingMiddleware)
	closeLogger := configureLogger()
	defer closeLogger()

	config.IS_SERVER = true
	StartTime = time.Now()

	log.Info().Msgf("Starting DCS Daemon at port %d", port)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logError(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	wait := 15 * time.Second

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info().Msg("shutting down...")
	os.Exit(0)
}
