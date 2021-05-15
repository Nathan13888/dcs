package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"dcs/config"
	"dcs/downloader"
	"dcs/scraper"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func getPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	dd, err := downloader.DownloadedDramas()
	if err != nil {
		logError(err)
	}
	de, err := downloader.DownloadedEpisodes()
	if err != nil {
		logError(err)
	}
	csize, err := downloader.DirSize("")
	if err != nil {
		logError(err)
	}
	res := StatusResponse{
		Uptime:             time.Since(StartTime).Seconds(),
		ProcessedRequests:  processedRequests,
		DownloadedDramas:   dd,
		DownloadedEpisodes: de,
		LibrarySize:        csize,
	}
	json.NewEncoder(w).Encode(res)
}

func getLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, found := vars["type"]
	if !found || !(t == "job" || t == "server") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, found := vars["id"]
	if !found && t != "server" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var p string
	if t == "server" {
		p = logFile.Name()
	} else { // job
		p = path.Join(logDir, getJobLogName(id))
	}

	res := LogLookupResponse{
		Found: true,
		Log:   []string{},
	}

	f, err := os.Open(p)
	if err != nil {
		internalError(w, err)
		res.Found = false
	} else {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			res.Log = append(res.Log, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			internalError(w, err)
			res.Found = false
		}
	}
	json.NewEncoder(w).Encode(res)
}

func getRecentDownloads(w http.ResponseWriter, r *http.Request) {
	res := config.GetRecentDownloads()
	json.NewEncoder(w).Encode(res)
}

func getLookupCollection(w http.ResponseWriter, r *http.Request) {
	name, found := mux.Vars(r)["name"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	num, epNames, err := downloader.CollectionLookup(name)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		internalError(w, err)
		return
	}
	size, err := downloader.DirSize(name)
	if err != nil {
		internalError(w, err)
		return
	}
	res := CollectionLookupResponse{
		NumOfEpisodes:      num,
		DownloadedEpisodes: epNames,
		Error:              err,
		Size:               size,
	}
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
		internalError(w, err)
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
		internalError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func internalError(w http.ResponseWriter, err error) {
	logError(err)
	w.WriteHeader(http.StatusInternalServerError)
}
