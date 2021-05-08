package daemon

import (
	"encoding/json"
	"net/http"
	"time"

	"dcs/config"
	"dcs/downloader"
)

type PingResponse struct {
	Uptime         float64 `json:"uptime"`
	Downloaded     int     `json:"downloaded"`
	CollectionSize int64   `json:"size"`
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	res := PingResponse{
		Uptime:     time.Since(StartTime).Minutes(),
		Downloaded: downloader.Size(config.DownloadPath()),
	}
	json.NewEncoder(w).Encode(res)
}
