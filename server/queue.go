package server

import (
	"dcs/downloader"
	"log"
)

var jobs = make(map[string]*DownloadJob)

func AddJob(job DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	jobs[job.ID] = &job
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
		// jobLogger:=log.New()
		downloader.Get(job.Req.DInfo, job.Req.Props)
	}()
}

func JobExists(id string) bool {
	_, ok := jobs[id]
	return ok
}
