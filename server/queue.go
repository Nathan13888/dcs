package server

import (
	"dcs/downloader"
	"os"

	"github.com/rs/zerolog/log"
)

var jobs = make(map[string]*DownloadJob)

func AddJob(job *DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	job.Req.Props.Interactive = false
	jobs[job.ID] = job
	job.Status = QueuedJob
}

func StartJob(id string) {
	if !JobExists(id) {
		log.Printf("WARN: Job with ID %s cannot be found", id)
		return
	}
	log.Printf("Starting job with ID %s\n", id)
	job := jobs[id]
	job.Status = RunningJob
	go func() {
		info := job.Req.DInfo
		jobLogger := getJobLogger(job)
		info.Logger = jobLogger

		err := downloader.Get(info, job.Req.Props)
		if err != nil {
			job.Status = FailedJob
			jobLogger.Error().Err(err).Msg("")
		}
		job.Status = CompleteJob
	}()
}

func GetJobInfo() ([]DownloadJob, []int64) {
	var ret []DownloadJob
	var sizes []int64
	for _, job := range jobs {
		ret = append(ret, *job)
		s, err := downloader.LookupEpisode(job.Req.DInfo)
		if err != nil && !os.IsNotExist(err) {
			logError(err)
		}
		sizes = append(sizes, s)
	}
	return ret, sizes
}

func JobExists(id string) bool {
	_, ok := jobs[id]
	return ok
}
