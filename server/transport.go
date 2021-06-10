package server

import (
	"dcs/downloader"
	"time"

	"gorm.io/gorm"
)

type StatusResponse struct {
	Uptime             time.Duration `json:"uptime"`
	ProcessedRequests  int64         `json:"processedRequests"`
	DownloadedDramas   int           `json:"downloadedDramas"`
	DownloadedEpisodes int           `json:"downloadedEpisodes"`
	LibrarySize        int64         `json:"size"`
	Version            string        `json:"version"`
	BuildInfo          string        `json:"buildinfo"`
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
	DInfo downloader.DownloadInfo       `gorm:"embedded;<-:create" json:"dinfo"`
	Props downloader.DownloadProperties `gorm:"embedded;<-:create" json:"props"`
}

type DownloadStatus string

const (
	QueuedJob    DownloadStatus = "queued"
	RunningJob   DownloadStatus = "running"
	FailedJob    DownloadStatus = "failed"
	CancelledJob DownloadStatus = "cancelled"
	CompleteJob  DownloadStatus = "complete"
)

type ProgressInfo struct {
	Completion float64        `json:"completion"`
	StartTime  time.Time      `json:"startTime"`
	EndTime    time.Time      `json:"endTime"`
	Status     DownloadStatus `gorm:"index" json:"status"`
}

type DownloadJob struct {
	gorm.Model
	ID       string          `gorm:"primaryKey" json:"id"`
	Date     time.Time       `gorm:"<-:create" json:"date"`
	Schedule time.Time       `gorm:"<-:update" json:"scheduledTime"`
	Progress ProgressInfo    `gorm:"embedded;embeddedPrefix:progress_" json:"progress"`
	Req      DownloadRequest `gorm:"embedded" json:"req"`
}

type JobsResponse struct {
	Jobs  []DownloadJob `json:"jobs"`
	Sizes []int64       `json:"sizes"`
	// Num  int           `json:"numOfJobs"`
}
