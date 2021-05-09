package daemon

import "log"

var jobs = make(map[string]DownloadJob)

func AddJob(job DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	jobs[job.ID] = job
	// StartJob(job.ID)
}
