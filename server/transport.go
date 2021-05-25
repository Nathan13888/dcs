package server

import (
	"dcs/downloader"
	"time"
)

type StatusResponse struct {
	Uptime             float64 `json:"uptime"`
	ProcessedRequests  int64   `json:"processedRequests"`
	DownloadedDramas   int     `json:"downloadedDramas"`
	DownloadedEpisodes int     `json:"downloadedEpisodes"`
	LibrarySize        int64   `json:"size"`
	Version            string  `json:"version"`
	BuildInfo          string  `json:"buildinfo"`
}

type LogLookupResponse struct {
	Found bool     `json:"found"`
	Log   []string `json:"log"`
}

type CollectionLookupResponse struct {
	NumOfEpisodes      int      `json:"numOfEpisodes"`
	DownloadedEpisodes []string `json:"downloadedEpisodes"`
	Error              error    `json:"err,omitempty"`
	Size               int64    `json:"size"`
}

type DownloadRequest struct {
	DInfo downloader.DownloadInfo       `json:"dinfo"`
	Props downloader.DownloadProperties `json:"props"`
}

type DownloadStatus string

const (
	QueuedJob   DownloadStatus = "queued"
	RunningJob  DownloadStatus = "running"
	FailedJob   DownloadStatus = "failed"
	CompleteJob DownloadStatus = "complete"
)

type ProgressInfo struct {
	Completion float64        `json:"completion"`
	StartTime  time.Time      `json:"startTime"`
	EndTime    time.Time      `json:"endTime"`
	Status     DownloadStatus `json:"status"`
}

type DownloadJob struct {
	ID       string          `json:"id"`
	Date     time.Time       `json:"date"`
	Schedule time.Time       `json:"scheduledTime"`
	Progress ProgressInfo    `json:"progress"`
	Req      DownloadRequest `json:"req"`
}

type JobsResponse struct {
	Jobs  []DownloadJob `json:"jobs"`
	Sizes []int64       `json:"sizes"`
	// Num  int           `json:"numOfJobs"`
}
