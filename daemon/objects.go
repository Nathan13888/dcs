package daemon

import (
	"dcs/downloader"
)

type DownloadRequest struct {
	DInfo downloader.DownloadInfo       `json:"dinfo"`
	Props downloader.DownloadProperties `json:"props"`
}

type DownloadStatus string

const (
	QueuedJob   DownloadStatus = "queued"
	RunningJob  DownloadStatus = "running"
	StaledJob   DownloadStatus = "staled"
	FailedJob   DownloadStatus = "failed"
	CompleteJob DownloadStatus = "complete"
)

type DownloadJob struct {
	ID     string          `json:"id"`
	Status DownloadStatus  `json:"status"`
	Req    DownloadRequest `json:"req"`
}
