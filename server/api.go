package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"dcs/config"
	"dcs/downloader"
	"dcs/scraper"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PingResponse struct {
	Uptime         float64 `json:"uptime"`
	Downloaded     int     `json:"downloaded"`
	CollectionSize int64   `json:"size"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	res := PingResponse{
		Uptime:     time.Since(StartTime).Minutes(),
		Downloaded: downloader.Size(config.DownloadPath()),
	}
	json.NewEncoder(w).Encode(res)
}

func getRecentDownloads(w http.ResponseWriter, r *http.Request) {
	res := config.GetRecentDownloads()
	json.NewEncoder(w).Encode(res)
}

func postRecentDownload(w http.ResponseWriter, r *http.Request) {
	info := scraper.DramaInfo{}
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info().Msgf("Adding new recent download `%s`", info.Name)

	config.AddRecentDownload(&info)

	response, err := json.Marshal(&info)
	if err != nil {
		logError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func postDownload(w http.ResponseWriter, r *http.Request) {
	var dreq DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&dreq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create
	job := DownloadJob{
		ID:     strings.ReplaceAll(uuid.New().String(), "-", ""),
		Status: QueuedJob,
		Req:    dreq,
	}
	log.Info().Msgf("Adding new downloading job for '%s EPISODE %v' (%s)",
		dreq.DInfo.Name, dreq.DInfo.Num, job.ID)
	AddJob(&job)
	StartJob(job.ID)

	// return information about job
	response, err := json.Marshal(&job)
	if err != nil {
		logError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
