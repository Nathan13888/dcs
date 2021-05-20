package downloader

import "io"

// DownloadInfo - Information you need to Download
type DownloadInfo struct {
	Link            string        `json:"link"`
	Name            string        `json:"name"`
	Num             float64       `json:"num"`
	Logger          io.Writer     `json:"-"`
	ProgressUpdater func(float64) `json:"-"`
}

type DownloadProperties struct {
	Overwrite   bool
	Interactive bool
	IgnoreM3U8  bool
	Remote      bool
}

// PathInfo - Path information about downloaded content
type PathInfo struct {
	Folder  string
	Episode string
	Dir     string
	Path    string
}
