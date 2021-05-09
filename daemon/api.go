package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"dcs/config"
	"dcs/downloader"
	"dcs/scraper"

	"github.com/google/uuid"
)

type PingResponse struct {
	Uptime         float64 `json:"uptime"`
	Downloaded     int     `json:"downloaded"`
	CollectionSize int64   `json:"size"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("PING from", r.UserAgent())
	res := PingResponse{
		Uptime:     time.Since(StartTime).Minutes(),
		Downloaded: downloader.Size(config.DownloadPath()),
	}
	json.NewEncoder(w).Encode(res)
}

func getRecentDownloads(w http.ResponseWriter, r *http.Request) {
	log.Println("GET recent downloads")
	res := config.GetRecentDownloads()
	json.NewEncoder(w).Encode(res)
}

func postRecentDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("POST recent download")
	info := scraper.DramaInfo{}
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Adding new recent download `%s`\n", info.Name)

	config.AddRecentDownload(&info)

	response, err := json.Marshal(&info)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func postDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("POST download")

	var dreq DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&dreq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create
	job := DownloadJob{
		ID:     strings.ReplaceAll(uuid.New().String(), "-", ""),
		Status: RunningJob,
		Req:    dreq,
	}
	log.Printf("Adding new downloading job for '%s EPISODE %v' (%s))\n",
		dreq.DInfo.Name, dreq.DInfo.Num, job.ID)
	AddJob(job)

	// return information about job
	response, err := json.Marshal(&job)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
