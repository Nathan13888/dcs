package server

import (
	"dcs/downloader"

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

func JobExists(id string) bool {
	_, ok := jobs[id]
	return ok
}
