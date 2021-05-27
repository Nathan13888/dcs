package server

import (
	"dcs/config"
	"dcs/downloader"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

var jobs = make(map[string]*DownloadJob)
var runningJobs = 0

func AddJob(job *DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	// configure job
	job.Req.Props.Interactive = false
	job.Progress.Status = QueuedJob
	job.Progress.Completion = 0.0
	job.Date = time.Now()
	// add to queue
	jobs[job.ID] = job

	if job.Schedule.IsZero() {
		job.Schedule = time.Now().Truncate(time.Minute)
	}

	CheckJob(job)
}

func CheckJob(job *DownloadJob) {
	if !job.Schedule.IsZero() && job.Schedule.After(time.Now()) {
		return
	}
	if runningJobs > config.DownloadLimit() {
		return
	}
	if job.Progress.Status != QueuedJob {
		return
	}
	RunJob(job.ID)
}

func CheckQueue() {
	for _, j := range jobs {
		CheckJob(j)
	}
}

func RunJob(id string) {
	if !JobExists(id) {
		logError(fmt.Errorf("job with ID %s cannot be found", id))
		return
	}
	runningJobs++
	job := jobs[id]
	job.Progress.Status = RunningJob
	job.Progress.StartTime = time.Now()
	log.Info().
		Str("job", job.ID).
		Str("collection", job.Req.DInfo.Name).
		Float64("num", job.Req.DInfo.Num).
		Msg("Job starting")
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
			log.Error().
				Err(err).
				Str("job", job.ID).
				Str("collection", job.Req.DInfo.Name).
				Float64("num", job.Req.DInfo.Num).
				Msg("Job experienced an error")
		} else {
			job.Progress.Status = CompleteJob
		}
		job.Progress.EndTime = time.Now()

		log.Info().
			Str("job", job.ID).
			Str("collection", job.Req.DInfo.Name).
			Float64("num", job.Req.DInfo.Num).
			Msg("Job completed")
		runningJobs--
		CheckQueue()
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
