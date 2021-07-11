package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"dcs/config"
	"dcs/downloader"
	"dcs/scraper"

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
		Uptime:             time.Since(StartTime),
		ProcessedRequests:  processedRequests,
		DownloadedDramas:   dd,
		DownloadedEpisodes: de,
		LibrarySize:        csize,
		Version:            fmt.Sprintf("%s (%s/%s)", config.BuildVersion, config.BuildGOOS, config.BuildARCH),
		BuildInfo:          fmt.Sprintf("Built on %s (by %s)", config.BuildTime, config.BuildUser),
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

	// create new job
	job := DownloadJob{
		ID:  GenerateID(),
		Req: dreq,
	}
	log.Info().
		Str("job", job.ID).
		Str("collection", dreq.DInfo.Name).
		Float64("num", job.Req.DInfo.Num).
		Msg("New job")

	AddJob(&job)

	// return information about job
	response, err := json.Marshal(&job)
	if err != nil {
		internalError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func getRemoveJob(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// job, exists := DBGetJob(id)
	// if !exists || job.ID != id {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	deleteJob(id, w)
}

func deleteJob(id string, w http.ResponseWriter) {
	res := GetDB().Delete(&DownloadJob{
		ID: id,
	})
	if res.Error != nil {
		internalError(w, res.Error)
		return
	}
	log.Info().Str("job", id).Msg("Deleting job")
}

func getPurgeJobs(w http.ResponseWriter, r *http.Request) {
	jobs, sizes := GetJobInfo()
	purged := JobsResponse{
		Jobs:  make([]DownloadJob, 0),
		Sizes: make([]int64, 0),
	}
	for i, j := range jobs {
		if j.Progress.Status == FailedJob ||
			(j.Progress.Status == RunningJob && time.Since(j.Progress.StartTime) > (time.Hour*24)) {
			deleteJob(j.ID, w)
			purged.Jobs = append(purged.Jobs, j)
			purged.Sizes = append(purged.Sizes, sizes[i])
		}
	}
	json.NewEncoder(w).Encode(purged)
}

func getJobsList(w http.ResponseWriter, r *http.Request) {
	// TODO: verify sizes returned after equal??
	jobs, sizes := GetJobInfo()
	res := JobsResponse{
		Jobs:  jobs,
		Sizes: sizes,
	}
	json.NewEncoder(w).Encode(res)
}

func internalError(w http.ResponseWriter, err error) {
	logError(err)
	w.WriteHeader(http.StatusInternalServerError)
}
