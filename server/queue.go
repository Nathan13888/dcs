package server

import (
	"dcs/downloader"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

var jobs = make(map[string]*DownloadJob)

func AddJob(job *DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	// configure job
	job.Req.Props.Interactive = false
	job.Progress.Status = QueuedJob
	job.Progress.Completion = 0.0
	job.Date = time.Now()
	// add to queue
	jobs[job.ID] = job

	fmt.Println(job.Schedule)
	if job.Schedule.Before(time.Now()) {
		RunJob(job.ID)
	}
}

func RunJob(id string) {
	if !JobExists(id) {
		log.Printf("WARN: Job with ID %s cannot be found", id)
		return
	}
	log.Printf("Starting job with ID %s\n", id)
	job := jobs[id]
	job.Progress.Status = RunningJob
	job.Progress.StartTime = time.Now()
	go func() {
		info := job.Req.DInfo
		info.ProgressUpdater = func(f float64) {
			job.Progress.Completion = math.Round(f*100) / 100
		}
		jobLogger := getJobLogger(job)
		info.Logger = jobLogger

		err := downloader.Get(info, job.Req.Props)
		if err != nil {
			job.Progress.Status = FailedJob
			jobLogger.Error().Err(err).Msg("")
		}
		job.Progress.Status = CompleteJob
		job.Progress.EndTime = time.Now()
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
