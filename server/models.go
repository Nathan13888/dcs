package server

import (
	"dcs/downloader"
)

type StatusResponse struct {
	Uptime             float64 `json:"uptime"`
	ProcessedRequests  int64   `json:"processedRequests"`
	DownloadedDramas   int     `json:"downloadedDramas"`
	DownloadedEpisodes int     `json:"downloadedEpisodes"`
	LibrarySize        int64   `json:"size"`
}

type LogLookupResponse struct {
	Found bool     `json:"found"`
	Log   []string `json:"log"`
}

type CollectionLookupResponse struct {
	NumOfEpisodes      int      `json:"numOfEpisodes"`
	DownloadedEpisodes []string `json:"downloadedEpisodes"`
	Error              error    `json:"error"`
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
	StaledJob   DownloadStatus = "staled"
	FailedJob   DownloadStatus = "failed"
	CompleteJob DownloadStatus = "complete"
)

type DownloadJob struct {
	ID     string          `json:"id"`
	Status DownloadStatus  `json:"status"`
	Req    DownloadRequest `json:"req"`
}
