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
	RunningJob  DownloadStatus = "running"
	FailedJob   DownloadStatus = "failed"
	StaledJob   DownloadStatus = "staled"
	CompleteJob DownloadStatus = "complete"
)

type DownloadJob struct {
	ID     string          `json:"id"`
	Status DownloadStatus  `json:"status"`
	Req    DownloadRequest `json:"req"`
}
