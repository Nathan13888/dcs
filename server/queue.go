package server

import (
	"dcs/config"
	"dcs/downloader"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var runningJobs = 0

func AddJob(job *DownloadJob) {
	log.Printf("Recording new job with ID %s\n", job.ID)

	// configure job
	job.Req.Props.Interactive = false
	job.Progress.Status = QueuedJob
	job.Progress.Completion = 0.0
	job.Date = time.Now()
	if job.Schedule.IsZero() {
		job.Schedule = time.Now().Truncate(time.Minute)
	}
	DBAddJob(job)

	CheckJob(job)
}

func CheckJob(job *DownloadJob) {
	if runningJobs > config.DownloadLimit() {
		return
	}
	if !job.Schedule.IsZero() && job.Schedule.After(time.Now()) {
		return
	}
	if job.Progress.Status != QueuedJob {
		return
	}
	RunJob(job)
}

func CheckQueue() {
	for _, j := range DBGetJobs() {
		CheckJob(&j)
	}
}

func RunJob(job *DownloadJob) {
	runningJobs++
	job.Progress.Status = RunningJob
	job.Progress.StartTime = time.Now()
	DBUpdateJob(job)
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
		jobLogger := getJobLogger(job.ID)
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
		DBUpdateJob(job)

		log.Info().
			Str("job", job.ID).
			Str("collection", job.Req.DInfo.Name).
			Float64("num", job.Req.DInfo.Num).
			Msg("Job completed")
		runningJobs--
		CheckQueue()
	}()
}

func RunUncompletedJobs() {
	for _, job := range DBGetJobs() {
		if job.Progress.Status == FailedJob {
			continue
		}
		if job.Progress.Completion != 100.0 || job.Progress.Status == RunningJob {
			RunJob(&job)
		}
	}
}

func GetJobInfo() ([]DownloadJob, []int64) {
	var jobs = DBGetJobs()
	var sizes []int64
	for _, job := range jobs {
		s, err := downloader.LookupEpisode(job.Req.DInfo)
		if err != nil && !os.IsNotExist(err) {
			logError(err)
		}
		sizes = append(sizes, s)
	}
	return jobs, sizes
}

func GenerateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
